package main

import (
	"os"

	"github.com/tktkc72/ouchi-dashboard/collector"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	deviceID := os.Getenv("DEVICE_ID")
	fetcher := collector.NewFetcher(accessToken, deviceID)
	repository := collector.NewRepository()

	service := collector.NewCollectorService(fetcher, repository)
	service.Collect()
}
