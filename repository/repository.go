package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"github.com/tktkc72/ouchi-dashboard/collector"
)

type (
	Repository struct {
		documentRef  *firestore.DocumentRef
		documentSnap *firestore.DocumentSnapshot
		time         timeInterface
	}
	ouchiLog struct {
		Value     float64
		UpdatedAt time.Time
		CreatedAt time.Time
	}
	noRoom interface {
		noRoom() bool
	}
	noRoomErr struct {
		s string
	}
	nowTime       struct{}
	timeInterface interface {
		now() time.Time
	}
)

func IsNoRoom(err error) bool {
	no, ok := errors.Cause(err).(noRoom)
	return ok && no.noRoom()
}

func (e *noRoomErr) Error() string { return e.s }

func (e *noRoomErr) noRoom() bool { return true }

func (*nowTime) now() time.Time { return time.Now() }

func NewRepository(client *firestore.Client, rootPath, roomName string, time timeInterface) (collector.IRepository, error) {
	ctx := context.Background()
	ref := client.Collection(rootPath).Doc(roomName)
	snap, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !snap.Exists() {
		return nil, &noRoomErr{fmt.Sprintf("no room name: %s", roomName)}
	}
	return &Repository{
		documentRef:  ref,
		documentSnap: snap,
		time:         time,
	}, nil
}

func (r *Repository) SourceID() (string, error) {
	sourceID, err := r.documentSnap.DataAt("sourceID")
	if err != nil {
		return "", err
	}
	return sourceID.(string), nil
}

func (r *Repository) Add(collectLogs []collector.CollectLog) error {
	ctx := context.Background()
	for _, c := range collectLogs {
		o := ouchiLog{
			c.Value,
			c.UpdatedAt,
			r.time.now(),
		}
		_, _, err := r.documentRef.Collection(c.LogType.String()).Add(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}
