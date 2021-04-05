// +build acceptance

package main_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	executablePath := os.Getenv("EXECUTABLE_PATH")
	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); nil != err {
		t.Fatal(fmt.Errorf("Could not start command '%s': %w", executablePath, err))
	}
	defer func() {
		if err := cmd.Process.Signal(os.Interrupt); nil != err {
			t.Fatal(fmt.Errorf("Could not stop command correctly '%s': %w", executablePath, err))
		}
	}()

	go func() {
		if err := cmd.Wait(); nil != err {
			t.Fatal(fmt.Errorf("Error when running command '%s': %w", executablePath, err))
		}
	}()

	time.Sleep(500 * time.Millisecond)

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
