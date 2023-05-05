package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	siotest "gitea.slauson.io/slausonio/go-testing/sio_test"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"github.com/stretchr/testify/assert"
)

var awUser = &siogeneric.AwUser{
	ID:    "10000069",
	Phone: "+15555555555",
	Email: "iam-integration@slauson.io",
	Name:  "Iam Integration",
}
var createdUsers = []*siogeneric.AwUser{}

func TestCreateUser_HappyScenarios(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
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
			result := &siogeneric.AwUser{}
			siotest.ParseHappyResponse(t, resp, result)

			assert.Equal(t, awUser.ID, result.ID)
			assert.Equal(t, awUser.Phone, result.Phone)
			assert.Equal(t, awUser.Email, result.Email)
			assert.Equal(t, awUser.Name, result.Name)

			createdUsers = append(createdUsers, result)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	for _, user := range createdUsers {
		t.Run(fmt.Sprintf("Delete user.  ID: %s", user.ID), func(t *testing.T) {
			req, err := http.NewRequest(
				"DELETE",
				fmt.Sprintf("%s%s%s", ts.URL, "/api/iam/v1/user/", user.ID),
				nil,
			)
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

			result := &siogeneric.SuccessResponse{}
			siotest.ParseHappyResponse(t, resp, result)

			assert.Truef(t, result.Success, "expected success to be true, got false")
		})
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	for _, user := range createdUsers {
		t.Run(fmt.Sprintf("Delete user.  ID: %s", user.ID), func(t *testing.T) {
			req, err := http.NewRequest(
				"DELETE",
				fmt.Sprintf("%s%s%s", ts.URL, "/api/iam/v1/user/", user.ID),
				nil,
			)
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

			siotest.ParseCheckIfCorrectError(t, resp, "user not found", http.StatusNotFound)
		})
	}
}
