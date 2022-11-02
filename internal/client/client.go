// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

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
		APIKEY  string
	}
)

// New returns a new client.
func New(url string, key string) *Client {
	return &Client{
		BaseURL: url,
		APIKEY:  key,
	}
}

// request is a helper function to request ctrlplane REST API. Example on how to use it. The long term goal is to make this
// available as an SDK.
//
//	import "go.breu.io/ctrlplane/internal/client"
//	import "go.breu.io/ctrlplane/internal/api/core"
//
//	c := client.New("https://api.ctrlplane.ai", <api key>)
//	url := "/core/login"
//	data := &Data{}
//	reply := &Reply{}
//
//	err := c.request("GET", url, reply, data)
func (c *Client) request(method, url string, reply, data interface{}) error {
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
	request.Header.Set("Authorization", auth.APIKeyPrefix+" "+c.APIKEY)
}

// url returns the full URL for the request.
func (c *Client) url(path string) string {
	return c.BaseURL + path
}
