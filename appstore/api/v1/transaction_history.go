package appstoreapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

// TransactionHistory see https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
// TODO Query Parameters
func (s *Service) TransactionHistory(ctx context.Context, bearer, transactionID string) (*datatypes.HistoryResponse, error) {
	req, err := http.NewRequest(http.MethodGet, s.BasePath+"history/"+transactionID, nil)
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

	var rsp datatypes.HistoryResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}
