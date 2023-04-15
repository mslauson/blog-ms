package main

import (
	"os"

	"gitea.slauson.io/slausonio/sio-loki/hooks"
	log "github.com/sirupsen/logrus"
)

type logBatch struct {
	entries      []*log.Entry
	maxEntries   int
	maxTimeLimit int
	sendSchedule chan bool
	labels       map[string]string
}

func init() {
	lh := hooks.NewLokiHook(10, 5, map[string]string{"app": "iam-ms", "environment": os.Getenv("ENV")})
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	log.AddHook(lh)
}

func main() {
	CreateRouter()
}
