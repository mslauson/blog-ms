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
var createdUsers []*siogeneric.AwUser

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

func TestCreateUser_Errors(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	tests := []struct {
		name       string
		request    *siogeneric.AwCreateUserRequest
		error      string
		statusCode int
	}{
		{
			name: "already exists",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "A user with the same email already exists in your project.",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "missing user id",
			request: &siogeneric.AwCreateUserRequest{
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "Key: 'AwCreateUserRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "missing phone",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "Key: 'AwCreateUserRequest.Phone' Error:Field validation for 'Phone' failed on the 'required' tag",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "bad phone short",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "d64f9807ee288ab37faf7150edf3eb08",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "please enter a ten digit mobile number",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "bad phone long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "9cab364d79ce0219b991a47678529b97084a12a88cd1b2e1aed9dd964476789f",
				Phone:    "d64f9807ee288ab37faf7150edf3eb08",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "please enter a ten digit mobile number",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "bad email",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "f418e9a1f5f0b1b052e7c2ba97cb6b96",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "invalid email",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Email too long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "184be7fa35e55150a59032b480fdd0ce015c67c375a1bb08a46f35dd178a944c0c2d5157a95b68ce6b339457c87c0c0be4e9b35039327fef8fee586773fe76c2b395c7496ced1aa20dfec178562dc1504cb154a5eeffbf876e7d775c8703754a9257499083d7e9c89d8987a6c81b9f112e49050ba018dd8e967a68ab23582d9c31d56df3ba47d15cc41dd1104d73d9253b2a4d03b3c48510a24706c5e6cf1b9037e25de2f2f4f8df92fa2802f37e661caea499a20c3295b431bdb15d7815ae417ca2bbc963fd7752ec47a4c9eb43894e1287ef18e0f73619405219f3d3c936a445bbd2a6b8a5ee342da211ee6fa65dd52bb41b5240992d745ad0ad1a79dd4f21fc673b22464b1772d0146f30fd220b98a9a7bc26862d8daa5a9d6fb09afa11cc672558151bb024cda12f7c9fdd18941480b5d139290236368bafbbe4986341565a1a0f090c9ebde17e8a443a0bfdca94",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "invalid email",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "no name",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "Key: 'AwCreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "name too long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "6f124736c6b70801fde7273624f7bb9d",
				Email:    "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:     "184be7fa35e55150a59032b480fdd0ce015c67c375a1bb08a46f35dd178a944c0c2d5157a95b68ce6b339457c87c0c0be4e9b35039327fef8fee586773fe76c2b395c7496ced1aa20dfec178562dc1504cb154a5eeffbf876e7d775c8703754a9257499083d7e9c89d8987a6c81b9f112e49050ba018dd8e967a68ab23582d9c31d56df3ba47d15cc41dd1104d73d9253b2a4d03b3c48510a24706c5e6cf1b9037e25de2f2f4f8df92fa2802f37e661caea499a20c3295b431bdb15d7815ae417ca2bbc963fd7752ec47a4c9eb43894e1287ef18e0f73619405219f3d3c936a445bbd2a6b8a5ee342da211ee6fa65dd52bb41b5240992d745ad0ad1a79dd4f21fc673b22464b1772d0146f30fd220b98a9a7bc26862d8daa5a9d6fb09afa11cc672558151bb024cda12f7c9fdd18941480b5d139290236368bafbbe498634156c5febfc8d65269da74ae6a5e2e57ee66",
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "invalid name",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "no password",
			request: &siogeneric.AwCreateUserRequest{
				UserID: "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:  "6f124736c6b70801fde7273624f7bb9d",
				Email:  "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
				Name:   "3e73bbe6b2605b01a3456022cc30688b",
			},
			error:      "Key: 'AwCreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Password Missing Number",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "97fd79d960696af638433104adb9b255",
				Email:    "efd623543c981069ca0c05c57b59d69bf97f684cd69a15f4cef7c8260447c82c",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "f83aac7a772ed502d2cecd4d1c91d900",
			},
			statusCode: http.StatusBadRequest,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Bad Password Missing Upper",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "97fd79d960696af638433104adb9b255",
				Email:    "efd623543c981069ca0c05c57b59d69bf97f684cd69a15f4cef7c8260447c82c",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "47422a14ff04a471b6829843fde489ae",
			},
			statusCode: http.StatusBadRequest,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Bad Password Missing Special",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "97fd79d960696af638433104adb9b255",
				Email:    "efd623543c981069ca0c05c57b59d69bf97f684cd69a15f4cef7c8260447c82c",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "c7510a622ca0a4f52576023c0ff7c7a6",
			},
			statusCode: http.StatusBadRequest,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Short Password",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "5ba9ac5847ee1d04c2dcce97a377ee3d",
				Phone:    "97fd79d960696af638433104adb9b255",
				Email:    "efd623543c981069ca0c05c57b59d69bf97f684cd69a15f4cef7c8260447c82c",
				Name:     "3e73bbe6b2605b01a3456022cc30688b",
				Password: "43dad3e484522e9252e30db80557d1d4",
			},
			statusCode: http.StatusBadRequest,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
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

			siotest.ParseCheckIfCorrectError(t, resp, test.error, test.statusCode)
		})
	}
}

func TestUpdateEmail_HappyScenarios(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	tests := []struct {
		name    string
		request *siogeneric.UpdateEmailRequest
	}{
		{
			name: "happy",
			request: &siogeneric.UpdateEmailRequest{
				Email: "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
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
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/email", ts.URL, createdUsers[0].ID),
				sr,
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

			// TODO: change when encrypt response
			// assert.Equal(t, awUser.ID, result.ID)
			// assert.Equal(t, awUser.Phone, result.Phone)
			// assert.Equal(t, awUser.Email, result.Email)
			// assert.Equal(t, awUser.Name, result.Name)

		})
	}
}

func TestUpdateEmail_Errors(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	id := createdUsers[0].ID

	tests := []struct {
		name       string
		request    *siogeneric.UpdateEmailRequest
		error      string
		statusCode int
		id         string
	}{
		{
			name: "not found",
			request: &siogeneric.UpdateEmailRequest{
				Email: "629ab286599f7b5b67ee1d88093b0280608ed13e9083b7ea883cba59c5330860",
			},
			error:      "User with the requested ID could not be found.",
			statusCode: http.StatusNotFound,
			id:         "123",
		},
		{
			name: "bad email",
			request: &siogeneric.UpdateEmailRequest{
				Email: "f418e9a1f5f0b1b052e7c2ba97cb6b96",
			},
			error:      "invalid email",
			statusCode: http.StatusBadRequest,
			id:         id,
		},
		{
			name: "Email too long",
			request: &siogeneric.UpdateEmailRequest{
				Email: "184be7fa35e55150a59032b480fdd0ce015c67c375a1bb08a46f35dd178a944c0c2d5157a95b68ce6b339457c87c0c0be4e9b35039327fef8fee586773fe76c2b395c7496ced1aa20dfec178562dc1504cb154a5eeffbf876e7d775c8703754a9257499083d7e9c89d8987a6c81b9f112e49050ba018dd8e967a68ab23582d9c31d56df3ba47d15cc41dd1104d73d9253b2a4d03b3c48510a24706c5e6cf1b9037e25de2f2f4f8df92fa2802f37e661caea499a20c3295b431bdb15d7815ae417ca2bbc963fd7752ec47a4c9eb43894e1287ef18e0f73619405219f3d3c936a445bbd2a6b8a5ee342da211ee6fa65dd52bb41b5240992d745ad0ad1a79dd4f21fc673b22464b1772d0146f30fd220b98a9a7bc26862d8daa5a9d6fb09afa11cc672558151bb024cda12f7c9fdd18941480b5d139290236368bafbbe4986341565a1a0f090c9ebde17e8a443a0bfdca94",
			},
			error:      "invalid email",
			statusCode: http.StatusBadRequest,
			id:         id,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rJSON, err := json.Marshal(test.request)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/email", ts.URL, test.id),
				sr,
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

			siotest.ParseCheckIfCorrectError(t, resp, test.error, test.statusCode)
		})
	}
}

func TestUpdatePassword_HappyScenarios(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	tests := []struct {
		name    string
		request *siogeneric.UpdatePasswordRequest
	}{
		{
			name: "happy",
			request: &siogeneric.UpdatePasswordRequest{
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
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/password", ts.URL, createdUsers[0].ID),
				sr,
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
		})
	}
}

func TestUpdatePassword_Errors(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	id := createdUsers[0].ID

	tests := []struct {
		name       string
		request    *siogeneric.UpdatePasswordRequest
		error      string
		statusCode int
		id         string
	}{
		{
			name: "not found",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "c73583a948d7662f30828d764834552b",
			},
			error:      "User with the requested ID could not be found.",
			statusCode: http.StatusNotFound,
			id:         "123",
		},
		{
			name:       "no password",
			request:    &siogeneric.UpdatePasswordRequest{},
			error:      "Key: 'UpdatePasswordRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			statusCode: http.StatusBadRequest,
			id:         id,
		},
		{
			name: "Bad Password Missing Number",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "f83aac7a772ed502d2cecd4d1c91d900",
			},
			statusCode: http.StatusBadRequest,
			id:         id,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Bad Password Missing Upper",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "efd623543c981069ca0c05c57b59d69bf97f684cd69a15f4cef7c8260447c82c",
			},
			statusCode: http.StatusBadRequest,
			id:         id,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Bad Password Missing Special",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "c7510a622ca0a4f52576023c0ff7c7a6",
			},
			statusCode: http.StatusBadRequest,
			id:         id,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
		{
			name: "Short Password",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "43dad3e484522e9252e30db80557d1d4",
			},
			statusCode: http.StatusBadRequest,
			id:         id,
			error:      "invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rJSON, err := json.Marshal(test.request)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/password", ts.URL, test.id),
				sr,
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

			siotest.ParseCheckIfCorrectError(t, resp, test.error, test.statusCode)
		})
	}
}

func TestUpdatePhone_HappyScenarios(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	tests := []struct {
		name    string
		request *siogeneric.UpdatePhoneRequest
	}{
		{
			name: "happy",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "6f124736c6b70801fde7273624f7bb9d",
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
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/phone", ts.URL, createdUsers[0].ID),
				sr,
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
		})
	}
}

func TestUpdatePhone_Errors(t *testing.T) {
	ts, token := siotest.RunTestServer(t, CreateRouter())
	defer ts.Close()

	id := createdUsers[0].ID

	tests := []struct {
		name       string
		request    *siogeneric.UpdatePhoneRequest
		error      string
		statusCode int
		id         string
	}{
		{
			name: "not found",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "6f124736c6b70801fde7273624f7bb9d",
			},
			error:      "User with the requested ID could not be found.",
			statusCode: http.StatusNotFound,
			id:         "123",
		},
		{
			name: "invalid short",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "d64f9807ee288ab37faf7150edf3eb08",
			},
			error:      "please enter a ten digit mobile number",
			statusCode: http.StatusBadRequest,
			id:         id,
		},
		{
			name: "invalid long",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "d64f9807ee288ab37faf7150edf3eb08",
			},
			error:      "please enter a ten digit mobile number",
			statusCode: http.StatusBadRequest,
			id:         id,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rJSON, err := json.Marshal(test.request)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PUT",
				fmt.Sprintf("%s/api/iam/v1/user/%s/phone", ts.URL, test.id),
				sr,
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

			siotest.ParseCheckIfCorrectError(t, resp, test.error, test.statusCode)
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

			siotest.ParseCheckIfCorrectError(
				t,
				resp,
				"User with the requested ID could not be found.",
				http.StatusNotFound,
			)
		})
	}
}
