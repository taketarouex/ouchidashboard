package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

type Request struct {
	DeviceIDs []string `json:"deviceIDs"`
}

func TestE2E(t *testing.T) {
	func() {
		deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")
		request := Request{
			DeviceIDs: []string{deviceID},
		}
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf("wrong request: %s", string(requestJson))
		}
		resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(requestJson))
		if err != nil {
			t.Errorf("err: %v", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("status: %s", resp.Status)
		}
	}()
	func() {
		deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")
		request := Request{
			DeviceIDs: []string{deviceID, "test"},
		}
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf("wrong request: %s", string(requestJson))
		}
		resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(requestJson))
		if err != nil {
			t.Errorf("err: %v", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 400 {
			t.Errorf("status: %s", resp.Status)
		}
	}()
}
