# motion stream and snapshot server

A simple Go HTTP server that acts as a proxy for motion snapshots. The server listens for HTTP connections on port 8082 and serves the latest motion snapshot image. This project was made for use with Octoprint so that timelapse creation functions properly with the motion software. NOTE: The Octopi image is 32-bit and you WILL NOT be able to use the Docker container functionality! If you wish to use Docker, you will need to install a 64-bit version of your OS (e.g. Raspberry Pi OS 64-bit) instead of the 32-bit version, and then install Octoprint on top of that base install.

## Features

- HTTP server listening on port 8082
- Endpoint `/current` that serves the latest motion snapshot
- Docker and Docker Compose support with Motion integration

## Prerequisites

- Go 1.21 or later
- The motion package, available from debian/alpine repositories
- Docker and Docker Compose (for containerized deployment)
- Video device (/dev/video0) for motion detection

## Configuration

The application can be configured using environment variables or command-line flags:

### Environment Variables
- `SNAPSHOT_HOST`: Host for motion snapshot endpoint (default: localhost)
- `SNAPSHOT_PORT`: Port for motion snapshot endpoint (default: 8080)

### Command-line Flags
```
-p, --port int         Port to listen on for /current endpoint (default 8082)
-h, --snapshot-host string  Host for motion snapshot endpoint (default "localhost")
-s, --snapshot-port int    Port for motion snapshot endpoint (default 8080)
```

Command-line flags take precedence over environment variables.

## Building and Running

### Local Development

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd motion-stream-snapshot
   ```

2. Build the binary:
   ```bash
   make build
   ```
   This will create a binary named `motion-snapshot-server` in the current directory.

3. Clean build artifacts:
   ```bash
   make clean
   ```

4. Run the application:
   ```bash
   # Using default settings
   ./motion-snapshot-server

   # Using command-line flags
   ./motion-snapshot-server -p 8083 -h example.com -s 8080

   # Using environment variables
   SNAPSHOT_HOST=example.com SNAPSHOT_PORT=8080 ./motion-snapshot-server
   ```

### Docker Deployment

1. Build and start the containers:
   ```bash
   docker-compose up -d
   ```

2. Stop the containers:
   ```bash
   docker-compose down
   ```

3. View logs:
   ```bash
   docker-compose logs -f
   ```

## Usage

The server exposes two ports:
- Port 8081: Motion stream service
- Port 8082: Access the latest snapshot at `http://localhost:8082/current`

The `/current` endpoint will:
1. Make a GET request to the configured snapshot URL
2. Read the latest snapshot from `/var/lib/motion/lastsnap.jpg`
3. Return the image with the appropriate content type

## Volume Mounting

The application uses the following volume mounts when running in Docker:

1. `/var/lib/motion:/var/lib/motion`: Shared volume for motion snapshots
2. `/dev/video0:/dev/video0:ro`: Read-only access to the video device (motion service only)

## License

MIT