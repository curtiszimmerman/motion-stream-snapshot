# motion stream and snapshot server

A simple Go HTTP server that acts as a proxy for motion snapshots. The server listens for HTTP connections on port 8082 and serves the latest motion snapshot image. This project was made for use with Octoprint so that timelapse creation functions properly with the motion software if, for example, you have hardware that doesn't work with mjpg-streamer. NOTE: The default Octopi image is 32-bit and you WILL NOT be able to use the Docker container functionality! If you wish to use Docker, you will need to install the [64-bit version of OctoPi](https://unofficialpi.org/Distros/OctoPi/octopi-bookworm-arm64-lite-1.1.0.zip) or another 64-bit OS instead of the 32-bit version, and then install Octoprint on top of that base install.

## Features

- HTTP server listening on port 8082
- Endpoint `/current` that serves the latest motion snapshot
- Docker and Docker Compose support with Motion integration

## Prerequisites

- Video device (/dev/video0)
- Go 1.21 or later and the motion package, available from debian/alpine repositories, unless you wanna use...
- Docker and Docker Compose (for containerized deployment)

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

### Local Deployment

1. (Optional) Use the provided installation script:
   ```bash
   ./install.sh
   ```

... or, if you want to do things manually ...

1. Ensure the motion application is installed and running:
   ```bash
   sudo apt-get install motion  # install motion from repo
   motion -b                    # run motion in background (as daemon)
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/curtiszimmerman/motion-stream-snapshot
   cd motion-stream-snapshot
   ```
3. Build the binary:
   ```bash
   make build
   ```
   This will create a binary named `motion-snapshot-server` in the `bin/` directory.

4. Clean build artifacts:
   ```bash
   make clean
   ```

5. Run the application:
   ```bash
   # Using default settings
   ./bin/motion-snapshot-server

   # Using command-line flags
   ./bin/motion-snapshot-server -p 8083 -h example.com -s 8080

   # Using environment variables
   SNAPSHOT_HOST=example.com SNAPSHOT_PORT=8080 ./bin/motion-snapshot-server
   ```
6. (Optional) Copy the systemd service file and enable the service to start services on next boot:
   ```bash
   sudo cp motion-snapshot.service /etc/systemd/system/
   systemctl daemon-reload
   systemctl enable motion-snapshot
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