//go:generate mockgen -source=$GOFILE -destination=collector_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi-dashboard

package collector

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
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

	logType int

	collectLog struct {
		Value     float64
		UpdatedAt time.Time
		LogType   logType
		SourceID  string
	}
)

func (t logType) String() string {
	switch t {
	case temperature:
		return "temperature"
	case humidity:
		return "humidity"
	case illumination:
		return "illumination"
	case motion:
		return "motion"
	default:
		return "Unknown"
	}
}

const (
	temperature = iota
	humidity
	illumination
	motion
)

func NewCollectorService(fetcher IFetcher, repository IRepository) ICollector {
	return &CollectorSevice{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (s *CollectorSevice) Collect() error {
	sourceID, err := s.repository.sourceID()
	if err != nil {
		return err
	}

	collected, err := s.fetcher.fetch(sourceID)
	if err != nil {
		return err
	}

	err = s.repository.add(collected)
	if err != nil {
		return err
	}

	return nil
}

type (
	IRepository interface {
		sourceID() (string, error)
		add([]collectLog) error
	}
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

func NewRepository(client *firestore.Client, rootPath, roomName string, time timeInterface) (IRepository, error) {
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

func (r *Repository) add(collectLogs []collectLog) error {
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

func (r *Repository) sourceID() (string, error) {
	sourceID, err := r.documentSnap.DataAt("sourceID")
	if err != nil {
		return "", err
	}
	return sourceID.(string), nil
}

type (
	IFetcher interface {
		fetch(deviceID string) ([]collectLog, error)
	}
	Fetcher struct {
		client *natureremo.Client
	}
	deviceSlice []*natureremo.Device
)

func (rcv deviceSlice) where(fn func(*natureremo.Device) bool) (result deviceSlice) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func (rcv deviceSlice) fetchLog() []collectLog {
	var collectLogs []collectLog
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

func parseNatureremoDevice(d *natureremo.Device) []collectLog {
	return []collectLog{
		{
			d.NewestEvents[natureremo.SensorTypeTemperature].Value,
			d.NewestEvents[natureremo.SensorTypeTemperature].CreatedAt,
			temperature,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorTypeHumidity].Value,
			d.NewestEvents[natureremo.SensorTypeHumidity].CreatedAt,
			humidity,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensortypeIllumination].Value,
			d.NewestEvents[natureremo.SensortypeIllumination].CreatedAt,
			illumination,
			d.ID,
		},
		{
			d.NewestEvents[natureremo.SensorType("mo")].Value,
			d.NewestEvents[natureremo.SensorType("mo")].CreatedAt,
			motion,
			d.ID,
		},
	}
}

type noDevice interface {
	noDevice() bool
}

func IsNoDevice(err error) bool {
	no, ok := errors.Cause(err).(noDevice)
	return ok && no.noDevice()
}

type noDeviceErr struct {
	s string
}

func (e *noDeviceErr) Error() string { return e.s }

func (e *noDeviceErr) noDevice() bool { return true }

func (f *Fetcher) fetch(deviceID string) ([]collectLog, error) {
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
