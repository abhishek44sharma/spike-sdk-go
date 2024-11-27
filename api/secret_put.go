//    \\ SPIKE: Secure your secrets with SPIFFE.
//  \\\\\ Copyright 2024-present SPIKE contributors.
// \\\\\\\ SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"errors"

	"github.com/spiffe/go-spiffe/v2/workloadapi"

	"github.com/spiffe/spike-sdk-go/api/internal/entity/v1/reqres"
	"github.com/spiffe/spike-sdk-go/api/internal/url"
	"github.com/spiffe/spike-sdk-go/net"
)

// PutSecret creates or updates a secret at the specified path with the given
// values using mTLS authentication.
//
// Parameters:
//   - source: X509Source for mTLS client authentication
//   - path: Path where the secret should be stored
//   - values: Map of key-value pairs representing the secret data
//
// Returns:
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	err := PutSecret(x509Source, "secret/path", map[string]string{"key": "value"})
func PutSecret(source *workloadapi.X509Source,
	path string, values map[string]string) error {

	r := reqres.SecretPutRequest{
		Path:   path,
		Values: values,
	}

	mr, err := json.Marshal(r)
	if err != nil {
		return errors.Join(
			errors.New("putSecret: I am having problem generating the payload"),
			err,
		)
	}

	var truer = func(string) bool { return true }
	client, err := net.CreateMtlsClient(source, truer)
	if err != nil {
		return err
	}

	body, err := net.Post(client, url.SecretPut(), mr)
	if err != nil {
		return err
	}

	res := reqres.SecretPutResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return errors.Join(
			errors.New("putSecret: Problem parsing response body"),
			err,
		)
	}
	if res.Err != "" {
		return errors.New(string(res.Err))
	}

	return nil
}
