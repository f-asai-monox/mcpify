package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Created string `json:"created"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Created: "2024-01-15T10:00:00Z"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Created: "2024-01-20T15:30:00Z"},
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

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Error encoding users", http.StatusInternalServerError)
		}
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		user.ID = len(users) + 1
		user.Created = time.Now().Format(time.RFC3339)
		users = append(users, user)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user *User
	var index int
	for i, u := range users {
		if u.ID == id {
			user = &u
			index = i
			break
		}
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(user)
	case "PUT":
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		updatedUser.ID = id
		updatedUser.Created = user.Created
		users[index] = updatedUser
		json.NewEncoder(w).Encode(updatedUser)
	case "DELETE":
		users = append(users[:index], users[index+1:]...)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", corsHandler(healthHandler))
	http.HandleFunc("/users", corsHandler(usersHandler))
	http.HandleFunc("/users/", corsHandler(userHandler))

	fmt.Printf("Mock API Server starting on port %s...\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  GET    /users")
	fmt.Println("  POST   /users")
	fmt.Println("  GET    /users/{id}")
	fmt.Println("  PUT    /users/{id}")
	fmt.Println("  DELETE /users/{id}")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
