// +build acceptance

package main_test

import (
	"fmt"
	"io"
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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(fmt.Errorf("Read response body: %w", err))
	}
	if string(bodyBytes) != "Hello World !" {
		t.Error(fmt.Errorf("Wrong message %s: %w", string(bodyBytes), err))
	}
}
