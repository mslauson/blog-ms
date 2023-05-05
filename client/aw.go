package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
	"gitea.slauson.io/slausonio/iam-ms/constants"
)

type AwClient struct {
	h              sioUtils.SioRestHelpers
	defaultHeaders map[string][]string
	host           string
	key            string
}

//go:generate mockery --name AppwriteClient
type AppwriteClient interface {
	ListUsers() (*siogeneric.AwlistResponse, error)
	GetUserByID(id string) (*siogeneric.AwUser, error)
	CreateUser(r *siogeneric.AwCreateUserRequest) (*siogeneric.AwUser, error)
	UpdateEmail(id string, r *siogeneric.UpdateEmailRequest) (*siogeneric.AwUser, error)
	UpdatePhone(id string, r *siogeneric.UpdatePhoneRequest) (*siogeneric.AwUser, error)
	UpdatePassword(id string, r *siogeneric.UpdatePasswordRequest) (*siogeneric.AwUser, error)
	DeleteUser(id string) error
	CreateEmailSession(r *siogeneric.AwEmailSessionRequest) (*siogeneric.AwSession, error)
	DeleteSession(ID, sID string) error
}

func NewAwClient() *AwClient {
	return &AwClient{
		h: sioUtils.NewRestHelpers(),
		defaultHeaders: map[string][]string{
			"Content-Type":                 {"application/json"},
			constants.AW_HEADER_PROJECT_ID: {os.Getenv("IAM_PROJECT")},
		},
		host: os.Getenv("IAM_HOST"),
		key:  os.Getenv("IAM_KEY"),
	}
}

func (c *AwClient) ListUsers() (*siogeneric.AwlistResponse, error) {
	url := fmt.Sprintf("%s/users", c.host)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)
	response := new(siogeneric.AwlistResponse)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) GetUserByID(id string) (*siogeneric.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s", c.host, id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	response := new(siogeneric.AwUser)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) CreateUser(r *siogeneric.AwCreateUserRequest) (*siogeneric.AwUser, error) {
	url := fmt.Sprintf("%s/users", c.host)
	r.Phone = fmt.Sprintf("+1%s", r.Phone)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	response := new(siogeneric.AwUser)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) UpdateEmail(
	id string,
	r *siogeneric.UpdateEmailRequest,
) (*siogeneric.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/email", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	response := new(siogeneric.AwUser)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) UpdatePassword(
	id string,
	r *siogeneric.UpdatePasswordRequest,
) (*siogeneric.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/password", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	response := new(siogeneric.AwUser)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) UpdatePhone(
	id string,
	r *siogeneric.UpdatePhoneRequest,
) (*siogeneric.AwUser, error) {
	url := fmt.Sprintf("%s/users/%s/phone", c.host, id)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("PATCH", url, sr)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	response := new(siogeneric.AwUser)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) DeleteUser(id string) error {
	url := fmt.Sprintf("%s/users/%s", c.host, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	return c.executeAndParseResponse(req, nil)
}

func (c *AwClient) CreateEmailSession(
	r *siogeneric.AwEmailSessionRequest,
) (*siogeneric.AwSession, error) {
	url := fmt.Sprintf("%s/account/sessions/email", c.host)
	rJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(rJSON))
	req, _ := http.NewRequest("POST", url, sr)

	req.Header = c.defaultHeaders

	response := new(siogeneric.AwSession)
	if err := c.executeAndParseResponse(req, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *AwClient) DeleteSession(ID, sID string) error {
	url := fmt.Sprintf("%s/users/%s/sessions/%s", c.host, ID, sID)
	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header = c.defaultHeaders
	req.Header.Add(constants.AW_HEADER_KEY, c.key)

	return c.executeAndParseResponse(req, nil)
}

func (c *AwClient) executeAndParseResponse(
	req *http.Request,
	response any,
) error {
	res, err := c.h.ExecuteRequest(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 300) {
		errRes := new(siogeneric.AppwriteError)
		if err := c.h.ParseResponse(res, errRes); err != nil {
			return err
		}

		return sioerror.NewSioIamError(errRes)
	} else if response != nil {
		if err := c.h.ParseResponse(res, response); err != nil {
			return err
		}
	}

	return nil
}
