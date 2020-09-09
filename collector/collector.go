//go:generate mockgen -source=$GOFILE -destination=collector_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi-dashboard

package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/tenntenn/natureremo"
)

type (
	ICollector interface {
		Collect() error
	}
	// CollectorSevice service
	CollectorSevice struct {
		fetcher    IFetcher
		repository IRepository
	}

	LogType int

	CollectLog struct {
		Value     float64
		UpdatedAt time.Time
		LogType   LogType
		SourceID  string
	}
	Message struct {
		RoomNames []string `json:"RoomNames"`
	}
)

func (t LogType) String() string {
	switch t {
	case Temperature:
		return "temperature"
	case Humidity:
		return "humidity"
	case Illumination:
		return "illumination"
	case Motion:
		return "motion"
	default:
		return "Unknown"
	}
}

const (
	Temperature = iota
	Humidity
	Illumination
	Motion
)

func NewCollectorService(fetcher IFetcher, repository IRepository) ICollector {
	return &CollectorSevice{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (s *CollectorSevice) Collect() error {
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
	IRepository interface {
		SourceID() (string, error)
		Add([]CollectLog) error
	}
	noRoom interface {
		noRoom() bool
	}
	NoRoomErr struct {
		S string
	}
	NowTime       struct{}
	TimeInterface interface {
		Now() time.Time
	}
	IFetcher interface {
		fetch(deviceID string) ([]CollectLog, error)
	}
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

func IsNoRoom(err error) bool {
	no, ok := errors.Cause(err).(noRoom)
	return ok && no.noRoom()
}

func (e *NoRoomErr) Error() string { return e.S }

func (e *NoRoomErr) noRoom() bool { return true }

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
			Temperature,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorTypeHumidity].Value,
			d.NewestEvents[natureremo.SensorTypeHumidity].CreatedAt,
			Humidity,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensortypeIllumination].Value,
			d.NewestEvents[natureremo.SensortypeIllumination].CreatedAt,
			Illumination,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorType("mo")].Value,
			d.NewestEvents[natureremo.SensorType("mo")].CreatedAt,
			Motion,
			d.ID,
		},
	}
}

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
