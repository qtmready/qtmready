package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/shared"
)

type (
	Client struct{}
)

func New() *Client { return &Client{} }

func (c *Client) UserAgent() string {
	version := shared.Service.Name + "/" + shared.Service.Version()
	return version
}

func (c *Client) Request(method, url string, data interface{}, reply interface{}) error {
	body, _ := json.Marshal(data)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", c.UserAgent())
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &reply)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid credentials")
	}
}

func (c *Client) SetAuthenticationHeaders(request *http.Request) {
	//
}
