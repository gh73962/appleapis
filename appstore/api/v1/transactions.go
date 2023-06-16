package appstoreapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

// TransactionInfo see https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (s *Service) TransactionInfo(ctx context.Context, bearer, transactionID string) (*datatypes.JWSTransaction, error) {
	req, err := http.NewRequest(http.MethodGet, s.BasePath+"transactions/"+transactionID, nil)
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

	var rsp datatypes.TransactionInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return DecodeToJWSTransaction(rsp.SignedTransactionInfo)
}

// SendConsumptionInformation see https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information
func (s *Service) SendConsumptionInformation(ctx context.Context, bearer, transactionID string, cr *datatypes.ConsumptionRequest) error {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(cr); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, s.BasePath+"transactions/consumption/"+transactionID, &buff)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Do(ctx, req)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	return nil
}
