# Product Management API with Image Processing and Caching

This project provides an API for managing product data with image processing and caching using Redis, RabbitMQ, and a PostgreSQL database. The API allows users to create products, retrieve product details, and get image processing done asynchronously using RabbitMQ.

## Table of Contents
- [Overview](#overview)
- [Technologies Used](#technologies-used)
- [Architecture](#architecture)
- [Setup Instructions](#setup-instructions)
- [Assumptions](#assumptions)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Overview
The API allows users to manage product data, including product name, description, price, and images. Images are processed asynchronously using RabbitMQ, and product data is cached in Redis to improve response times. The system is built using the Go programming language and connects to a PostgreSQL database for persistent storage.

### Key Features:
- **Product Management**: Add, retrieve, and manage product data.
- **Image Processing**: Store product image URLs and process images asynchronously (e.g., compress and store in S3).
- **Caching**: Use Redis to cache frequently accessed product data for fast retrieval.
- **Asynchronous Processing**: Utilize RabbitMQ to handle image processing jobs asynchronously.

## Technologies Used
- **Go (Golang)**: Main programming language for backend development.
- **PostgreSQL**: Used as the relational database for storing product data.
- **Redis**: Used to cache product data to reduce database load.
- **RabbitMQ**: Message broker for handling asynchronous image processing tasks.
- **Docker**: For containerizing the application and its dependencies.
- **Gin/Gorilla Mux**: (Optional for better routing) web framework for Go.
  
## Architecture
### Components:
1. **Go API Server**: Handles HTTP requests for creating and retrieving products.
2. **PostgreSQL Database**: Stores product data, including names, descriptions, image URLs, and prices.
3. **Redis Cache**: Caches product details to reduce the load on the PostgreSQL database.
4. **RabbitMQ**: Acts as a message queue for asynchronous image processing tasks.
5. **Image Processing Microservice** (Not implemented here but described in the architecture): A service that consumes image URLs from RabbitMQ, downloads images, compresses them, and stores them in a remote storage solution (like S3). Once processing is done, the database is updated with the compressed image URLs.

### Data Flow:
1. A product is created via the `/products/create` endpoint.
2. The product is saved in PostgreSQL, and its details are cached in Redis.
3. Image URLs are published to a RabbitMQ queue for processing.
4. The image processing microservice consumes the messages, downloads the images, compresses them, stores them in S3, and updates the product data in the database.

## Setup Instructions

### Prerequisites:
- **Go**: Install Go (Golang) on your machine. [Go Installation Guide](https://golang.org/doc/install)
- **PostgreSQL**: Set up a PostgreSQL database instance.
- **Redis**: Install and configure Redis. [Redis Installation Guide](https://redis.io/download)
- **RabbitMQ**: Install and set up RabbitMQ. [RabbitMQ Installation Guide](https://www.rabbitmq.com/download.html)

### Steps to run the project:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/product-management-api.git
   cd product-management-api
