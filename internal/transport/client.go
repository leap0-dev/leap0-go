package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/leap0dev/leap0-go/internal/sse"
)

type Config struct {
	BaseURL       string
	SandboxDomain string
	APIKey        string
	AuthHeader    string
	Bearer        bool
	Source        string
	Version       string
	HTTPClient    *http.Client
}

type Options struct {
	Headers        map[string]string
	Query          map[string]string
	ExpectedStatus []int
}

type APIError struct {
	Status  int
	Headers http.Header
	Body    string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("leap0 api: status %d: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("leap0 api: status %d: %s", e.Status, http.StatusText(e.Status))
}

type Client struct{ cfg *Config }

func New(cfg *Config) *Client { return &Client{cfg: cfg} }

func (c *Client) SandboxURL(sandboxID string, port int, path string) string {
	host := sandboxID + "." + c.cfg.SandboxDomain
	if port > 0 {
		host = fmt.Sprintf("%s-%d.%s", sandboxID, port, c.cfg.SandboxDomain)
	}
	return "https://" + host + path
}

func (c *Client) SandboxWSURL(sandboxID string, port int, path string) string {
	host := sandboxID + "." + c.cfg.SandboxDomain
	if port > 0 {
		host = fmt.Sprintf("%s-%d.%s", sandboxID, port, c.cfg.SandboxDomain)
	}
	return "wss://" + host + path
}

func (c *Client) APIKey() string { return c.cfg.APIKey }

func (c *Client) authValue() string {
	if c.cfg.Bearer {
		return "Bearer " + c.cfg.APIKey
	}
	return c.cfg.APIKey
}

func (c *Client) buildURL(path string, query map[string]string) string {
	return c.applyQuery(strings.TrimRight(c.cfg.BaseURL, "/")+path, query)
}

func (c *Client) applyQuery(u string, query map[string]string) string {
	if len(query) == 0 {
		return u
	}
	params := url.Values{}
	for k, v := range query {
		if v != "" {
			params.Set(k, v)
		}
	}
	if enc := params.Encode(); enc != "" {
		return u + "?" + enc
	}
	return u
}

func (c *Client) do(ctx context.Context, method, url string, body io.Reader, opts *Options) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(c.cfg.AuthHeader, c.authValue())
	req.Header.Set("Leap0-Source", c.cfg.Source)
	req.Header.Set("Leap0-SDK-Version", c.cfg.Version)
	req.Header.Set("User-Agent", "leap0-go/"+c.cfg.Version)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request timed out")
		}
		return nil, err
	}
	return resp, nil
}

func (c *Client) checkStatus(resp *http.Response, opts *Options) error {
	expected := []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}
	if opts != nil && len(opts.ExpectedStatus) > 0 {
		expected = opts.ExpectedStatus
	}
	for _, code := range expected {
		if resp.StatusCode == code {
			return nil
		}
	}
	body, _ := io.ReadAll(resp.Body)
	return &APIError{Status: resp.StatusCode, Headers: resp.Header, Body: string(body)}
}

func (c *Client) JSON(ctx context.Context, method, path string, payload, result any, opts *Options) error {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.buildURL(path, query), body, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := c.checkStatus(resp, opts); err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) JSONAbsolute(ctx context.Context, method, absURL string, payload, result any, opts *Options) error {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.applyQuery(absURL, query), body, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := c.checkStatus(resp, opts); err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) Text(ctx context.Context, method, path string, payload any, opts *Options) (string, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.buildURL(path, query), body, opts)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := c.checkStatus(resp, opts); err != nil {
		return "", err
	}
	raw, err := io.ReadAll(resp.Body)
	return string(raw), err
}

func (c *Client) Bytes(ctx context.Context, method, path string, payload any, opts *Options) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.buildURL(path, query), body, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := c.checkStatus(resp, opts); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (c *Client) BytesAbsolute(ctx context.Context, method, absURL string, payload any, opts *Options) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.applyQuery(absURL, query), body, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := c.checkStatus(resp, opts); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (c *Client) SSE(ctx context.Context, method, path string, payload any, opts *Options) (*sse.Stream, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.buildURL(path, query), body, opts)
	if err != nil {
		return nil, err
	}
	if err := c.checkStatus(resp, opts); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return sse.New(resp.Body), nil
}

func (c *Client) SSEAbsolute(ctx context.Context, method, absURL string, payload any, opts *Options) (*sse.Stream, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	var query map[string]string
	if opts != nil {
		query = opts.Query
	}
	resp, err := c.do(ctx, method, c.applyQuery(absURL, query), body, opts)
	if err != nil {
		return nil, err
	}
	if err := c.checkStatus(resp, opts); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return sse.New(resp.Body), nil
}

func (c *Client) Raw(ctx context.Context, method, url string, body io.Reader, opts *Options) (*http.Response, error) {
	return c.do(ctx, method, url, body, opts)
}
