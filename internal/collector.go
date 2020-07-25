package ouchidashboard

//go:generate mkdir -p mock
//go:generate mockgen -source=collector.go -destination=mock/collector.go

import (
	"time"

	"github.com/tenntenn/natureremo"
)

type natureremoClient interface {
	GetAll() ([]*natureremo.Device, error)
}

type Collector struct {
	natureRemoClient natureremoClient
}

type temperatureLog struct {
	temparture int
	createdAt  time.Time
}

type humidityLog struct {
	humidity  int
	createdAt time.Time
}

type illuminationLog struct {
	illumination int
	createdAt    time.Time
}

type CollectLog struct {
	temperatureLog  temperatureLog
	humidityLog     humidityLog
	illuminationLog illuminationLog
}
