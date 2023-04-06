package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	createdAt    time.Time
	sendSchedule chan bool
}

func (b *logBatch) flush() error {
	// Send the batched log entries to their destination here.
	// ...
	// Clear the batched entries after they've been sent.
	b.entries = nil
	return nil
}

type lokiHook struct {
	batch *logBatch
}

func NewLokiHook(maxEntries int, timeLimit int) *lokiHook {
	return &lokiHook{
		batch: &logBatch{
			maxEntries:   maxEntries,
			maxTimeLimit: timeLimit,
			sendSchedule: make(chan bool),
			createdAt:    time.Now(),
		},
		// Other fields for the hook
	}
}

func (h *lokiHook) Fire(entry *logrus.Entry) error {
	// If the batch is nil, create a new batch and set the creation time.
	if h.batch == nil {
		h.batch = &logBatch{
			maxEntries:   h.batch.maxEntries,
			maxTimeLimit: h.batch.maxTimeLimit,
			createdAt:    time.Now(),
			sendSchedule: make(chan bool),
		}
		go func() {
			// Wait for the sendSchedule signal or the time limit to be reached.
			select {
			case <-h.batch.sendSchedule:
				// Batch has been sent manually.
			case <-time.After(5 * time.Second):
				// Time limit has been reached.
				_ = h.batch.flush()
			}
			h.batch = nil
		}()
	}

	// Add the log entry to the batch.
	h.batch.entries = append(h.batch.entries, entry)

	// If the batch is full, flush it and create a new batch.
	if len(h.batch.entries) >= h.batch.maxEntries {
		err := h.batch.flush()
		if err != nil {
			return err
		}
		h.batch = nil
	}

	return nil
}

type lokiClient struct {
	rh *sioUtils.RestHelpers
}

type LokiClient interface{}

func newLokiClient() *lokiClient {
	return &lokiClient{
		rh: sioUtils.NewRestHelpers(),
	}
}

func (lc *lokiClient) SendLog(request *LokiRequest) error {
	rh := sioUtils.NewRestHelpers()

	url := fmt.Sprintf("%s/loki/api/v1/push", os.Getenv("LOKI_URL"))

	rJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}
	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)
	req.Header.Add("Content-Type", "application/json")

	_, err = rh.DoHttpRequest(req)
	if err != nil {
		return err
	}

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
	temp := LokiRequest{
		Streams: []LokiStream{
			{
				Stream: map[string]string{"t": "a"},
				Values: [][]string{
					{strconv.FormatInt(time.Now().UnixNano(), 10), "test"},
				},
			},
		},
	}

	lc := newLokiClient()
	err := lc.SendLog(&temp)
	if err != nil {
		log.Fatal(err)
	}
	// return response, nil

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
	// log.AddHook(lh)
}

func main() {
	CreateRouter()
}
