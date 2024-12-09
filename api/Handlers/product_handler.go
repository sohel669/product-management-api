package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"your_project/models"
	"your_project/queue"
	"your_project/utils"
	"your_project/db" 
)


func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Parse JSON body into product struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input fields
	if product.UserID == 0 || product.ProductName == "" || product.ProductPrice <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}


	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
	err := db.DB.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription, product.ProductImages, product.ProductPrice).Scan(&product.ID)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		log.Println("Error inserting product:", err)
		return
	}

	for _, imageURL := range product.ProductImages {
		err := queue.PublishMessage("image_queue", imageURL)
		if err != nil {
			http.Error(w, "Failed to enqueue image for processing", http.StatusInternalServerError)
			log.Println("Error publishing to RabbitMQ:", err)
			return
		}
	}


	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}


package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"your_project/models"
	"your_project/db"
	"your_project/cache"
	"log"
)

// GetProduct fetches product details by its ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL
	productIDStr := r.URL.Query().Get("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Check if the product is in cache
	product, found := cache.GetProductCache(productID)
	if !found {
		// If not found in cache, query the database
		query := `SELECT id, user_id, product_name, product_description, product_images, product_price, compressed_images, created_at, updated_at
				  FROM products WHERE id = $1`
		err = db.DB.Get(&product, query, productID)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			log.Println("Error fetching product:", err)
			return
		}

		// Cache the product for future requests
		cache.SetProductCache(productID, product)
	}

	// Return the product details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}


package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"your_project/db"
	"your_project/models"
	"log"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")
	productNameFilter := r.URL.Query().Get("product_name")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid min_price", http.StatusBadRequest)
			return
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid max_price", http.StatusBadRequest)
			return
		}
	}


	query := "SELECT id, user_id, product_name, product_description, product_images, product_price, compressed_images, created_at, updated_at FROM products WHERE user_id = $1"
	args := []interface{}{userID}

	if productNameFilter != "" {
		query += " AND product_name ILIKE $2"
		args = append(args, "%"+productNameFilter+"%")
	}

	if minPrice > 0 {
		query += " AND product_price >= $3"
		args = append(args, minPrice)
	}

	if maxPrice > 0 {
		query += " AND product_price <= $4"
		args = append(args, maxPrice)
	}

	var products []models.Product
	err = db.DB.Select(&products, query, args...)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		log.Println("Error fetching products:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

package main

import (
    "log"
    "net/http"
    "your_project/api/routers"
    "your_project/database"
    "your_project/utils"
)

func main() {

    utils.LoadEnv()

  
    database.InitDB()


    router := routers.NewRouter()


    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}


func CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    log.Println("Received a request to create a new product")


    var product models.Product
    err := json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        log.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid product data", http.StatusBadRequest)
        return
    }
    
    log.Printf("Product data: %+v\n", product)


    query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
    var productID int
    err = database.DB.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription, 
        product.ProductImages, product.ProductPrice).Scan(&productID)
    if err != nil {
        log.Printf("Error inserting product into database: %v\n", err)
        http.Error(w, "Failed to create product", http.StatusInternalServerError)
        return
    }


    product.ID = productID
    response, err := json.Marshal(product)
    if err != nil {
        log.Printf("Error marshalling product response: %v\n", err)
        http.Error(w, "Failed to create product response", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write(response)
    log.Println("Product created successfully")
}
