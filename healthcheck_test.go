package main

import (
	"net/http"
	"testing"

	siotest "gitea.slauson.io/slausonio/go-testing/sio_test"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoot(t *testing.T) {
	ts, _ := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	assert.Equalf(
		t,
		resp.StatusCode,
		http.StatusOK,
		"expected status code %d, got %d",
		http.StatusOK,
		resp.StatusCode,
	)
}

func TestHealthCheckContextPath(t *testing.T) {
	ts, _ := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/iam", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	assert.Equalf(
		t,
		resp.StatusCode,
		http.StatusOK,
		"expected status code %d, got %d",
		http.StatusOK,
		resp.StatusCode,
	)
}
