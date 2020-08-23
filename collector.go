//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=$GOPACKAGE

package ouchidashboard

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

type (
	IRepository interface {
		add(collectedLog) error
	}
	Repository struct {
	}
)

func NewRepository() IRepository {
	return &Repository{}
}

func (*Repository) add(collected collectedLog) error {
	return nil
}

type (
	IFetcher interface {
		fetch() (collectedLog, error)
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

func (c *Fetcher) fetch() (collectedLog, error) {
	ctx := context.Background()
	devices, err := c.client.DeviceService.GetAll(ctx)
	if err != nil {
		return collectedLog{}, err
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

	return collectedLog{
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

type collectedLog struct {
	temperatureLog  historyLog
	humidityLog     historyLog
	illuminationLog historyLog
	motionLog       historyLog
}
