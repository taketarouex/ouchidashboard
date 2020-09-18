// +build integration

package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/tktkc72/ouchidashboard/collector"
	"github.com/tktkc72/ouchidashboard/enum"
	"github.com/tktkc72/ouchidashboard/ouchi"
)

func TestRepository(t *testing.T) {
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
	mockTime := collector.NewMockTimeInterface(ctrl)
	mockNow := time.Date(2020, 7, 31, 10, 0, 0, 0, time.Local)
	mockTime.EXPECT().Now().AnyTimes().Return(mockNow)
	repository, err := NewRepository(client, rootPath, doc.ID, mockTime)
	if err != nil {
		t.Fatalf("failed to create repository because of: %v", err)
	}

	gotSourceID, err := repository.SourceID()
	if err != nil {
		t.Error("failed to get sourceID")
	}
	if sourceID != gotSourceID {
		t.Errorf("expect: %s, got: %s", sourceID, gotSourceID)
	}

	collectLogs := []collector.CollectLog{
		{Value: 0, UpdatedAt: time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), LogType: enum.Temperature, SourceID: "test"},
		{Value: 1, UpdatedAt: time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local), LogType: enum.Humidity, SourceID: "test"},
		{Value: 2, UpdatedAt: time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local), LogType: enum.Illumination, SourceID: "test"},
		{Value: 3, UpdatedAt: time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local), LogType: enum.Motion, SourceID: "test"},
	}

	if err = repository.Add(collectLogs); err != nil {
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

	fetched, err := repository.Fetch(doc.ID, enum.Temperature,
		time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 10, 0, 0, 0, time.Local),
		1, enum.Asc)
	if err != nil {
		t.Errorf("failed to fetch due to %v", err)
	}
	expected := []ouchi.Log{
		{Value: 0, UpdatedAt: time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), CreatedAt: mockNow},
	}

	if !cmp.Equal(expected, fetched) {
		t.Errorf("expected: %v, got: %v", expected, fetched[0])
	}

	// delete test data
	if _, err = doc.Delete(ctx); err != nil {
		t.Fatalf("failed to delete test data")
	}
}

func helperParseDocument(t *testing.T, d *firestore.DocumentRef, ctx context.Context) map[enum.LogType]ouchi.Log {
	t.Helper()
	returnMap := map[enum.LogType]ouchi.Log{}
	for _, l := range []enum.LogType{enum.Temperature, enum.Humidity, enum.Illumination, enum.Motion} {
		docs, err := d.Collection(l.String()).Documents(ctx).GetAll()
		if err != nil {
			t.Errorf("failed to get log type: %s", l)
			returnMap[l] = ouchi.Log{}
			continue
		}
		if len(docs) != 1 {
			t.Errorf("unexpected log length expected 1, got: %d", len(docs))
			returnMap[l] = ouchi.Log{}
			continue
		}
		returnMap[l] = ouchi.Log{
			Value:     docs[0].Data()["Value"].(float64),
			UpdatedAt: docs[0].Data()["UpdatedAt"].(time.Time),
			CreatedAt: docs[0].Data()["CreatedAt"].(time.Time),
		}
	}
	return returnMap
}
