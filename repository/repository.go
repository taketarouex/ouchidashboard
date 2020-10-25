package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"github.com/tktkc72/ouchidashboard/collector"
	"github.com/tktkc72/ouchidashboard/enum"
	"github.com/tktkc72/ouchidashboard/ouchi"
)

type (
	// Repository data store
	Repository struct {
		rootCollection *firestore.CollectionRef
		time           collector.TimeInterface
	}
)

// NewRepository creates repository which has the name specified
func NewRepository(client *firestore.Client, rootPath string, time collector.TimeInterface) (ouchi.IRepository, error) {
	return &Repository{
		rootCollection: client.Collection(rootPath),
		time:           time,
	}, nil
}

// SourceID gets a sourceID at the root document
func (r *Repository) SourceID(roomName string) (string, error) {
	ctx := context.Background()
	ref := r.rootCollection.Doc(roomName)
	snap, err := ref.Get(ctx)
	if err != nil {
		return "", err
	}
	if !snap.Exists() {
		return "", &ouchi.NoRoomErr{S: fmt.Sprintf("no room name: %s", roomName)}
	}

	sourceID, err := snap.DataAt("sourceID")
	if err != nil {
		return "", err
	}
	return sourceID.(string), nil
}

func (r *Repository) existRoom(roomName string) error {
	ctx := context.Background()
	ref := r.rootCollection.Doc(roomName)
	_, err := ref.Get(ctx)
	if err != nil {
		return &ouchi.NoRoomErr{S: fmt.Sprintf("no room name: %s", roomName)}
	}
	return nil
}

// Add adds CollectLogs to repository
func (r *Repository) Add(roomName string, collectLogs []collector.CollectLog) error {
	err := r.existRoom(roomName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ref := r.rootCollection.Doc(roomName)
	for _, c := range collectLogs {
		o := ouchi.Log{
			Value:     c.Value,
			UpdatedAt: c.UpdatedAt,
			CreatedAt: r.time.Now(),
		}
		_, _, err := ref.Collection(c.LogType.String()).Add(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fetch fetches logs from repository
func (r *Repository) Fetch(
	roomName string,
	logType enum.LogType,
	start, end time.Time,
	limit int, order enum.Order) ([]ouchi.Log, error) {
	err := r.existRoom(roomName)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ref := r.rootCollection.Doc(roomName)

	collection := ref.Collection(logType.String())
	if collection == nil {
		return nil, errors.Errorf("no collection type: %s", logType.String())
	}

	direction := firestore.Asc
	if order == enum.Desc {
		direction = firestore.Desc
	}

	query := collection.Where("UpdatedAt", ">=", start).Where("UpdatedAt", "<", end).OrderBy("UpdatedAt", direction)
	if limit > 0 {
		query = query.Limit(limit)
	}
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Errorf("failed to get documents due to: %v", err)
	}

	return parse(docs), nil
}

func parse(docs []*firestore.DocumentSnapshot) []ouchi.Log {
	parsed := []ouchi.Log{}
	for _, doc := range docs {
		log := ouchi.Log{
			Value:     doc.Data()["Value"].(float64),
			UpdatedAt: doc.Data()["UpdatedAt"].(time.Time),
			CreatedAt: doc.Data()["CreatedAt"].(time.Time),
		}
		parsed = append(parsed, log)
	}
	return parsed
}

func (r *Repository) FetchRoomNames() (roomNames []string, err error) {
	ctx := context.Background()
	roomDocs, err := r.rootCollection.DocumentRefs(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	for _, roomDoc := range roomDocs {
		roomNames = append(roomNames, roomDoc.ID)
	}
	return roomNames, nil
}
