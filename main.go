package main

import (
	"net/http"
	"os"

	"gitea.slauson.io/blog/blog-ms/handler"
	"gitea.slauson.io/slausonio/go-prom/sioprom"
	"gitea.slauson.io/slausonio/sio-loki/hooks"

	// _ "gitea.slauson.io/slausonio/customer-ms/docs"

	_ "gitea.slauson.io/slausonio/go-types/siogeneric"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	lh := hooks.NewLokiHook(
		10,
		5,
		map[string]string{"app": "customer-ms", "environment": os.Getenv("ENV")},
	)

	log.AddHook(lh)
}

// @title Blog Microservice
// @description This MS handles blog posts and comments
// @version 1.0

// @contact.name Matthew Slauson
// @contact.email matthew@slauson.io
func main() {
	go func() { sioprom.InitPrometheus() }()
	r := handler.CreateRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
