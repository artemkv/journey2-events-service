package main

import (
	"log"
	"os"

	"artemkv.net/journey2/app"
	"artemkv.net/journey2/health"
	"artemkv.net/journey2/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	LoadDotEnv()
	port := GetPort(":8600")

	router := gin.Default() // logging and recovery attached
	app.SetupRouter(router)

	server.Serve(router, port, func() {
		health.SetIsReadyGlobally()
	})
}

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func GetPort(def string) string {
	port := os.Getenv("JOURNEY2_PORT")
	if port == "" {
		log.Printf("Using default port %s\n", def)
		return def
	}
	return port
}
