//go:generate mockgen -source=$GOFILE -destination=collector_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi-dashboard

package collector

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
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
		value     float64
		updatedAt time.Time
		logType   logType
		sourceId  string
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
	collected, err := s.fetcher.fetch()
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
		add([]collectLog) error
	}
	Repository struct {
		document *firestore.DocumentRef
	}
)

func NewRepository(client *firestore.Client, document string) IRepository {
	return &Repository{
		document: client.Doc(document),
	}
}

func (r *Repository) add(collectLogs []collectLog) error {
	ctx := context.Background()
	for _, c := range collectLogs {
		_, _, err := r.document.Collection(c.logType.String()).Add(ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}

type (
	IFetcher interface {
		fetch() ([]collectLog, error)
	}
	Fetcher struct {
		client   *natureremo.Client
		deviceID string
	}
)

// NewFetcher creates Fetcher
func NewFetcher(client *natureremo.Client, deviceID string) IFetcher {
	return &Fetcher{
		client:   client,
		deviceID: deviceID,
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

func (f *Fetcher) fetch() ([]collectLog, error) {
	ctx := context.Background()
	devices, err := f.client.DeviceService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var device *natureremo.Device
	for _, d := range devices {
		if d.ID == f.deviceID {
			device = d
			break
		}
	}
	if device == nil {
		log.Fatalf("not found deviceID: %s", f.deviceID)
	}

	return parseNatureremoDevice(device), nil
}
