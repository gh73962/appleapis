package appstoreapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

// TestNotification see https://developer.apple.com/documentation/appstoreserverapi/request_a_test_notification
func (s *Service) TestNotification(ctx context.Context, bearer string) (*datatypes.SendTestNotificationResponse, error) {
	req, err := http.NewRequest(http.MethodPost, s.BasePath+"notifications/test", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rsp datatypes.SendTestNotificationResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// GetTestNotificationStatus see https://developer.apple.com/documentation/appstoreserverapi/get_test_notification_status
func (s *Service) GetTestNotificationStatus(ctx context.Context, bearer, testNotificationToken string) (*datatypes.NotificationHistoryResponseItem, error) {
	req, err := http.NewRequest(http.MethodGet, s.BasePath+"notifications/test/"+testNotificationToken, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rsp datatypes.NotificationHistoryResponseItem
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// NotificationHistory see https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
func (s *Service) NotificationHistory(ctx context.Context, bearer, paginationToken string,
	nhr *datatypes.NotificationHistoryRequest) (*datatypes.NotificationHistoryResponse, error) {

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(nhr); err != nil {
		return nil, err
	}

	u := s.BasePath + "notifications/history"
	if paginationToken != "" {
		u += fmt.Sprintf("?paginationToken=%s", paginationToken)
	}

	req, err := http.NewRequest(http.MethodPost, u, &buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rsp datatypes.NotificationHistoryResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}
