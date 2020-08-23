package main

import (
	"log"
	"os"

	ouchidashboard "github.com/tktkc72/ouchi-dashboard"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	deviceID := os.Getenv("DEVICE_ID")
	fetcher := ouchidashboard.NewFetcher(accessToken, deviceID)
	collectedLog, err := fetcher.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("temp:%v", collectedLog.TemperatureLog.Value)
}
