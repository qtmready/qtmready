package slack

import (
	"net/http"
)

// HTTPClient implements the httpClient interface.
type (
	HTTPClient struct{}
)

// Do executes an HTTP request using the custom HTTP client.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}
