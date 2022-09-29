// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

// client provides methods to interact with the ctrlplane API.
// the long term goal is to make this available as an SDK.
//
// The main client.go provides the client.New() method to create a new client along with the call() method to be
// used internally by the client. All of the exported methods are wrappers around the call() method.
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.breu.io/ctrlplane/internal/api/auth"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	Client struct {
		BaseURL string
	}
)

// New returns a new client.
func New() *Client {
	return &Client{
		BaseURL: shared.Service.CLI.BaseURL,
	}
}

// call is a helper function to call ctrlplane REST API. Example on how to use it. The long term goal is to make this
// available as an SDK.
//
//	import "go.breu.io/ctrlplane/internal/client"
//	import "go.breu.io/ctrlplane/internal/api/core"
//
//	c := client.New()
//	url := "/core/login"
//	data := &Data{}
//	reply := &Reply{}
//
//	err := c.call("GET", url, reply, data)
func (c *Client) call(method, url string, reply, data interface{}) error {
	var (
		err     error
		request *http.Request
	)

	url = c.url(url)

	if data != nil {
		data, _ = json.Marshal(data)
		request, err = http.NewRequest(method, url, bytes.NewReader(data.([]byte)))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return err
	}

	c.headers(request)

	httpclient := &http.Client{}
	response, err := httpclient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		body, _ := io.ReadAll(response.Body)

		if err = json.Unmarshal(body, &reply); err != nil {
			return err
		}

		return nil
	default:
		return ErrInvalidCredentials
	}
}

// headers sets the headers for the request.
func (c *Client) headers(request *http.Request) {
	request.Header.Set("User-Agent", shared.Service.Name+"/"+shared.Service.Version())
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", auth.APIKeyPrefix+" "+shared.Service.CLI.APIKEY)
}

// url returns the full URL for the request.
func (c *Client) url(path string) string {
	return c.BaseURL + path
}
