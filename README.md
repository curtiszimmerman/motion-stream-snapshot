# motion stream and snapshot server

A simple Go HTTP server that acts as a proxy for motion snapshots. The server listens for HTTP connections on port 8082 and serves the latest motion snapshot image.

## Features

- HTTP server listening on port 8082 with endpoint `/current` that serves the latest motion snapshot
- motion streaming server listening on port 8081

## Prerequisites

- Docker and Docker Compose (for containerized deployment)
- Video device (/dev/video0) for motion detection

## Configuration

The HTTP snapshot server can be configured using environment variables:

- `SNAPSHOT_URL`: URL for the snapshot endpoint (default: http://motion:8080/00000/action/snapshot)

## Building and Running

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
- Port 8081: Stream / motion detection service
- Port 8082: Access the latest snapshot at `http://localhost:8082/current`

## License

MIT