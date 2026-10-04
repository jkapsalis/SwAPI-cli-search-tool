// Package tests holds black-box tests for the api and cli packages. Every
// test talks to a local fake SWAPI server, so no network access is needed.
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("encode response: %v", err)
	}
}
