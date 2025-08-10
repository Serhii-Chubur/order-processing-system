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

## API Endpoints
### Product Service (Port: 8001)

- GET /api/products - List all products
- GET /api/products/{id} - Get product by ID
- POST /api/products - Create product (admin)
- PUT /api/products/{id} - Update product (admin)
- DELETE /api/products/{id} - Delete product (admin)
- GET /api/products/{id}/stock - Get product stock level

### Order Service (Port: 8002)

- POST /api/orders - Create new order
- GET /api/orders/{id} - Get order by ID
- GET /api/orders/user/{id} - Get orders by user ID
- PUT /api/orders/{id}/status - Update order status

### User Service (Port: 8003)

- POST /api/users/register - Register new user
- POST /api/users/login - User login
- GET /api/users/{id} - Get user profile
- PUT /api/users/{id} - Update user profile (name)
- POST /api/users/logout - User logout
