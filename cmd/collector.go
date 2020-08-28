package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/tenntenn/natureremo"
	"github.com/tktkc72/ouchi-dashboard/collector"
)

func handler(w http.ResponseWriter, r *http.Request) {
	err := collect()
	if err != nil {
		log.Printf("%v", err)
		status := http.StatusInternalServerError
		text := http.StatusText(status)
		http.Error(w, fmt.Sprintf("%s", text), status)
	}
}

func collect() error {
	accessToken := os.Getenv("ACCESS_TOKEN")
	deviceID := os.Getenv("DEVICE_ID")
	projectID := os.Getenv("GCP_PROJECT")
	documentPath := os.Getenv("FIRESTORE_DOC_PATH")

	natureremoClient := natureremo.NewClient(accessToken)
	fetcher := collector.NewFetcher(natureremoClient, deviceID)

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer firestoreClient.Close()
	repository := collector.NewRepository(firestoreClient, documentPath)

	service := collector.NewCollectorService(fetcher, repository)
	err = service.Collect()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("collector: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
