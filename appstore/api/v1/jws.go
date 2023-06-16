package appstoreapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

func DecodeToJWSTransaction(data string) (*datatypes.JWSTransaction, error) {
	header, payload, sig, err := DecodeSignedData(data)
	if err != nil {
		return nil, err
	}

	t := datatypes.JWSTransaction{
		Signature: sig,
	}
	if err = json.Unmarshal(header, &t.Header); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(payload, &t.Payload); err != nil {
		return nil, err
	}

	return &t, nil
}

func DecodeToJWSRenewalInfo(data string) (*datatypes.JWSRenewalInfo, error) {
	header, payload, sig, err := DecodeSignedData(data)
	if err != nil {
		return nil, err
	}

	t := datatypes.JWSRenewalInfo{
		Signature: sig,
	}
	if err = json.Unmarshal(header, &t.Header); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(payload, &t.Payload); err != nil {
		return nil, err
	}

	return &t, nil
}

func DecodeSignedData(data string) ([]byte, []byte, string, error) {
	array := strings.Split(data, ".")
	if len(array) != 3 {
		return nil, nil, "", errors.New("invalid signed data")
	}

	header, err := base64.RawStdEncoding.DecodeString(array[0])
	if err != nil {
		return nil, nil, "", err
	}
	payload, err := base64.RawStdEncoding.DecodeString(array[1])
	if err != nil {
		return nil, nil, "", err
	}

	return header, payload, array[2], nil
}
