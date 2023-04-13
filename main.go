package main

import (
	"fmt"
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

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// lc := new(siologger.LokiConfig)
	//
	// lc.UseDefaults("customer-ms")
	// fmt.Println(lc)
	//
	// lh, err := siologger.NewLokiHook(lc)
	// if err != nil {
	// 	log.Errorf("Error creating Loki hook: %s", err)
	// 	return
	// }
	log.AddHook(lh)
}

func main() {
	fmt.Println(os.Environ())
	CreateRouter()
}
