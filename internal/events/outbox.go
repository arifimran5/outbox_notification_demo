package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/arifimran5/outbox_notification_demo/internal/database"
	"github.com/arifimran5/outbox_notification_demo/internal/models"
)

// StartOutboxProcessor starts a background goroutine to process events
func StartOutboxProcessor(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Outbox Processor")
				return
			case <-ticker.C:
				processBatch()
			}
		}
	}()
}

func processBatch() {
	// 1. Fetch Pending Events
	// In production, use SELECT ... FOR UPDATE SKIP LOCKED to allow multiple processor instances
	rows, err := database.DB.Query(`
		SELECT id, event_type, payload 
		FROM outbox 
		WHERE status = 'PENDING' 
		ORDER BY created_at ASC 
		LIMIT 10
	`)
	if err != nil {
		log.Printf("Error fetching outbox events: %v", err)
		return
	}
	defer rows.Close()

	var eventsToProcess []models.OutboxEvent
	for rows.Next() {
		var e models.OutboxEvent
		var payloadStr string
		if err := rows.Scan(&e.ID, &e.EventType, &payloadStr); err != nil {
			continue
		}
		e.Payload = json.RawMessage(payloadStr)
		eventsToProcess = append(eventsToProcess, e)
	}

	// 2. Process Each Event
	for _, e := range eventsToProcess {
		if err := handleEvent(e); err != nil {
			log.Printf("Failed to process event %s: %v", e.ID, err)
			// Ideally update status to 'FAILED'
			continue
		}

		// 3. Mark as Processed
		_, err := database.DB.Exec("UPDATE outbox SET status = 'PROCESSED' WHERE id = $1", e.ID)
		if err != nil {
			log.Printf("Failed to update status for event %s: %v", e.ID, err)
		}
	}
}

func handleEvent(e models.OutboxEvent) error {
	if e.EventType == "POST_CREATED" {
		// Determine who needs to be notified
		var payload models.NotificationPayload
		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return err
		}

		// Find all users subscribed to this topic
		// Note: In highly scalable systems, this query might be heavy.
		// You might fan-out to a queue per user or use a dedicated pub/sub system.
		rows, err := database.DB.Query(`
			SELECT s.user_id 
			FROM subscriptions s
			JOIN topics t ON s.topic_id = t.id
			WHERE t.name = $1 AND s.user_id != $2
		`, payload.TopicName, payload.UserID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var userID string
			if err := rows.Scan(&userID); err == nil {

				// Send via SSE
				if userID == payload.UserID {
					continue // Skip sending notification to self
				}
				fmt.Printf("Relaying Outbox Event %s to User %s\n", e.ID, userID)
				Manager.SendNotification(userID, payload)
			}
		}
	}
	return nil
}
