package main

import (
	"context"
	"go-test/db"
	"go-test/server"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	PORT, DB_URL, JWT_ACCESS_SECRET, JWT_REFRESH_SECRET := setupEnv()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(DB_URL))
	if err != nil {
		log.Fatal("ERROR [DB]:", err)
	}
	defer client.Disconnect(context.Background())

	log.Println("DB connected")
	db.NewDAO(*client)

	server.InitializeServer(PORT, JWT_ACCESS_SECRET, JWT_REFRESH_SECRET)
}
