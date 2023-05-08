package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	_ "gitea.slauson.io/slausonio/go-types/siogeneric"
	_ "gitea.slauson.io/slausonio/iam-ms/docs"
	"gitea.slauson.io/slausonio/sio-loki/hooks"
)

func init() {
	lh := hooks.NewLokiHook(
		10,
		5,
		map[string]string{"app": "iam-ms", "environment": os.Getenv("ENV")},
	)
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	log.AddHook(lh)
}

// @title IAM Microservice
// @description This MS handles all IAM related requests with the IAM provider
// @version 1.0

// @contact.name Matthew Slauson
// @contact.email matthew@slauson.io
func main() {
	r := CreateRouter()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
