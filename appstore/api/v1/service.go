package appstoreapi

import (
	"context"
	"net/http"
	"time"
)

// Backoff see https://aws.amazon.com/cn/blogs/architecture/exponential-backoff-and-jitter/
type Backoff interface {
	Pause() time.Duration
}

type Service struct {
	client    *http.Client
	BasePath  string
	UserAgent string
	BackOff   Backoff
	NeedRetry bool
}

func NewAppStoreService(ctx context.Context, options ...Option) *Service {
	var clientOpt ClientOption
	for _, opt := range options {
		opt(&clientOpt)
	}
	s := Service{
		client:    clientOpt.HTTPClient,
		BasePath:  clientOpt.GetBasePath(),
		UserAgent: clientOpt.GetUserAgent(),
		BackOff:   clientOpt.GetBackoff(),
		NeedRetry: clientOpt.NeedRetry,
	}

	return &s
}
