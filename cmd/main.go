package main

import (
	"context"
	"log"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	cli := natureremo.NewClient(accessToken)
	ctx := context.Background()

	devices, err := cli.DeviceService.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range devices {
		log.Printf("%v", d)
	}
}
