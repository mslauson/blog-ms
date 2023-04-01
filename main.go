package main

import (
	"fmt"
	"os"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
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

	lc := new(sioUtils.LokiConfig)

	lc.UseDefaults("customer-ms")
	fmt.Println(lc)

	lh, err := sioUtils.NewLokiHook(lc)
	if err != nil {
		log.Errorf("Error creating Loki hook: %s", err)
		return
	}
	log.AddHook(lh)
}

func main() {
	CreateRouter()
}
