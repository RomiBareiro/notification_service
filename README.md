# Notification Service

This project is a notification service implemented in Go, designed to handle rate-limited notifications efficiently.

## Features

- **Rate Limiting**: Implements a rate-limiting mechanism to control the frequency of notifications sent.
- **HTTP Server**: Includes an HTTP server to receive notification requests.
- **Logging**: Utilizes logging to track notification activities and errors.
- **Docker Support**: Provides Dockerfiles for easy deployment in containerized environments.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker (optional, for containerization)
- PostreSql +16

### Installation

1. Clone the repository.
2. Install dependencies with `go mod download`.
3. Set up your environment variables for database connections and API keys.
4. Build and run the service with `go run main.go`.
5. Make POST HTTP requests to `http://localhost:8080/sendNotif` to send notifications.


## Usage
### Running Locally
To run the service locally:

```code
go build -o notification_service main.go
```

```code 
export POSTGRES_HOST=localhost  # use your IP
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=admin
export POSTGRES_DB=notifications
export POSTGRES_SSL_MODE=disable

go run main.go
```

## Using Docker
To run the service using Docker:
### 1- Build the Docker image:
```code 
docker build -t notification-service .
```
### 2- Run the Docker container:
```code
docker run --name notification-service \
  -e POSTGRES_HOST=192.168.1.16 \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=admin \
  -e POSTGRES_DB=notifications \
  -e POSTGRES_SSL_MODE=disable \
  -p 8080:8080 \
  -d notification-service
```

## API Endpoints
The service exposes the following API endpoints:

* POST /sendNotif: Sends a notification. Requires a JSON payload with the notification details.

### Example payload:
Input body
```json 
{
 "recipient": "romi",
 "group" :"NEWS"
}
```

Output:
```json 
{
    "recipient": "romi",
    "message": "success",
    "timestamp": "0001-01-01T00:00:00Z" -- record date
}

```

### What could be improved?
* Authentication (JWT or OAuth or AWS cognito etc)
* Monitoring (Grafana could be a good option)
* API versioning
* Fix docker compose
* Add more restrictions for messages
* Separate by go modules
* Load balancer implementation
