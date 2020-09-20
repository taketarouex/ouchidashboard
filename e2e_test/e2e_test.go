package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/tktkc72/ouchidashboard/collector"
)

func TestOuchi_E2E(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")
	deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")
	baseUrl, _ := url.Parse("http://localhost:8080")

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
		request := collector.Message{
			RoomNames: []string{doc.ID},
		}
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf("wrong request: %s", string(requestJson))
		}
		resp, err := http.Post(baseUrl.String(), "application/json", bytes.NewBuffer(requestJson))
		if err != nil {
			t.Errorf("err: %v", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("status: %s", resp.Status)
		}
	})
	t.Run("fail invalid roomName", func(t *testing.T) {
		request := collector.Message{
			RoomNames: []string{doc.ID, "test"},
		}
		requestJson, err := json.Marshal(request)
		if err != nil {
			t.Errorf("wrong request: %s", string(requestJson))
		}
		resp, err := http.Post(baseUrl.String(), "application/json", bytes.NewBuffer(requestJson))
		if err != nil {
			t.Errorf("err: %v", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 500 {
			t.Errorf("status: %s", resp.Status)
		}
	})

	t.Run("success", func(t *testing.T) {
		start := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		end := time.Now().Format(time.RFC3339)
		baseUrl.Path += fmt.Sprintf("/rooms/%s/logs/temperature", doc.ID)
		params := url.Values{}
		params.Add("start", start)
		params.Add("end", end)
		baseUrl.RawQuery = params.Encode()
		resp, err := http.Get(baseUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("failed to get logs due to: %v", resp.Status)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		t.Logf("body: %v", string(body))
	})

	// delete test data
	if _, err = doc.Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}
