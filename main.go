package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type logBatch struct {
	entries      []*logrus.Entry
	maxEntries   int
	maxTimeLimit int
	sendSchedule chan bool
	labels       map[string]string
}

func (b *logBatch) Flush() error {
	request := &LokiRequest{
		Streams: []LokiStream{},
	}
	for _, entry := range b.entries {
		labels := b.labels
		labels["level"] = entry.Level.String()

		// Add a stack trace label if the entry has a stack trace.
		if entry.HasCaller() {
			pc := make([]uintptr, 15)
			n := runtime.Callers(6, pc)
			frames := runtime.CallersFrames(pc[:n])
			var buf bytes.Buffer
			for {
				frame, more := frames.Next()
				fmt.Printf("%s:%d\n", frame.File, frame.Line)
				if !more {
					break
				}
			}
			labels["stack_trace"] = buf.String()
		}

		values := [][]string{
			{
				strconv.FormatInt(entry.Time.UnixNano(), 10),
				entry.Message,
			},
		}
		stream := LokiStream{
			Stream: labels,
			Values: values,
		}
		request.Streams = append(request.Streams, stream)
	}

	// Send the batched log entries to Loki.
	client := NewLokiClient()
	if err := client.SendLog(request); err != nil {
		return err
	}
	b.entries = nil
	return nil
}

func (b *logBatch) startPeriodicFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := b.Flush(); err != nil {
				log.Error(err)
			}
		}
	}()
}

type LokiHook struct {
	batch *logBatch
}

func NewLokiHook(maxEntries int, interval int, labels map[string]string) *LokiHook {
	h := &LokiHook{
		batch: &logBatch{
			maxEntries: maxEntries,
			labels:     labels,
		},
		// Other fields for the hook
	}
	h.batch.startPeriodicFlush(time.Duration(interval) * time.Second)

	return h
}

func (h *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LokiHook) Fire(entry *logrus.Entry) error {
	// If the batch is nil, create a new batch and set the creation time.
	if h.batch == nil {
		h.batch = &logBatch{
			maxEntries:   h.batch.maxEntries,
			maxTimeLimit: h.batch.maxTimeLimit,
			sendSchedule: make(chan bool),
		}
	}

	// Add the log entry to the batch.
	h.batch.entries = append(h.batch.entries, entry)

	// If the batch is full, flush it and create a new batch.
	if len(h.batch.entries) >= h.batch.maxEntries {
		err := h.batch.Flush()
		if err != nil {
			return err
		}
		h.batch = nil
	}

	return nil
}

type LokiClient struct {
	rh *sioUtils.RestHelpers
}

// type LokiClient interface{}
func NewLokiClient() *LokiClient {
	return &LokiClient{
		rh: sioUtils.NewRestHelpers(),
	}
}

func (lc *LokiClient) SendLog(request *LokiRequest) error {
	rh := sioUtils.NewRestHelpers()

	url := fmt.Sprintf("%s/loki/api/v1/push", os.Getenv("LOKI_URL"))

	rJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}
	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)
	req.Header.Add("Content-Type", "application/json")

	res, err := rh.DoHttpRequest(req)
	if err != nil {
		return err
	}

	log.Debugln("Response: %s", res)

	return nil
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiRequest struct {
	Streams []LokiStream `json:"streams"`
}

func init() {
	lh := NewLokiHook(10, 5, map[string]string{"app": "customer-ms", "environment": os.Getenv("ENV")})
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
	CreateRouter()
}
