package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/tktkc72/ouchi/collector"
)

type (
	// Repository data store
	Repository struct {
		documentRef  *firestore.DocumentRef
		documentSnap *firestore.DocumentSnapshot
		time         collector.TimeInterface
	}
	ouchiLog struct {
		Value     float64
		UpdatedAt time.Time
		CreatedAt time.Time
	}
)

// NewRepository creates repository which has the name specified
func NewRepository(client *firestore.Client, rootPath, roomName string, time collector.TimeInterface) (collector.IRepository, error) {
	ctx := context.Background()
	ref := client.Collection(rootPath).Doc(roomName)
	snap, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !snap.Exists() {
		return nil, &collector.NoRoomErr{S: fmt.Sprintf("no room name: %s", roomName)}
	}
	return &Repository{
		documentRef:  ref,
		documentSnap: snap,
		time:         time,
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
		o := ouchiLog{
			c.Value,
			c.UpdatedAt,
			r.time.Now(),
		}
		_, _, err := r.documentRef.Collection(c.LogType.String()).Add(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}
