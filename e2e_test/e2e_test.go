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
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/tktkc72/ouchidashboard/collector"
	"github.com/tktkc72/ouchidashboard/ouchi"
)

func TestCollector_E2E(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")
	deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")
	baseUrl, _ := url.Parse("http://localhost:8080")
	roomName := "testRoom"

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create client because of:%v", err)
	}
	defer client.Close()
	// setup test data
	_, err = client.Doc(fmt.Sprintf("%s/%s", rootPath, roomName)).Create(ctx, map[string]string{"sourceID": deviceID})
	if err != nil {
		t.Fatalf("failed to create test data")
	}

	t.Run("success to collect", func(t *testing.T) {
		request := collector.Message{
			RoomNames: []string{roomName},
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
			RoomNames: []string{roomName, "invalidRoom"},
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

	// delete test data
	if _, err = client.Doc(fmt.Sprintf("%s/%s", rootPath, roomName)).Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}

func TestOuchi_E2E(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")
	deviceID := os.Getenv("NATURE_REMO_DEVICE_ID")
	baseUrl := "http://localhost:8080"
	roomName := "testRoom"

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create client because of:%v", err)
	}
	defer client.Close()
	// setup test data
	_, err = client.Doc(fmt.Sprintf("%s/%s", rootPath, roomName)).Create(ctx, map[string]string{"sourceID": deviceID})
	if err != nil {
		t.Fatalf("failed to create test data")
	}

	testLogs := []ouchi.Log{
		{Value: 0, UpdatedAt: time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC)},
		{Value: 1, UpdatedAt: time.Date(2020, 1, 23, 1, 0, 0, 0, time.UTC), CreatedAt: time.Date(2020, 1, 23, 1, 0, 0, 0, time.UTC)},
		{Value: 2, UpdatedAt: time.Date(2020, 1, 23, 2, 0, 0, 0, time.UTC), CreatedAt: time.Date(2020, 1, 23, 2, 0, 0, 0, time.UTC)},
	}

	for _, testLog := range testLogs {
		_, _, err := client.Collection(fmt.Sprintf("%s/%s/temperature", rootPath, roomName)).Add(ctx, testLog)
		if err != nil {
			t.Fatalf("failed to create test logs due to: %v", err)
		}
	}

	t.Run("success to get logs", func(t *testing.T) {
		apiUrl, err := url.Parse(baseUrl)
		if err != nil {
			t.Fatalf("failed to parse url %v", err)
		}
		apiUrl.Path += fmt.Sprintf("/api/rooms/%s/logs/temperature", roomName)
		start := time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		end := time.Date(2020, 1, 23, 2, 0, 0, 0, time.UTC).Format(time.RFC3339)
		params := url.Values{}
		params.Add("start", start)
		params.Add("end", end)
		apiUrl.RawQuery = params.Encode()
		resp, err := http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("failed to get logs due to: %v", resp.Status)
		}
		defer resp.Body.Close()
		actual := new([]ouchi.Log)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read body due to: %v", err)
		}
		if err = json.Unmarshal(body, actual); err != nil {
			t.Errorf("failed to unmarshal due to: %v", err)
		}
		expected := testLogs[0:2]
		if !reflect.DeepEqual(actual, &expected) {
			t.Errorf("got: %v, expect: %v", actual, &expected)
		}

		// add options
		params.Add("limit", "1")
		params.Add("order", "DESC")
		apiUrl.RawQuery = params.Encode()
		resp, err = http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("failed to get logs due to: %v", resp.Status)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read body due to: %v", err)
		}
		if err = json.Unmarshal(body, actual); err != nil {
			t.Errorf("failed to unmarshal due to: %v", err)
		}
		expected = []ouchi.Log{testLogs[1]}
		if !reflect.DeepEqual(actual, &expected) {
			t.Errorf("got: %v, expect: %v", actual, &expected)
		}
	})

	t.Run("fail invalid query params", func(t *testing.T) {
		apiUrl, err := url.Parse(baseUrl)
		if err != nil {
			t.Fatalf("failed to parse url %v", err)
		}
		apiUrl.Path += fmt.Sprintf("/api/rooms/%s/logs/temperature", roomName)
		start := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		end := time.Now().Format(time.RFC3339)
		params := url.Values{}
		params.Add("start", start)
		apiUrl.RawQuery = params.Encode()
		resp, err := http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("expected 400 but got: %v", resp.Status)
		}
		params.Del("start")
		params.Add("end", end)
		apiUrl.RawQuery = params.Encode()
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("expected 400 but got: %v", resp.Status)
		}

		// invalid format
		invalidStart := time.Now().AddDate(0, 0, -1).Format(time.RFC1123)
		invalidEnd := time.Now().Format(time.RFC1123)
		params.Add("start", invalidStart)
		apiUrl.RawQuery = params.Encode()
		resp, err = http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("expected 400 but got: %v", resp.Status)
		}

		params.Set("start", start)
		params.Set("end", invalidEnd)
		apiUrl.RawQuery = params.Encode()
		resp, err = http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("expected 400 but got: %v", resp.Status)
		}
	})

	t.Run("fail invalid url parameter", func(t *testing.T) {
		apiUrl, err := url.Parse(baseUrl)
		if err != nil {
			t.Fatalf("failed to parse url %v", err)
		}
		apiUrl.Path += fmt.Sprintf("/api/rooms/%s/logs/temperature", "notFoundRoom")
		start := time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		end := time.Date(2020, 1, 23, 2, 0, 0, 0, time.UTC).Format(time.RFC3339)
		params := url.Values{}
		params.Add("start", start)
		params.Add("end", end)
		apiUrl.RawQuery = params.Encode()
		resp, err := http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("expected 400 but got: %v", resp.Status)
		}
	})
	t.Run("success to get room names", func(t *testing.T) {
		apiUrl, err := url.Parse(baseUrl)
		if err != nil {
			t.Fatalf("failed to parse url %v", err)
		}
		apiUrl.Path += "/api/rooms"
		resp, err := http.Get(apiUrl.String())
		if err != nil {
			t.Errorf("failed to http get due to: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("failed to get room names due to: %v", resp.Status)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read body due to: %v", err)
		}
		actual := new([]string)
		if err = json.Unmarshal(body, actual); err != nil {
			t.Errorf("failed to unmarshal due to: %v", err)
		}
		expected := []string{roomName}
		if !reflect.DeepEqual(actual, &expected) {
			t.Errorf("got: %v, expect: %v", actual, &expected)
		}
	})

	// delete test data
	//if _, err = client.Doc(fmt.Sprintf("%s/%s", rootPath, roomName)).Delete(ctx); err != nil {
	//	t.Fatalf("failed to delete test data")
	//}
}
