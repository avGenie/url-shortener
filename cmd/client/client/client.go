// Package client implements client for shortener
package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Client Contains information for sending request to client
type Client struct {
	Client   *http.Client
	postAddr string
}

// NewClient Creates client object
func NewClient(postAddr string) *Client {
	return &Client{
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		postAddr: postAddr,
	}
}

// SendPostRequest Sends POST request
func (c *Client) SendPostRequest(data []byte, cookie *http.Cookie) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, c.postAddr, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendGetRequest Sends GET request
func (c *Client) SendGetRequest(netAddr, baseURIPrefix string) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s/%s", baseURIPrefix, netAddr)
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, nil
}

// ReadBody Reads response body
func (c *Client) ReadBody(response *http.Response) ([]byte, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
