package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

type Request struct {
	DeviceIDs []string `json:"deviceIDs"`
}

func TestCollector_E2E(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")
	deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create client because of:%v", err)
	}
	defer client.Close()
	// setup test data
	doc, _, err := client.Collection(rootPath).Add(ctx, map[string]string{"sourceID": deviceID})
	if err != nil {
		t.Fatalf("failed to create test data")
	}

	t.Run("success", func(t *testing.T) {
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
	})
	t.Run("fail invalid deviceID", func(t *testing.T) {
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
		if resp.StatusCode != 500 {
			t.Errorf("status: %s", resp.Status)
		}
	})

	// delete test data
	if _, err = doc.Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}
