//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=$GOPACKAGE

package ouchidashboard

import (
	"context"
	"log"
	"time"

	"github.com/tenntenn/natureremo"
)

type (
	// Collector collect log
	Collector interface {
		CollectLog() (CollectedLog, error)
	}
	// NatureClient nature remo
	NatureClient struct {
		client    *natureremo.Client
		deviceID  string
		Collector Collector
	}
)

func NewNatureClient(accessToken, deviceID string) *NatureClient {
	cli := natureremo.NewClient(accessToken)
	return &NatureClient{
		client:   cli,
		deviceID: deviceID,
	}
}

func (c *NatureClient) CollectLog() (CollectedLog, error) {
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
		Log{
			device.NewestEvents[natureremo.SensorTypeTemperature].Value,
			device.NewestEvents[natureremo.SensorTypeTemperature].CreatedAt,
		},
		Log{
			device.NewestEvents[natureremo.SensorTypeHumidity].Value,
			device.NewestEvents[natureremo.SensorTypeHumidity].CreatedAt,
		},
		Log{
			device.NewestEvents[natureremo.SensortypeIllumination].Value,
			device.NewestEvents[natureremo.SensortypeIllumination].CreatedAt,
		},
		Log{
			device.NewestEvents[natureremo.SensorType("mo")].Value,
			device.NewestEvents[natureremo.SensorType("mo")].CreatedAt,
		},
	}, nil
}

type Log struct {
	Value     float64
	UpdatedAt time.Time
}

// CollectedLog collected log
type CollectedLog struct {
	TemperatureLog  Log
	HumidityLog     Log
	IlluminationLog Log
	MotionLog       Log
}
