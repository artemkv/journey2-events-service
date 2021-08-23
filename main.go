package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"artemkv.net/journey2/app"
	"artemkv.net/journey2/health"
	"artemkv.net/journey2/reststats"
	"artemkv.net/journey2/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var version = "1.0"

func main() {
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	LoadDotEnv()
	port := GetPort(":8600")

	router := gin.New()
	app.SetupRouter(router)

	reststats.Initialize(version)

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
