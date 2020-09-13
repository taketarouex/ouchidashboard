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
	Ouchi struct {
		repository IRepository
	}
	// IRepository is an interface of repository
	IRepository interface {
		fetch(roomName string, start, end time.Time, limit int, order enum.Order) ([]Log, error)
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
func NewOuchi(repository IRepository) IOuchi {
	return &Ouchi{repository}
}

// GetTemperature gets temperature log
func (o *Ouchi) GetTemperature(roomName string, start, end time.Time, opts ...getOption) ([]Log, error) {
	options := &getOpts{
		limit: 0,
		order: enum.Asc,
	}
	for _, setOpt := range opts {
		setOpt(options)
	}

	logs, err := o.repository.fetch(roomName, start, end, options.limit, options.order)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
