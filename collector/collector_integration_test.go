// +build integration

package collector

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
)

func TestRepository(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	documentPath := os.Getenv("FIRESTORE_DOC_PATH")

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Errorf("cant create client because of:%v", err)
	}
	defer client.Close()

	repository := NewRepository(client, documentPath)
	collectLogs := []collectLog{
		{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), temperature, "test"},
		{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), humidity, "test"},
		{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), illumination, "test"},
		{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), motion, "test"},
	}
	err = repository.add(collectLogs)
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
