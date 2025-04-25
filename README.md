# Pack Optimizer

## Overview

This application solves the pack size optimization problem using dynamic programming and backtracking.  
It determines the optimal combination of packs to fulfill orders while minimizing both the total number of items and the
number of packs used.

## Demo

A live demo is available at [retask.tural.pro](https://retask.tural.pro)  
Watch video demo on [YouTube](https://www.youtube.com/watch?v=qrGtGzoaioM)

![assignment](assignment.png)

## Problem Statement

When a customer orders a specific quantity of items, the system needs to determine the best combination of available
pack sizes to fulfill the order according to these rules:

1. Only complete packs can be sent (no partial packs)
2. The total number of items must be minimized (must be ≥ order amount)
3. The number of packs must be minimized (for the chosen minimum total items)

See the [pdf file](/re-partners-software-challenge.pdf) for the details of task

## Features

### Core Functionality

- **Dynamic Pack Size Management**: Pack sizes are flexible and can be added or removed through API
- **Optimal Solution Algorithm**: Uses dynamic programming with backtracking to find the most efficient pack combination
- **Interactive UI**: User-friendly interface to interact with the API
- **RESTful API**: Clean API endpoints for all operations

### Technical Features

- **TLS Support**: Secure communication over HTTPS
- **Graceful Shutdown**: The server stops accepting new requests while allowing in-flight requests to complete
- **Panic Recovery**: Middleware that catches panics and converts them to proper error responses
- **Health Check Endpoint**: Monitoring endpoint for system health
- **Rate Limiting**: Protection against excessive requests
- **Dynamic CORS Settings**: Configurable trusted origins via command-line flags
- **Metrics Endpoint**: Protected metrics endpoint accessible via `localhost`
- **Database Migrations**: Version-controlled database schema using migration files
- **SQLite Database**: Lightweight, portable database solution

## Technology Stack

- **Backend**: Go 1.24+
- **Database**: SQLite with migration support
- **Router**: [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- **Rate Limiting**: [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)
- **Structured Logging**: log/slog
- **Metrics**: expvar

## API Endpoints

| Method | Endpoint                 | Description                                             |
|--------|--------------------------|---------------------------------------------------------|
| GET    | /                        | User interface                                          |
| GET    | /v1/healthcheck          | Health check endpoint                                   |
| GET    | /v1/pack/calculate/:size | Calculate optimal pack combination for given order size |
| GET    | /v1/pack/size            | Get all available pack sizes                            |
| POST   | /v1/pack/size            | Add a new pack size                                     |
| DELETE | /v1/pack/size/:size      | Delete a pack size                                      |
| GET    | /debug/vars              | Metrics endpoint (localhost only)                       |

## Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite

### Installation

1. Clone the repository
   ```shell
    git clone https://github.com/TuralAsgar/retask.git
   ```
2. Configure the environment:
   ```
   cp .envrc-example .envrc
   ```
3. Run database migrations:
   ```
   make db/migrations/up
   ```
4. Build and run the application:
   ```
   make run/api
   ```

### Configuration

The application can be configured using the following command-line flags:

- `port`: API server port (default: 4004)
- `env`: Environment (development|staging|production)
- `db-dsn`: SQLite database path
- `cors-trusted-origins`: Trusted CORS origins (space separated)
- `limiter-enabled`: Enable rate limiter
- `limiter-rps`: Rate limiter maximum requests per second
- `limiter-burst`: Rate limiter maximum burst

### Deployment

1. Cross-build the application
   ```shell
   make build/api
   ```
2. Deploy
   ```shell
   make production/deploy/api
   ```

## Future Improvements

- [ ] Adding OpenAPI/Swagger documentation for the API
- [ ] Extracting code into more packages as the project grows
- [ ] Adding authentication and authorization
- [ ] Implementing CI/CD pipelines
