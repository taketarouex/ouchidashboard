package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tktkc72/ouchi-dashboard/collector"
)

func main() {
	http.HandleFunc("/", collector.CollectorHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("collector: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
