package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/stretchr/testify/assert"
)

var awUser = &siogeneric.AwUser{
	ID:    "10000069",
	Phone: "5555555555",
	Email: "test-integration@slauson.io",
	Name:  "Iam Integration",
}

func runTestServer(t *testing.T) (*httptest.Server, string) {
	s := httptest.NewServer(CreateRouter())
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}
	return s, token.AccessToken
}

func TestCreateUser_HappyScenarios(t *testing.T) {
	ts, token := runTestServer(t)
	defer ts.Close()

	tests := []struct {
		name    string
		request *siogeneric.AwCreateUserRequest
	}{
		{
			name: "happy",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rJSON, err := json.Marshal(test.request)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			// Test
			req, err := http.NewRequest("POST", ts.URL+"/api/iam/v1/user", sr)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
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
		})
	}
}
