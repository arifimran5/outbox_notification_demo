package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Topic struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Subscription struct {
	UserID  string `json:"user_id"`
	TopicID string `json:"topic_id"`
}

type Post struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type OutboxEvent struct {
	ID            string          `json:"id"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

// NotificationPayload is what gets sent via SSE
type NotificationPayload struct {
	Message   string `json:"message"`
	TopicName string `json:"topic_name"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
}
