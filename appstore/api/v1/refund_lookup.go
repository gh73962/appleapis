package appstoreapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

// RefundHistory see https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (s *Service) RefundHistory(ctx context.Context, bearer, transactionID, revision string) (*datatypes.OrderLookupResponse, error) {
	req, err := http.NewRequest(http.MethodGet, s.BasePath+"lookup/"+transactionID+"?revision="+revision, nil)
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

	var rsp datatypes.OrderLookupResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}
