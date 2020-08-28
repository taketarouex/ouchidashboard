package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/tenntenn/natureremo"
	"github.com/tktkc72/ouchi-dashboard/collector"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	deviceID := os.Getenv("DEVICE_ID")
	projectID := os.Getenv("GCP_PROJECT")
	documentPath := os.Getenv("FIRESTORE_DOC_PATH")

	natureremoClient := natureremo.NewClient(accessToken)
	fetcher := collector.NewFetcher(natureremoClient, deviceID)

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed get firestore client due to: %v", err)
	}
	defer firestoreClient.Close()
	repository := collector.NewRepository(firestoreClient, documentPath)

	service := collector.NewCollectorService(fetcher, repository)
	err = service.Collect()
	if err != nil {
		log.Fatalf("failed collect log due to: %v", err)
	}
}
