// +build integration

package collector

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
)

func TestRepository_add(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	sourceID := os.Getenv("NATURE_REMO_DEVICE_ID")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create client because of:%v", err)
	}
	defer client.Close()

	// setup test data
	doc, _, err := client.Collection(rootPath).Add(ctx, map[string]string{"sourceID": sourceID})
	if err != nil {
		t.Fatalf("failed to create test data")
	}

	repository, err := NewRepository(client, rootPath, sourceID)
	if err != nil {
		t.Fatalf("failed to create repository because of: %v", err)
	}
	collectLogs := []collectLog{
		{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), temperature, "test"},
		{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), humidity, "test"},
		{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), illumination, "test"},
		{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), motion, "test"},
	}

	if err = repository.add(collectLogs); err != nil {
		t.Fatalf("error: %v", err)
	}

	// delete test data
	if _, err = doc.Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}
