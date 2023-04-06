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
	log "github.com/sirupsen/logrus"
)

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

	rh := sioUtils.NewRestHelpers()
	url := fmt.Sprintf("%s/loki/api/v1/push", os.Getenv("LOKI_URL"))

	rJSON, err := json.Marshal(temp)
	if err != nil {
		log.Error(err)
		// return nil, err
	}
	log.Println(temp)

	sr := strings.NewReader(string(rJSON))
	println(sr)
	req, _ := http.NewRequest("POST", url, sr)
	req.Header.Add("Content-Type", "application/json")

	res, err := rh.DoHttpRequest(req)
	if err != nil {
		// return nil, err
		log.Error(err)
	}

	log.Println(res.StatusCode)

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
