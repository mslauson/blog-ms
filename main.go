package main

import (
	"net/http"
	"os"

	"gitea.slauson.io/slausonio/sio-loki/hooks"
	log "github.com/sirupsen/logrus"
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

func main() {
	r := CreateRouter()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
