package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	sioRest "gitea.slauson.io/slausonio/go-utils/rest"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
)

type AwClient struct {
	h              sioRest.RestHelpers
	defaultHeaders map[string][]string
	host           string
	key            string
}

type AppwriteClient interface {
	ListUsers() (*sioModel.AwlistResponse, error)
	GetUserByID(id string) (*sioModel.AwUser, error)
	CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error)
	UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error)
	UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error)
	UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error)

	CreateSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error)
	DeleteSession(sID string) error
}

func NewAwClient() *AwClient {
	return &AwClient{
		h: sioRest.RestHelpers{},
		defaultHeaders: map[string][]string{
			"Content-Type":       {"application/json"},
			"X-Appwrite-Project": {os.Getenv("IAM_PROJECT")},
		},
		host: os.Getenv("IAM_HOST"),
		key:  os.Getenv("IAM_KEY"),
	}
}

func (c *AwClient) ListUsers() (*sioModel.AwlistResponse, error) {
	url := fmt.Sprintf("%s/users", c.host)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)
	response := new(sioModel.AwlistResponse)
	err := c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) GetUserByID(id string) (*sioModel.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s", c.host, id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	response := new(sioModel.AwUser)
	err := c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error) {
	url := fmt.Sprintf("%s/users", c.host)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	response := new(sioModel.AwUser)
	err = c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/email", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	response := new(sioModel.AwUser)
	err = c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/password", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	response := new(sioModel.AwUser)
	err = c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/phone", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	response := new(sioModel.AwUser)
	err = c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) DeleteUser(id string) error {
	url := fmt.Sprintf("%s/users/%s", c.host, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header = c.defaultHeaders
	req.Header.Add("X-Appwrite-Key", c.key)

	_, err := c.h.DoHttpRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *AwClient) CreateSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error) {
	url := fmt.Sprintf("%s/account/sessions/email", c.host)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)

	req.Header = c.defaultHeaders

	response := new(sioModel.AwSession)
	err = c.h.DoHttpRequestAndParse(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AwClient) DeleteSession(sID string) error {
	url := fmt.Sprintf("%s/account/sessions/%s", c.host, sID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header = c.defaultHeaders

	_, err := c.h.DoHttpRequest(req)
	if err != nil {
		return err
	}

	return nil
}
