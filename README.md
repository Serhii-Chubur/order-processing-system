# Order Processing System

A simplified e-commerce order processing system using microservices architecture. The system will handle product catalog, order management, and user authentication through event-driven communication, all containerized with Docker.

## Features

- User registration and authentication
- Product catalog
- Stock management
- Order creation, and management
- Secure API endpoints

## Technologies Used

- Go 1.19+
- PostgreSQL 16+
- Redis
- JWT Authentication
- Docker

### Prerequisites

- Go 1.19+
- PostgreSQL 16+

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/Serhii-Chubur/order-processing-system
    cd order_processing_system
    ```
2. Install dependencies:
    ```bash
    go mod download
    ```
3. Configure your `.env` file according to `.env.sample`.

### Running the Application
1. With Docker
```bash
docker-compose up
```
2. Locally
```bash
go run ./cmd/main.go
```

## Users Credentials
### Admin
- **Email:** admin@admin.com
- **Password:** 000000

### User
- **Email:** test@test.com
- **Password:** 333333
