package appstoreapi

import (
	"reflect"
	"testing"

	"github.com/gh73962/appleapis/appstore/api/v1/datatypes"
)

const (
	signedTransaction = `eyJhbGciOiJFUzI1NiIsIng1YyI6WyJleGFtcGxlMSIsImV4YW1wbGUyIiwiZXhhbXBsZTMxIl19.eyJ0cmFuc2FjdGlvbklkIjoiMTIzNDU2Nzg5Iiwib3JpZ2luYWxUcmFuc2FjdGlvbklkIjoiMTIzNDU2Nzg5IiwiYnVuZGxlSWQiOiJjb20ueHh4Lnh4eHgiLCJwcm9kdWN0SWQiOiJjb20ueHh4Lnh4eHgiLCJwdXJjaGFzZURhdGUiOjE2NzgwMzAyOTMwMDAsIm9yaWdpbmFsUHVyY2hhc2VEYXRlIjoxNjc4MDMwMjkzMDAwLCJxdWFudGl0eSI6MSwidHlwZSI6IkNvbnN1bWFibGUiLCJpbkFwcE93bmVyc2hpcFR5cGUiOiJQVVJDSEFTRUQiLCJzaWduZWREYXRlIjoxNjg3OTM3OTgyNjM0LCJlbnZpcm9ubWVudCI6IlByb2R1Y3Rpb24iLCJ0cmFuc2FjdGlvblJlYXNvbiI6IlBVUkNIQVNFIiwic3RvcmVmcm9udCI6IkNBTiIsInN0b3JlZnJvbnRJZCI6IjE0MzQ1NSJ9.nQe_caQQQdRH6HJQ8ZfugR_hh9xxxxxxohkVCjDbBwYXwRnBdmlKbxW3sE9MFnAMONzyE0AA`
)

func TestDecodeToJWSTransaction(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *datatypes.JWSTransaction
		wantErr bool
	}{
		{
			name: "test decode",
			args: args{signedTransaction},
			want: &datatypes.JWSTransaction{
				Header: datatypes.JWSDecodedHeader{
					Alg: "ES256",
					X5c: []string{"example1", "example2", "example31"},
				},
				Payload: datatypes.JWSTransactionDecodedPayload{
					BundleID:              "com.xxx.xxxx",
					Environment:           "Production",
					InAppOwnershipType:    "PURCHASED",
					OriginalPurchaseDate:  1678030293000,
					OriginalTransactionID: "123456789",
					ProductID:             "com.xxx.xxxx",
					PurchaseDate:          1678030293000,
					Quantity:              1,
					SignedDate:            1687937982634,
					Storefront:            "CAN",
					StorefrontID:          "143455",
					TransactionID:         "123456789",
					TransactionReason:     "PURCHASE",
					Type:                  "Consumable",
				},
				Signature: "nQe_caQQQdRH6HJQ8ZfugR_hh9xxxxxxohkVCjDbBwYXwRnBdmlKbxW3sE9MFnAMONzyE0AA",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeToJWSTransaction(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeToJWSTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeToJWSTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const signedRenewalInfo = `eyJhbGciOiJFUzI1NiIsIng1YyI6WyJleGFtcGxlMSIsImV4YW1wbGUyIiwiZXhhbXBsZTMxIl19.eyJvcmlnaW5hbFRyYW5zYWN0aW9uSWQiOiIxMjM0NTY3ODkxMTExMTEiLCJhdXRvUmVuZXdQcm9kdWN0SWQiOiJjb20ueHh4eC5zdWIiLCJwcm9kdWN0SWQiOiJjb20ueHh4eC5zdWIiLCJhdXRvUmVuZXdTdGF0dXMiOjEsInNpZ25lZERhdGUiOjE2ODY4ODIwNjM5MjksImVudmlyb25tZW50IjoiUHJvZHVjdGlvbiIsInJlY2VudFN1YnNjcmlwdGlvblN0YXJ0RGF0ZSI6MTY4MzA3MjgxNDAwMCwicmVuZXdhbERhdGUiOjE2ODc0ODY0NjAwMDB9.bRPbFO0cX3XhE1XuUoV8UgdZOD3vVjxxxxxxyW5qK7diTCqJ3A`

func TestDecodeToJWSRenewalInfo(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *datatypes.JWSRenewalInfo
		wantErr bool
	}{
		{
			name: "test decode",
			args: args{
				data: signedRenewalInfo,
			},
			want: &datatypes.JWSRenewalInfo{
				Header: datatypes.JWSDecodedHeader{
					Alg: "ES256",
					X5c: []string{"example1", "example2", "example31"},
				},
				Payload: datatypes.JWSRenewalInfoDecodedPayload{
					AutoRenewProductID:          "com.xxxx.sub",
					AutoRenewStatus:             1,
					Environment:                 "Production",
					OriginalTransactionID:       "123456789111111",
					ProductID:                   "com.xxxx.sub",
					RecentSubscriptionStartDate: 1683072814000,
					RenewalDate:                 1687486460000,
					SignedDate:                  1686882063929,
				},
				Signature: "bRPbFO0cX3XhE1XuUoV8UgdZOD3vVjxxxxxxyW5qK7diTCqJ3A",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeToJWSRenewalInfo(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeToJWSRenewalInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeToJWSRenewalInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
