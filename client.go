package bclient

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*http.Client
	Headers map[string]string
}

func (c *Client) Do(req *http.Request, headers ...map[string]string) (*http.Response, []byte, error) {
	setHeader(req, c.Headers)
	setHeader(req, headers...)

	res, err := c.Client.Do(req)
	if err != nil {
		return res, nil, err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return res, nil, err
		}
		defer reader.Close()
	default:
		reader = res.Body
	}

	data, err := io.ReadAll(reader)
	return res, data, err
}

func (c *Client) Get(ctx context.Context, url string, headers ...map[string]string) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	return c.Do(req, headers...)
}

func (c *Client) PostForm(ctx context.Context, url string, data url.Values, headers ...map[string]string) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req, headers...)
}

func setHeader(req *http.Request, headers ...map[string]string) {
	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}
}
