package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/arifimran5/outbox_notification_demo/internal/database"
	"github.com/arifimran5/outbox_notification_demo/internal/middleware"
	"github.com/arifimran5/outbox_notification_demo/internal/models"
	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var req models.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start Transaction (The Core of Outbox Pattern)
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Insert Post
	var postID string
	var topicName string

	// Get Topic Name for notification
	err = tx.QueryRow("SELECT name FROM topics WHERE id = $1", req.TopicID).Scan(&topicName)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusBadRequest)
		return
	}

	err = tx.QueryRow(`
		INSERT INTO posts (topic_id, author_id, title, content) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`, req.TopicID, userID, req.Title, req.Content).Scan(&postID)

	if err != nil {
		http.Error(w, "Failed to save post", http.StatusInternalServerError)
		return
	}

	// 2. Insert Outbox Event (Same Transaction)
	payload := models.NotificationPayload{
		Message:   "New post in " + topicName + ": " + req.Title,
		TopicName: topicName,
		PostID:    postID,
		UserID:    userID,
	}
	payloadBytes, _ := json.Marshal(payload)

	_, err = tx.Exec(`
		INSERT INTO outbox (aggregate_id, aggregate_type, event_type, payload)
		VALUES ($1, 'POST', 'POST_CREATED', $2)
	`, postID, payloadBytes)

	if err != nil {
		http.Error(w, "Failed to save event", http.StatusInternalServerError)
		return
	}

	// Commit Transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "Post created and event queued"})
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DB.Query("SELECT id, name, description FROM topics")
	defer rows.Close()
	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		rows.Scan(&t.ID, &t.Name, &t.Description)
		topics = append(topics, t)
	}
	json.NewEncoder(w).Encode(topics)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	vars := mux.Vars(r)
	topicID := vars["topicId"]

	_, err := database.DB.Exec("INSERT INTO subscriptions (user_id, topic_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	vars := mux.Vars(r)
	topicID := vars["topicId"]

	_, err := database.DB.Exec("DELETE FROM subscriptions WHERE user_id=$1 AND topic_id=$2", userID, topicID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	rows, _ := database.DB.Query(`
		SELECT t.id, t.name, t.description 
		FROM subscriptions s 
		JOIN topics t ON s.topic_id = t.id 
		WHERE s.user_id = $1`, userID)
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		rows.Scan(&t.ID, &t.Name, &t.Description)
		topics = append(topics, t)
	}
	json.NewEncoder(w).Encode(topics)
}

func GetTopicPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]

	rows, _ := database.DB.Query("SELECT id, title, content, created_at FROM posts WHERE topic_id = $1 ORDER BY created_at DESC", topicID)
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
		posts = append(posts, p)
	}
	json.NewEncoder(w).Encode(posts)
}
