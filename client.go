package getmac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ListResponse struct {
	Total int `json:"total"`
}

type ClientOpts func(*Client)

func WithBaseURL(url string) ClientOpts {
	f := func(c *Client) {
		c.baseURL = url
	}

	return f
}

func WithToken(token string) ClientOpts {
	f := func(c *Client) {
		c.token = token
	}

	return f
}

func WithHTTPClient(httpClient *http.Client) ClientOpts {
	f := func(c *Client) {
		c.httpClient = httpClient
	}

	return f
}

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(opts ...ClientOpts) *Client {
	c := &Client{
		baseURL:    "https://api.getmac.io/v1",
		token:      "",
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.baseURL, path), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (c *Client) VirtualMachines() *virtualMachinesService {
	return &virtualMachinesService{client: c}
}
