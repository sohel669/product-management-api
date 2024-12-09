package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type Product struct {
	ID                 int            `json:"id" db:"id"`
	UserID             int            `json:"user_id" db:"user_id"`
	ProductName        string         `json:"product_name" db:"product_name"`
	ProductDescription string         `json:"product_description" db:"product_description"`
	ProductImages      pq.StringArray `json:"product_images" db:"product_images"`
	ProductPrice       float64        `json:"product_price" db:"product_price"`
}

var (
	db          *sqlx.DB
	redisClient *redis.Client
	rabbitMQCh  *amqp.Channel
	ctx         = context.Background()
)

func init() {
	var err error

	// Initialize PostgreSQL database
	db, err = sqlx.Connect("postgres", "user=postgres dbname=Project sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	log.Println("Database initialized successfully")

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("Redis initialized successfully")

	// Initialize RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}

	rabbitMQCh, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v\n", err)
	}
	log.Println("RabbitMQ initialized successfully")
}

// CacheProduct caches a product in Redis
func CacheProduct(product Product) {
	cacheKey := "product:" + strconv.Itoa(product.ID)
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshalling product for cache: %v", err)
		return
	}
	err = redisClient.Set(ctx, cacheKey, productJSON, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Error caching product: %v", err)
	}
}

// GetCachedProduct retrieves a product from Redis
func GetCachedProduct(id int) (Product, bool) {
	cacheKey := "product:" + strconv.Itoa(id)
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return Product{}, false
	}

	var product Product
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		log.Printf("Error unmarshalling product from cache: %v", err)
		return Product{}, false
	}

	return product, true
}

func InvalidateCache(id int) {
	cacheKey := "product:" + strconv.Itoa(id)
	err := redisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("Error invalidating cache for product ID %d: %v", id, err)
	}
}

func PublishToQueue(queueName string, message string) {
	err := rabbitMQCh.Publish(
		"",          // exchange
		queueName,   // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish a message to the queue: %v\n", err)
	} else {
		log.Printf("Message published to queue '%s': %s\n", queueName, message)
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Fetching products from the database...")

	var products []Product
	query := "SELECT id, user_id, product_name, product_description, product_images, product_price FROM products"

	err := db.Select(&products, query)
	if err != nil {
		log.Printf("Error fetching products: %v\n", err)
		http.Error(w, "Unable to fetch products", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Printf("Error encoding products to JSON: %v\n", err)
		http.Error(w, "Unable to process response", http.StatusInternalServerError)
		return
	}
	log.Println("Products fetched successfully")
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Creating new product...")

	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("Error decoding product: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = db.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription, product.ProductImages, product.ProductPrice).Scan(&product.ID)
	if err != nil {
		log.Printf("Error inserting product into database: %v\n", err)
		http.Error(w, "Unable to create product", http.StatusInternalServerError)
		return
	}

	for _, imageURL := range product.ProductImages {
		PublishToQueue("image_queue", imageURL)
	}

	CacheProduct(product)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Printf("Error encoding product to JSON: %v\n", err)
		http.Error(w, "Unable to process response", http.StatusInternalServerError)
		return
	}
	log.Printf("Product created successfully: %+v\n", product)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID: %v\n", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Check Redis cache first
	product, found := GetCachedProduct(id)
	if found {
		log.Printf("Product found in cache: %+v\n", product)
		json.NewEncoder(w).Encode(product)
		return
	}

	query := "SELECT id, user_id, product_name, product_description, product_images, product_price FROM products WHERE id = $1"
	err = db.Get(&product, query, id)
	if err != nil {
		log.Printf("Error fetching product by ID: %v\n", err)
		http.Error(w, "Unable to fetch product", http.StatusNotFound)
		return
	}

	CacheProduct(product)

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Printf("Error encoding product to JSON: %v\n", err)
		http.Error(w, "Unable to process response", http.StatusInternalServerError)
		return
	}
	log.Printf("Product with ID %d fetched successfully\n", product.ID)
}

func main() {
	http.HandleFunc("/products", GetProducts)
	http.HandleFunc("/products/create", CreateProduct)
	http.HandleFunc("/products/", GetProductByID)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
