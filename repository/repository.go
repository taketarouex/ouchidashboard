package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"github.com/tktkc72/ouchi/collector"
	"github.com/tktkc72/ouchi/enum"
	"github.com/tktkc72/ouchi/ouchi"
)

type (
	// Repository data store
	Repository struct {
		rootCollection *firestore.CollectionRef
		documentRef    *firestore.DocumentRef
		documentSnap   *firestore.DocumentSnapshot
		time           collector.TimeInterface
	}
)

// NewRepository creates repository which has the name specified
func NewRepository(client *firestore.Client, rootPath, roomName string, time collector.TimeInterface) (ouchi.IRepository, error) {
	ctx := context.Background()
	ref := client.Collection(rootPath).Doc(roomName)
	snap, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !snap.Exists() {
		return nil, &ouchi.NoRoomErr{S: fmt.Sprintf("no room name: %s", roomName)}
	}
	return &Repository{
		rootCollection: client.Collection(rootPath),
		documentRef:    ref,
		documentSnap:   snap,
		time:           time,
	}, nil
}

// SourceID gets a sourceID at the root document
func (r *Repository) SourceID() (string, error) {
	sourceID, err := r.documentSnap.DataAt("sourceID")
	if err != nil {
		return "", err
	}
	return sourceID.(string), nil
}

// Add adds CollectLogs to repository
func (r *Repository) Add(collectLogs []collector.CollectLog) error {
	ctx := context.Background()
	for _, c := range collectLogs {
		o := ouchi.Log{
			Value:     c.Value,
			UpdatedAt: c.UpdatedAt,
			CreatedAt: r.time.Now(),
		}
		_, _, err := r.documentRef.Collection(c.LogType.String()).Add(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fetch fetches logs from repository
func (r *Repository) Fetch(roomName string,
	logType enum.LogType,
	start, end time.Time,
	limit int, order enum.Order) ([]ouchi.Log, error) {
	roomDoc := r.rootCollection.Doc(roomName)
	if roomDoc == nil {
		return nil, &ouchi.NoRoomErr{S: fmt.Sprintf("no room name: %s", roomName)}
	}
	collection := roomDoc.Collection(logType.String())
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
	ctx := context.Background()
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
