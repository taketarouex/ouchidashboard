//go:generate mockgen -source=$GOFILE -destination=ouchi_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi

package ouchi

import (
	"time"

	"github.com/tktkc72/ouchi/enum"
)

type (
	// IOuchi is an interface of the ouchi service
	IOuchi interface {
		GetTemperature(roomName string, start, end time.Time, opts ...getOption) ([]Log, error)
	}
	// Ouchi service
	Ouchi struct{}
	// IRepository is an interface of repository
	IRepository interface {
		fetch(roomName string, start, end time.Time, limit int, order enum.Order)
	}
	// Log ouchi log
	Log struct {
		Value     float64
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

// NewOuchi creates a Ouchi service
func NewOuchi() IOuchi {
	return &Ouchi{}
}

// GetTemperature gets temperature log
func (*Ouchi) GetTemperature(roomName string, start, end time.Time, opts ...getOption) ([]Log, error) {
	o := &getOpts{
		limit: 0,
		order: enum.Asc,
	}
	for _, opt := range opts {
		opt(o)
	}

	return []Log{}, nil
}
