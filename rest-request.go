package sdk

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// DefaultHTTPClient initialized Grafana with appropriate conditions.
// It allows you globally redefine HTTP client.
var DefaultHTTPClient = http.DefaultClient

// Client uses Grafana REST API for interacting with Grafana server.
type Client struct {
	baseURL       string
	key           string
	basicAuth     bool
	client        *http.Client
	customHeaders map[string]string
}

// SetCustomHeaders - set additional headers that will be sent with each request
func (r *Client) SetCustomHeaders(headers map[string]string) {
	r.customHeaders = headers
}

// SetCustomHeader - set additional header that will be sent with each request
func (r *Client) SetCustomHeader(key, value string) {
	if r.customHeaders == nil {
		r.customHeaders = map[string]string{}
	}
	r.customHeaders[key] = value
}

// SetOrgIDHeader - set `X-Grafana-Org-Id`  to specified value
func (r *Client) SetOrgIDHeader(oid uint) {
	r.SetCustomHeader("X-Grafana-Org-Id", strconv.FormatUint(uint64(oid), 10))
}

func (r *Client) WithOrgIDHeader(oid uint) *Client {
	c := &Client{
		baseURL:   r.baseURL,
		key:       r.key,
		basicAuth: r.basicAuth,
		client:    r.client,
	}
	c.SetOrgIDHeader(oid)
	return c
}

// StatusMessage reflects status message as it returned by Grafana REST API.
type StatusMessage struct {
	ID      *uint   `json:"id"`
	OrgID   *uint   `json:"orgId"`
	Message *string `json:"message"`
	Slug    *string `json:"slug"`
	Version *int    `json:"version"`
	Status  *string `json:"status"`
	UID     *string `json:"uid"`
	URL     *string `json:"url"`
}

// NewClient initializes client for interacting with an instance of Grafana server;
// apiKeyOrBasicAuth accepts either 'username:password' basic authentication credentials,
// or a Grafana API key
func NewClient(apiURL, apiKeyOrBasicAuth string, client *http.Client) (*Client, error) {
	key := ""
	basicAuth := strings.Contains(apiKeyOrBasicAuth, ":")
	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	if !basicAuth {
		key = fmt.Sprintf("Bearer %s", apiKeyOrBasicAuth)
	} else {
		parts := strings.Split(apiKeyOrBasicAuth, ":")
		baseURL.User = url.UserPassword(parts[0], parts[1])
	}
	return &Client{baseURL: baseURL.String(), basicAuth: basicAuth, key: key, client: client}, nil
}

func (r *Client) get(ctx context.Context, query string, params url.Values) ([]byte, int, error) {
	return r.doRequest(ctx, "GET", query, "", params, nil)
}

func (r *Client) getWithRawPath(ctx context.Context, query, rawPath string, params url.Values) ([]byte, int, error) {
	return r.doRequest(ctx, "GET", query, rawPath, params, nil)
}

func (r *Client) patch(ctx context.Context, query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest(ctx, "PATCH", query, "", params, bytes.NewBuffer(body))
}

func (r *Client) put(ctx context.Context, query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest(ctx, "PUT", query, "", params, bytes.NewBuffer(body))
}

func (r *Client) post(ctx context.Context, query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest(ctx, "POST", query, "", params, bytes.NewBuffer(body))
}

func (r *Client) delete(ctx context.Context, query string) ([]byte, int, error) {
	return r.doRequest(ctx, "DELETE", query, "", nil, nil)
}

func (r *Client) doRequest(ctx context.Context, method, query, rawPath string, params url.Values, buf io.Reader) ([]byte, int, error) {
	u, _ := url.Parse(r.baseURL)
	u.Path = path.Join(u.Path, query)
	if rawPath != "" {
		u.RawPath = path.Join(u.RawPath, rawPath)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, 0, err
	}
	req = req.WithContext(ctx)
	if !r.basicAuth {
		req.Header.Set("Authorization", r.key)
	}
	if r.customHeaders != nil {
		for key, value := range r.customHeaders {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "autograf")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, resp.StatusCode, err
}
