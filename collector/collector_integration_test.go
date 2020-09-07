// +build integration

package collector

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	gomock "github.com/golang/mock/gomock"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTime := NewMocktimeInterface(ctrl)
	mockNow := time.Date(2020, 7, 31, 10, 0, 0, 0, time.Local)
	mockTime.EXPECT().now().AnyTimes().Return(mockNow)
	repository, err := NewRepository(client, rootPath, doc.ID, mockTime)
	if err != nil {
		t.Fatalf("failed to create repository because of: %v", err)
	}

	gotSourceID, err := repository.sourceID()
	if err != nil {
		t.Error("failed to get sourceID")
	}
	if sourceID != gotSourceID {
		t.Errorf("expect: %s, got: %s", sourceID, gotSourceID)
	}

	collectLogs := []collectLog{
		{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), temperature, "test"},
		{1, time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local), humidity, "test"},
		{2, time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local), illumination, "test"},
		{3, time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local), motion, "test"},
	}

	if err = repository.add(collectLogs); err != nil {
		t.Fatalf("error: %v", err)
	}

	ouchiLogMap := helperParseDocument(t, doc, ctx)
	for _, l := range collectLogs {
		if l.Value != ouchiLogMap[l.LogType].Value {
			t.Errorf("assert error value expect: %v, got: %v", l.Value, ouchiLogMap[l.LogType].Value)
		}
		if !l.UpdatedAt.Equal(ouchiLogMap[l.LogType].UpdatedAt) {
			t.Errorf("assert error updatedAt expect: %v, got: %v", l.UpdatedAt, ouchiLogMap[l.LogType].UpdatedAt)
		}
		if !mockNow.Equal(ouchiLogMap[l.LogType].CreatedAt) {
			t.Errorf("assert error createdAt expect: %v, got: %v", mockNow, ouchiLogMap[l.LogType].CreatedAt)
		}
	}

	// delete test data
	if _, err = doc.Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}

func helperParseDocument(t *testing.T, d *firestore.DocumentRef, ctx context.Context) map[logType]ouchiLog {
	t.Helper()
	returnMap := map[logType]ouchiLog{}
	for _, l := range []logType{temperature, humidity, illumination, motion} {
		docs, err := d.Collection(l.String()).Documents(ctx).GetAll()
		if err != nil {
			t.Errorf("failed to get log type: %s", l)
			returnMap[l] = ouchiLog{}
			continue
		}
		if len(docs) != 1 {
			t.Errorf("unexpected log length expected 1, got: %d", len(docs))
			returnMap[l] = ouchiLog{}
			continue
		}
		returnMap[l] = ouchiLog{
			docs[0].Data()["Value"].(float64),
			docs[0].Data()["UpdatedAt"].(time.Time),
			docs[0].Data()["CreatedAt"].(time.Time),
		}
	}
	return returnMap
}
