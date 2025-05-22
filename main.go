package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func handleCurrent(w http.ResponseWriter, r *http.Request) {
	// Get snapshot URL from environment variable
	snapshotURL := os.Getenv("SNAPSHOT_URL")
	if snapshotURL == "" {
		snapshotURL = "http://localhost:8080/00000/action/snapshot" // default fallback
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make synchronous HTTP GET request to the snapshot endpoint
	resp, err := client.Get(snapshotURL)
	if err != nil {
		http.Error(w, "Failed to get snapshot", http.StatusInternalServerError)
		log.Printf("Error getting snapshot: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Snapshot request failed with status: %d", resp.StatusCode), http.StatusInternalServerError)
		log.Printf("Snapshot request failed with status: %d", resp.StatusCode)
		return
	}

	// Sleep for 500ms to allow motion to write the file
	time.Sleep(100 * time.Millisecond)

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
