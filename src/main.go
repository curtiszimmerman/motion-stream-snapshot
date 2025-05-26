package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	port         = flag.Int("port", 8082, "Port to listen on for /current endpoint")
	snapshotHost = flag.String("snapshot-host", "localhost", "Host for motion snapshot endpoint (default localhost, overridden by SNAPSHOT_HOST env var)")
	snapshotPort = flag.Int("snapshot-port", 8080, "Port for motion snapshot endpoint (default 8080, overridden by SNAPSHOT_PORT env var)")
	daemon       = flag.Bool("daemon", false, "Run as daemon (background process)")
)

func init() {
	// Add single-letter options
	flag.IntVar(port, "p", 8082, "Port to listen on for /current endpoint")
	flag.StringVar(snapshotHost, "h", "localhost", "Host for motion snapshot endpoint (default localhost, overridden by SNAPSHOT_HOST env var)")
	flag.IntVar(snapshotPort, "s", 8080, "Port for motion snapshot endpoint (default 8080, overridden by SNAPSHOT_PORT env var)")
	flag.BoolVar(daemon, "d", false, "Run as daemon (background process)")
}

func handleCurrent(w http.ResponseWriter, r *http.Request) {
	// Get snapshot URL from flag or environment variable
	host := *snapshotHost
	if os.Getenv("SNAPSHOT_HOST") != "" {
		host = os.Getenv("SNAPSHOT_HOST")
	}
	port := *snapshotPort
	if os.Getenv("SNAPSHOT_PORT") != "" {
		var err error
		port, err = strconv.Atoi(os.Getenv("SNAPSHOT_PORT"))
		if err != nil {
			log.Fatalf("Invalid SNAPSHOT_PORT: %v", err)
		}
	}

	url := fmt.Sprintf("http://%s:%d/00000/action/snapshot", host, port)
	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make synchronous HTTP GET request to the snapshot endpoint
	resp, err := client.Get(url)
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

	// Sleep for 100ms to allow motion to write the file
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
	// Parse command line flags
	flag.Parse()

	// Set up logging
	logFile := "/var/log/motion-snapshot-server.log"
	if *daemon {
		// Open log file
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	http.HandleFunc("/current", handleCurrent)

	// Create server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	if *daemon {
		log.Printf("Server starting in daemon mode on port %d...\n", *port)
	} else {
		fmt.Printf("Server starting on port %d...\n", *port)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
