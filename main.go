package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handleCurrent(w http.ResponseWriter, r *http.Request) {
	// Get snapshot URL from environment variable
	snapshotURL := os.Getenv("SNAPSHOT_URL")
	if snapshotURL == "" {
		snapshotURL = "http://motion:8080/00000/action/snapshot" // default fallback
	}

	// Make HTTP GET request to the snapshot endpoint
	resp, err := http.Get(snapshotURL)
	if err != nil {
		http.Error(w, "Failed to get snapshot", http.StatusInternalServerError)
		log.Printf("Error getting snapshot: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read the image file
	imageData, err := os.ReadFile("/var/lib/motion/lastsnap.jpg")
	if err != nil {
		http.Error(w, "Failed to read image file", http.StatusInternalServerError)
		log.Printf("Error reading image file: %v", err)
		return
	}

	// Set content type and send the image
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData)
}

func main() {
	http.HandleFunc("/current", handleCurrent)

	port := 8082
	fmt.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
