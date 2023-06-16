package appstoreapi

import (
	"net/http"
	"time"

	appleapigoclient "github.com/gh73962/appleapis"
	"github.com/gh73962/appleapis/appstore/api/internal/httputils"
	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

type Option func(*ClientOption)

type ClientOption struct {
	NeedRetry    bool          // retry use backoff and jitter
	RetryInitial time.Duration // retry first retry pause duration , default 100ms
	RetryMax     time.Duration // retry max duration, default 30s
	HTTPClient   *http.Client  // default use http.DefaultClient
	UserAgent    string        // default apple-api-go-client
	IsSandbox    bool          // default false
}

func (c *ClientOption) GetBackoff() *httputils.BackoffImpl {
	if c.RetryInitial == 0 {
		return nil
	}

	return &httputils.BackoffImpl{
		Initial: c.RetryInitial,
		Max:     c.RetryMax,
	}
}

func (c *ClientOption) GetBasePath() string {
	if c.IsSandbox {
		return datatypes.SandboxBasePath
	}
	return datatypes.BasePath
}

func (c *ClientOption) GetUserAgent() string {
	if c.UserAgent == "" {
		return appleapigoclient.UserAgent
	}

	return c.UserAgent
}

func WithRetry(initial, max time.Duration) Option {
	return func(c *ClientOption) {
		c.NeedRetry = true
		if initial == 0 {
			c.RetryInitial = 100 * time.Millisecond
		}
		if max == 0 {
			c.RetryMax = 10 * time.Second
		}
	}
}

func WithSandbox() Option {
	return func(c *ClientOption) {
		c.IsSandbox = true
	}
}

func WithUserAgent(data string) Option {
	return func(c *ClientOption) {
		c.UserAgent = data
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *ClientOption) {
		c.HTTPClient = client
	}
}
