package appstoreapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gh73962/appleapis/appstore/api/internal/httputils"
	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

func (s *Service) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if s.NeedRetry {
		return sendAndRetry(ctx, s.client, req, s.BackOff)
	}

	return send(ctx, s.client, req)
}

func SendRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	return send(ctx, client, req)
}

func send(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		errResp := datatypes.ErrorResponse{
			HTTPStatus: resp.StatusCode,
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		err = errors.Join(err, &errResp)
	}
	return resp, err
}

func SendAndRetry(ctx context.Context, client *http.Client, req *http.Request, bo Backoff) (*http.Response, error) {
	return sendAndRetry(ctx, client, req, bo)
}

func sendAndRetry(ctx context.Context, client *http.Client, req *http.Request, bo Backoff) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	if bo == nil {
		bo = httputils.NewBackoffImpl()
	}

	var (
		needRetry bool
		resp      *http.Response
		err       error
		interval  time.Duration
	)
	for {
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			if err != nil {
				return resp, errors.Join(ctx.Err(), err)
			}
			return resp, ctx.Err()
		case <-t.C:
		}

		if ctx.Err() != nil {
			if err != nil {
				return resp, errors.Join(ctx.Err(), err)
			}
			return resp, ctx.Err()
		}

		resp, err = client.Do(req.WithContext(ctx))
		if needRetry, err = shouldRetry(resp, err); !needRetry {
			break
		}

		interval = bo.Pause()
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
	return resp, err
}

func shouldRetry(resp *http.Response, err error) (bool, error) {
	if err == nil && resp.StatusCode == http.StatusOK {
		return false, err
	}

	if http.StatusInternalServerError <= resp.StatusCode && resp.StatusCode <= 599 {
		return true, err
	}

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout {
		return true, err
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true, err
	}

	errResp := datatypes.ErrorResponse{
		HTTPStatus: resp.StatusCode,
	}
	_ = json.NewDecoder(resp.Body).Decode(&errResp)
	err = errors.Join(err, &errResp)
	// see https://developer.apple.com/documentation/appstoreserverapi/error_codes
	if errResp.ErrorCode == 4040002 || errResp.ErrorCode == 4040004 || errResp.ErrorCode == 5000001 ||
		errResp.ErrorCode == 4040006 {
		return true, err
	}
	return false, err
}
