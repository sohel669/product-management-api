<<<<<<< HEAD
# product-management-api
=======
Product Management API with Image Processing and Caching
This project provides an API for managing product data with image processing and caching using Redis, RabbitMQ, and a PostgreSQL database. The API allows users to create products, retrieve product details, and get image processing done asynchronously using RabbitMQ.

Table of Contents
Overview
Technologies Used
Architecture
Setup Instructions
Assumptions
API Endpoints
Contributing
License
Overview
The API allows users to manage product data, including product name, description, price, and images. Images are processed asynchronously using RabbitMQ, and product data is cached in Redis to improve response times. The system is built using the Go programming language and connects to a PostgreSQL database for persistent storage.

Key Features:
Product Management: Add, retrieve, and manage product data.
Image Processing: Store product image URLs and process images asynchronously (e.g., compress and store in S3).
Caching: Use Redis to cache frequently accessed product data for fast retrieval.
Asynchronous Processing: Utilize RabbitMQ to handle image processing jobs asynchronously.
Technologies Used
Go (Golang): Main programming language for backend development.
PostgreSQL: Used as the relational database for storing product data.
Redis: Used to cache product data to reduce database load.
RabbitMQ: Message broker for handling asynchronous image processing jobs.
Docker: For containerizing the application and its dependencies.
Gin/Gorilla Mux: (Optional for better routing) web framework for Go.
Architecture
Components:
Go API Server: Handles HTTP requests for creating and retrieving products.
PostgreSQL Database: Stores product data, including names, descriptions, image URLs, and prices.
Redis Cache: Caches product details to reduce the load on the PostgreSQL database.
RabbitMQ: Acts as a message queue for asynchronous image processing tasks.
Image Processing Microservice (Not implemented here but described in the architecture): A service that consumes image URLs from RabbitMQ, downloads images, compresses them, and stores them in a remote storage solution (like S3). Once processing is done, the database is updated with the compressed image URLs.
Data Flow:
A product is created via the /products/create endpoint.
The product is saved in PostgreSQL, and its details are cached in Redis.
Image URLs are published to a RabbitMQ queue for processing.
The image processing microservice consumes the messages, downloads the images, compresses them, stores them in S3, and updates the product data in the database.
Setup Instructions
Prerequisites:
Go: Install Go (Golang) on your machine. Go Installation Guide
PostgreSQL: Set up a PostgreSQL database instance.
Redis: Install and configure Redis. Redis Installation Guide
RabbitMQ: Install and set up RabbitMQ. RabbitMQ Installation Guide
Steps to run the project:
Clone the repository:

bash
Copy code
git clone https://github.com/yourusername/product-management-api.git
cd product-management-api
Set up environment variables: You can use a .env file to set up environment variables such as database credentials and Redis connection details.

Example .env file:

plaintext
Copy code
DB_USER=postgres
DB_NAME=project
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
REDIS_ADDR=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
Install dependencies: Install Go dependencies:

bash
Copy code
go mod tidy
Start services: Make sure PostgreSQL, Redis, and RabbitMQ are running.

Start PostgreSQL:

bash
Copy code
docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
Start Redis:

bash
Copy code
docker run --name redis -p 6379:6379 -d redis
Start RabbitMQ:

bash
Copy code
docker run --name rabbitmq -p 5672:5672 -p 15672:15672 -d rabbitmq:management
Run the Go server:

bash
Copy code
go run main.go
The server should now be running on http://localhost:8080.

Assumptions
Image URLs: The product_images field in the Product model stores an array of image URLs. The image processing microservice will handle the download and compression of these images.
Message Queue (RabbitMQ): The image URLs are pushed to RabbitMQ for asynchronous processing. The system assumes that there is a microservice listening to the queue and performing the image processing tasks.
Product Creation: The POST /products/create endpoint is responsible for creating new products and handling associated image URLs.
Cache Invalidation: Redis caching is used for performance optimization. Cache invalidation is not implemented in this example, but it should be implemented when product data changes (e.g., when product price or description is updated).
API Endpoints
POST /products/create
Description: Creates a new product with the provided details (name, description, price, and image URLs).
Request Body:
json
Copy code
{
  "user_id": 1,
  "product_name": "Sample Product",
  "product_description": "This is a sample product",
  "product_images": ["image_url_1", "image_url_2"],
  "product_price": 199.99
}
Response:
Status: 201 Created
Body:
json
Copy code
{
  "id": 1,
  "user_id": 1,
  "product_name": "Sample Product",
  "product_description": "This is a sample product",
  "product_images": ["image_url_1", "image_url_2"],
  "product_price": 199.99
}
GET /products
Description: Retrieves all products.
Response:
Status: 200 OK
Body:
json
Copy code
[
  {
    "id": 1,
    "user_id": 1,
    "product_name": "Sample Product",
    "product_description": "This is a sample product",
    "product_images": ["image_url_1", "image_url_2"],
    "product_price": 199.99
  }
]
GET /products/{id}
Description: Retrieves a product by ID.
Response:
Status: 200 OK
Body:
json
Copy code
{
  "id": 1,
  "user_id": 1,
  "product_name": "Sample Product",
  "product_description": "This is a sample product",
  "product_images": ["image_url_1", "image_url_2"],
  "product_price": 199.99
}






7b5acb2 (Initial commit with project files)
