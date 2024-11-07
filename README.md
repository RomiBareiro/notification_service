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
- Docker
- PostreSql +16

### Installation

1. Clone the repository.
2. Install dependencies with `go mod download`.
3. Set up your environment variables for database connections and API keys.
4. Build and run the service with `go run main.go`.
5. Make POST HTTP requests to `http://localhost:8080/V1/notify` to send notifications.


## Usage
### Using Docker compose

```code
docker compose up --build
```

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

## API Endpoints
The service exposes the following API endpoints:

1. POST /login: Get auth token. It must be refreshed every 5 min. *keys: username, password*

2. POST /V1/notify: Sends a notification. Requires a JSON payload with the notification details.

## Message types
Currently, we have:
'''
    ('STATUS', 2, EXTRACT(EPOCH FROM INTERVAL '1 minute')),
    ('NEWS', 1, EXTRACT(EPOCH FROM INTERVAL '1 day')),
    ('MARKETING', 3, EXTRACT(EPOCH FROM INTERVAL '1 hour'));
'''

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
* Monitoring (Grafana could be a good option)
* Add more restrictions for messages
* Separate by go modules
* Load balancer implementation

