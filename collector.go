//go:generate mockgen -source=$GOFILE -destination=collector_mock.go -package=$GOPACKAGE -self_package=github.com/tktkc72/ouchi-dashboard

package collector

import (
	"context"
	"log"
	"time"

	"github.com/tenntenn/natureremo"
)

type (
	ICollector interface {
		Collect() error
	}
	// CollectorSevice service
	CollectorSevice struct {
		fetcher    Fetcher
		repository Repository
	}
)

func NewCollectorService(fetcher IFetcher, repository IRepository) ICollector {
	return &CollectorSevice{}
}

func (*CollectorSevice) Collect() error {
	return nil
}

type (
	IRepository interface {
		add(CollectedLog) error
	}
	Repository struct {
	}
)

func NewRepository() IRepository {
	return &Repository{}
}

func (*Repository) add(collected CollectedLog) error {
	return nil
}

type (
	IFetcher interface {
		fetch() (CollectedLog, error)
	}
	Fetcher struct {
		client   *natureremo.Client
		deviceID string
	}
)

// NewFetcher creates Fetcher
func NewFetcher(accessToken, deviceID string) IFetcher {
	cli := natureremo.NewClient(accessToken)
	return &Fetcher{
		client:   cli,
		deviceID: deviceID,
	}
}

func (c *Fetcher) fetch() (CollectedLog, error) {
	ctx := context.Background()
	devices, err := c.client.DeviceService.GetAll(ctx)
	if err != nil {
		return CollectedLog{}, err
	}

	var device *natureremo.Device
	for _, d := range devices {
		if d.ID == c.deviceID {
			device = d
			break
		}
	}
	if device == nil {
		log.Fatalf("not found deviceID: %s", c.deviceID)
	}

	return CollectedLog{
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

type CollectedLog struct {
	temperatureLog  historyLog
	humidityLog     historyLog
	illuminationLog historyLog
	motionLog       historyLog
}
