// +build acceptance

package main_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpServer(t *testing.T) {
	resp, err := http.Get("http://localhost:18080")
	if err != nil {
		t.Fatal(fmt.Errorf("Querying server: %w", err))
	}
	if http.StatusOK != resp.StatusCode {
		t.Error(fmt.Errorf("Wrong status %d: %w", resp.StatusCode, err))
	}
}
