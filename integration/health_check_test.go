package integration

import (
	"net/http"
	"testing"

	"gitea.slauson.io/blog/post-ms/handler"
	"gitea.slauson.io/slausonio/go-testing/siotest"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoot(t *testing.T) {
	ts, _ := siotest.RunTestServer(t, handler.CreateRouter())
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
	ts, _ := siotest.RunTestServer(t, handler.CreateRouter())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/post", nil)
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
