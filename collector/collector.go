//go:generate mockgen -source=$GOFILE -destination=collector_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi

package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/tenntenn/natureremo"
	"github.com/tktkc72/ouchi/enum"
)

type (
	// ICollector is an interface of the collector service
	ICollector interface {
		Collect() error
	}
	// Sevice collector service
	Sevice struct {
		fetcher    IFetcher
		repository IRepository
	}

	// CollectLog is a model of collected logs
	CollectLog struct {
		Value     float64
		UpdatedAt time.Time
		LogType   enum.LogType
		SourceID  string
	}
	// Message is a struct of a requests
	Message struct {
		RoomNames []string `json:"RoomNames"`
	}
)

// NewCollectorService creates a service
func NewCollectorService(fetcher IFetcher, repository IRepository) ICollector {
	return &Sevice{
		fetcher:    fetcher,
		repository: repository,
	}
}

// Collect is a use case
func (s *Sevice) Collect() error {
	sourceID, err := s.repository.SourceID()
	if err != nil {
		return err
	}

	collected, err := s.fetcher.fetch(sourceID)
	if err != nil {
		return err
	}

	err = s.repository.Add(collected)
	if err != nil {
		return err
	}

	return nil
}

type (
	// IRepository is an interface of repository
	IRepository interface {
		SourceID() (string, error)
		Add([]CollectLog) error
	}
	noRoom interface {
		noRoom() bool
	}
	// NoRoomErr is an error represents no doc with a specified room name
	NoRoomErr struct {
		S string
	}
	// NowTime is a utility to return current time
	NowTime struct{}
	// TimeInterface is an interface of NowTime
	TimeInterface interface {
		Now() time.Time
	}
	// IFetcher is an interface of fetching logs
	IFetcher interface {
		fetch(deviceID string) ([]CollectLog, error)
	}
	// Fetcher is a struct which fetches logs
	Fetcher struct {
		client *natureremo.Client
	}
	deviceSlice []*natureremo.Device
	noDevice    interface {
		noDevice() bool
	}
	noDeviceErr struct {
		s string
	}
)

// IsNoRoom judge no room error
func IsNoRoom(err error) bool {
	no, ok := errors.Cause(err).(noRoom)
	return ok && no.noRoom()
}

func (e *NoRoomErr) Error() string { return e.S }

func (e *NoRoomErr) noRoom() bool { return true }

// Now returns current time
func (*NowTime) Now() time.Time { return time.Now() }

func (rcv deviceSlice) where(fn func(*natureremo.Device) bool) (result deviceSlice) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func (rcv deviceSlice) fetchLog() []CollectLog {
	var collectLogs []CollectLog
	for _, d := range rcv {
		collectLogs = append(collectLogs, parseNatureremoDevice(d)...)
	}
	return collectLogs
}

// NewFetcher creates Fetcher
func NewFetcher(client *natureremo.Client) IFetcher {
	return &Fetcher{
		client: client,
	}
}

func parseNatureremoDevice(d *natureremo.Device) []CollectLog {
	return []CollectLog{
		{
			d.NewestEvents[natureremo.SensorTypeTemperature].Value,
			d.NewestEvents[natureremo.SensorTypeTemperature].CreatedAt,
			enum.Temperature,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorTypeHumidity].Value,
			d.NewestEvents[natureremo.SensorTypeHumidity].CreatedAt,
			enum.Humidity,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensortypeIllumination].Value,
			d.NewestEvents[natureremo.SensortypeIllumination].CreatedAt,
			enum.Illumination,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorType("mo")].Value,
			d.NewestEvents[natureremo.SensorType("mo")].CreatedAt,
			enum.Motion,
			d.ID,
		},
	}
}

// IsNoDevice judges no device error
func IsNoDevice(err error) bool {
	no, ok := errors.Cause(err).(noDevice)
	return ok && no.noDevice()
}

func (e *noDeviceErr) Error() string { return e.s }

func (e *noDeviceErr) noDevice() bool { return true }

func (f *Fetcher) fetch(deviceID string) ([]CollectLog, error) {
	ctx := context.Background()
	var devices deviceSlice
	devices, err := f.client.DeviceService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	targetDevice := devices.where(func(d *natureremo.Device) bool {
		return d.ID == deviceID
	})
	if targetDevice == nil {
		return nil, &noDeviceErr{fmt.Sprintf("no device id: %s", deviceID)}
	}
	collectLogs := targetDevice.fetchLog()

	return collectLogs, nil
}
