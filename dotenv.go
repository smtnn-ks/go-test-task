package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func setupEnv() (PORT, DB_URL, JWT_ACCESS_SECRET, JWT_REFRESH_SECRET string) {
	godotenv.Load(".env")
	PORT = os.Getenv("PORT")
	DB_URL = os.Getenv("DB_URL")
	JWT_ACCESS_SECRET = os.Getenv("JWT_ACCESS_SECRET")
	JWT_REFRESH_SECRET = os.Getenv("JWT_REFRESH_SECRET")

	flag := false
	if PORT == "" {
		flag = true
		log.Println("ERROR: PORT is not specified")
	}
	if DB_URL == "" {
		flag = true
		log.Println("ERROR: DB_URL is not specified")
	}
	if JWT_ACCESS_SECRET == "" {
		flag = true
		log.Println("ERROR: JWT_ACCESS_SECRET is not specified")
	}
	if JWT_REFRESH_SECRET == "" {
		flag = true
		log.Println("ERROR: JWT_REFRESH_SECRET is not specified")
	}

	if flag {
		os.Exit(1)
	}

	return
}
