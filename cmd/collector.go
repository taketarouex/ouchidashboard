package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/tenntenn/natureremo"
	"github.com/tktkc72/ouchi-dashboard/collector"
)

type message struct {
	DeviceIDs []string `json:"deviceIDs"`
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := os.Getenv("NATURE_REMO_ACCESS_TOKEN")
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")

	var m message
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(b, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	errorChannel := make(chan error, len(m.DeviceIDs))
	for _, deviceID := range m.DeviceIDs {
		go collect(accessToken, deviceID, projectID, rootPath, errorChannel)
	}
	for range m.DeviceIDs {
		err := <-errorChannel
		if err != nil {
			log.Printf("collect: %v", err)
			if collector.IsNoDevice(err) {
				http.Error(w,
					"Bad Request",
					http.StatusBadRequest)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}

func collect(accessToken, deviceID, projectID, rootPath string, c chan error) {
	natureremoClient := natureremo.NewClient(accessToken)
	fetcher := collector.NewFetcher(natureremoClient, deviceID)

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		c <- err
		return
	}
	defer firestoreClient.Close()
	repository, err := collector.NewRepository(firestoreClient, rootPath, deviceID)
	if err != nil {
		c <- err
		return
	}

	service := collector.NewCollectorService(fetcher, repository)
	err = service.Collect()
	if err != nil {
		c <- err
		return
	}
	c <- nil
}

func main() {
	http.HandleFunc("/", collectorHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("collector: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
