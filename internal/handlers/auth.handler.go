package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/arifimran5/outbox_notification_demo/internal/database"
	"github.com/arifimran5/outbox_notification_demo/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Fetch User (ID and Hashed Password)
	var id string
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", req.Username).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 2. Compare Hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Ensure we use the same key source as the middleware
	// Note: In a real app, ensure this ENV var is set.
	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "user_id": id})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Hash Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	hashedPassword := string(hashedBytes)

	// 2. Insert User
	var id string
	err = database.DB.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", req.Username, hashedPassword).Scan(&id)
	if err != nil {
		// Postgres error code 23505 is unique_violation
		http.Error(w, "Username likely taken", http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
