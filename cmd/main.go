package main

import (
	"log"
	"os"

	ouchidashboard "github.com/tktkc72/ouchi-dashboard"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	deviceID := os.Getenv("DEVICE_ID")
	collector := ouchidashboard.NewNatureClient(accessToken, deviceID)
	collectedLog, err := collector.CollectLog()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("temp:%v", collectedLog.TemperatureLog.Value)
}
