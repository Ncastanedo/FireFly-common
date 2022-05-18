// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ffcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/ffresty"
	"github.com/stretchr/testify/assert"
)

func newTestClient(t *testing.T, response ffcapiResponse) (*apiClient, func()) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		resBytes, err := json.Marshal(response)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(resBytes)))
		if response.ErrorMessage() != "" {
			w.WriteHeader(500)
		}
		_, err = w.Write(resBytes)
		assert.NoError(t, err)
	}))
	section := config.RootSection("unittest")
	ffresty.InitConfig(section)
	section.Set(ffresty.HTTPConfigURL, fmt.Sprintf("http://%s", server.Listener.Addr()))
	ctx := context.Background()
	api := NewFFCAPIClient(ctx, section, VariantEVM)
	return api.(*apiClient), server.Close
}

func TestBadResponseContentType(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not JSON"))
	}))
	defer server.Close()
	section := config.RootSection("unittest")
	ffresty.InitConfig(section)
	section.Set(ffresty.HTTPConfigURL, fmt.Sprintf("http://%s", server.Listener.Addr()))
	ctx := context.Background()

	api := NewFFCAPIClient(ctx, section, VariantEVM).(*apiClient)
	_, err := api.invokeAPI(ctx, &ExecQueryRequest{}, &ResponseBase{})
	assert.Regexp(t, "FF00156", err)

}

func TestBadResponseError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not JSON"))
	}))
	config.RootConfigReset()
	section := config.RootSection("ffcapi_tests")
	ffresty.InitConfig(section)
	section.Set(ffresty.HTTPConfigURL, fmt.Sprintf("http://%s", server.Listener.Addr()))
	ctx := context.Background()

	server.Close()
	api := NewFFCAPIClient(ctx, section, VariantEVM)
	_, _, err := api.ExecQuery(ctx, &ExecQueryRequest{})
	assert.Regexp(t, "FF00155", err)

}