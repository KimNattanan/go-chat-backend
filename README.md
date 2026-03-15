# go-chat-backend

Chat backend implemented in Go, following a modular monolith architecture with support for REST, gRPC, WebSocket APIs, and RabbitMQ-based messaging.

## Features

- **Clean Architecture** with clear separation of concerns
- **Access and refresh token** authentication with token rotation
- **REST API** built with Echo
- **gRPC API** support
- **WebSocket API** for real-time communication
- **RabbitMQ** integration for asynchronous messaging
- **PostgreSQL** for persistent data
- **Redis** for refresh token and session management
- **Rate Limiter** implementing Token Bucket algorithm
- **Logging and Recovery** middlewares
- Centralized error mapping and handling

## Prerequisites

- Go 1.26+
- Docker & Docker Compose

## Getting Started

1. Clone the repository:

    ```sh
    git clone https://github.com/KimNattanan/go-chat-backend.git
    cd go-chat-backend
    ```

2. Install Go module dependencies:

    ```sh
    go mod tidy
    ```

3. Configure environment variables

    Copy `.env.example` to `.env` and set the required variables (see `.env.example` for keys).

4. Start the databases using Docker Compose:

    ```sh
    docker-compose up -d
    ```

5. Run the application (database migrations run automatically on startup):

    ```sh
    go run ./cmd/app
    ```

## Project Structure

```
.
├── cmd/app/main.go
├── internal
│   ├── app/
│   ├── auth/
│   ├── profile/
│   ├── chat/
│   │   ├── entity/
│   │   ├── handler/
│   │   │   ├── amqp_rpc/
│   │   │   ├── grpc/
│   │   │   └── rest/
│   │   │       ├── v1/
│   │   │       └── router.go
│   │   ├── proto/
│   │   ├── repo/
│   │   │   ├── persistent/
│   │   │   └── contracts.go
│   │   └── usecase/
│   │       ├── membership/
│   │       ├── room/
│   │       └── contracts.go
│   ├── message/
│   ├── realtime/
│   │   └── handler/
│   │       ├── amqp_rpc/
│   │       └── ws/
|   └── platform/
│       ├── config/
│       ├── middleware
│       │   ├── jwt.go
│       │   ├── logger.go
│       │   └── recovery.go
│       └── wsserver/
├── pkg
│   ├── apperror/
│   ├── grpcserver/
│   ├── httpserver/
│   ├── logger/
│   ├── postgres/
│   ├── rabbitmq/
│   ├── ratelimit/
│   ├── redisclient/
│   ├── responses/
│   └── token/
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── LICENSE
└── README.md
```

## License

This project is licensed under the MIT License.\
See the `LICENSE` file for details.
