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
		add(CollectLog) error
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

func (r *Repository) add(collected CollectLog) error {
	ctx := context.Background()
	_, _, err := r.document.Collection("temperature").Add(ctx, collected.temperatureLog)
	if err != nil {
		return err
	}
	_, _, err = r.document.Collection("humidity").Add(ctx, collected.humidityLog)
	if err != nil {
		return err
	}
	_, _, err = r.document.Collection("illumination").Add(ctx, collected.illuminationLog)
	if err != nil {
		return err
	}
	_, _, err = r.document.Collection("motion").Add(ctx, collected.motionLog)
	if err != nil {
		return err
	}

	return nil
}

type (
	IFetcher interface {
		fetch() (CollectLog, error)
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

func (f *Fetcher) fetch() (CollectLog, error) {
	ctx := context.Background()
	devices, err := f.client.DeviceService.GetAll(ctx)
	if err != nil {
		return CollectLog{}, err
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

	return CollectLog{
		historyLog{
			device.NewestEvents[natureremo.SensorTypeTemperature].Value,
			device.NewestEvents[natureremo.SensorTypeTemperature].CreatedAt,
		},
		historyLog{
			device.NewestEvents[natureremo.SensorTypeHumidity].Value,
			device.NewestEvents[natureremo.SensorTypeHumidity].CreatedAt,
		},
		historyLog{
			device.NewestEvents[natureremo.SensortypeIllumination].Value,
			device.NewestEvents[natureremo.SensortypeIllumination].CreatedAt,
		},
		historyLog{
			device.NewestEvents[natureremo.SensorType("mo")].Value,
			device.NewestEvents[natureremo.SensorType("mo")].CreatedAt,
		},
	}, nil
}

type historyLog struct {
	Value     float64
	UpdatedAt time.Time
}

type CollectLog struct {
	temperatureLog  historyLog
	humidityLog     historyLog
	illuminationLog historyLog
	motionLog       historyLog
}
