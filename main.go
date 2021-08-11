package main

import (
	"log"
	"os"

	"artemkv.net/journey2/server"
	"github.com/joho/godotenv"
)

func main() {
	loadDotEnv()
	port := getPort(":8600")

	server.Serve(port)
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func getPort(def string) string {
	port := os.Getenv("JOURNEY2_PORT")
	if port == "" {
		log.Printf("Using default port %s\n", def)
		return def
	}
	return port
}
