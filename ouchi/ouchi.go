//go:generate mockgen -source=$GOFILE -destination=ouchi_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi

package ouchi

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tktkc72/ouchidashboard/collector"
	"github.com/tktkc72/ouchidashboard/enum"
)

type (
	// IOuchi is an interface of the ouchi service
	IOuchi interface {
		GetLogs(logType enum.LogType, start, end time.Time, opts ...getOption) ([]Log, error)
	}
	// Ouchi service
	Ouchi struct {
		repository IRepository
	}
	noRoom interface {
		noRoom() bool
	}
	// NoRoomErr is an error represents no doc with a specified room name
	NoRoomErr struct {
		S string
	}
	// IRepository is an interface of repository
	IRepository interface {
		SourceID() (string, error)
		Add([]collector.CollectLog) error
		Fetch(logType enum.LogType, start, end time.Time, limit int, order enum.Order) ([]Log, error)
	}
	// Log ouchi log
	Log struct {
		Value     float64
		UpdatedAt time.Time
		CreatedAt time.Time
	}
	getOpts struct {
		limit int
		order enum.Order
	}
	getOption func(*getOpts)
)

// Limit sets
func Limit(v int) getOption {
	return func(g *getOpts) {
		g.limit = v
	}
}

// Order sets order desc or asc.
func Order(v enum.Order) getOption {
	return func(g *getOpts) {
		g.order = v
	}
}

// IsNoRoom judge no room error
func IsNoRoom(err error) bool {
	no, ok := errors.Cause(err).(noRoom)
	return ok && no.noRoom()
}

func (e *NoRoomErr) Error() string { return e.S }

func (e *NoRoomErr) noRoom() bool { return true }

// NewOuchi creates a Ouchi service
func NewOuchi(repository IRepository) IOuchi {
	return &Ouchi{repository}
}

// GetLogs gets log
func (o *Ouchi) GetLogs(logType enum.LogType, start, end time.Time, opts ...getOption) ([]Log, error) {
	options := &getOpts{
		limit: 0,
		order: enum.Asc,
	}
	for _, setOpt := range opts {
		setOpt(options)
	}

	logs, err := o.repository.Fetch(logType, start, end, options.limit, options.order)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
