package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Created     string  `json:"created"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Description: "High-performance laptop", Price: 999.99, Category: "Electronics", Created: "2024-01-15T10:00:00Z"},
	{ID: 2, Name: "Smartphone", Description: "Latest smartphone", Price: 699.99, Category: "Electronics", Created: "2024-01-20T15:30:00Z"},
	{ID: 3, Name: "Coffee Mug", Description: "Ceramic coffee mug", Price: 12.99, Category: "Home", Created: "2024-01-25T09:15:00Z"},
}

func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if err := json.NewEncoder(w).Encode(products); err != nil {
			http.Error(w, "Error encoding products", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/products/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product *Product
	for _, p := range products {
		if p.ID == id {
			product = &p
			break
		}
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(product)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Get port from environment variable, default to 8082
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/products", corsHandler(productsHandler))
	http.HandleFunc("/products/", corsHandler(productHandler))

	fmt.Printf("Mock Products API Server starting on port %s...\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /products")
	fmt.Println("  GET    /products/{id}")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
