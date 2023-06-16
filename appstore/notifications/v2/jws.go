package notifications

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func DecodeToJWSNotification(data string) (*JWSNotification, error) {
	header, payload, sig, err := DecodeSignedData(data)
	if err != nil {
		return nil, err
	}

	t := JWSNotification{
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
