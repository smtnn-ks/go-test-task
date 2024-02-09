package server

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-test/db"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inputId is not provided"))
		return
	}

	accessToken, refreshToken, err := generateTokens(userId)

	if err != nil {
		log.Println("ERROR:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	refHash := sha512.Sum512([]byte(refreshToken))
	refHashString := base64.StdEncoding.EncodeToString(refHash[:])

	payload, err := json.Marshal(Tokens{accessToken, refreshToken})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = db.DAOInstance.Insert(db.Record{
		Id:           userId,
		RefreshToken: refHashString,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("such ID is taken already"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Authorization header is not provided"))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secrets.jwtRefreshSecret), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	var userId string
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	userId = claims["sub"].(string)
	expiresAt := int64(claims["exp"].(float64))

	if expiresAt < time.Now().UTC().Unix() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	user, err := db.DAOInstance.FindById(userId)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	refHash := sha512.Sum512([]byte(tokenString))
	refHashString := base64.StdEncoding.EncodeToString(refHash[:])

	if refHashString != user.RefreshToken {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	accessToken, refreshToken, err := generateTokens(userId)

	if err != nil {
		log.Println("ERROR:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	refHash = sha512.Sum512([]byte(refreshToken))
	refHashString = base64.StdEncoding.EncodeToString(refHash[:])

	err = db.DAOInstance.Update(userId, refHashString)

	if err != nil {
		fmt.Println("ERROR:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	payload, err := json.Marshal(Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	if err != nil {
		fmt.Println("ERROR:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}
