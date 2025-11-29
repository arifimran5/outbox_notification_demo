package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// SSEManager handles active connections
type SSEManager struct {
	// Map of UserID -> List of ResponseWriters (allows multiple tabs per user)
	clients map[string]map[chan string]bool
	mu      sync.RWMutex
}

var Manager *SSEManager

func InitSSE() {
	Manager = &SSEManager{
		clients: make(map[string]map[chan string]bool),
	}
}

// Subscribe adds a client to the manager
func (m *SSEManager) Subscribe(userID string, w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	if _, ok := m.clients[userID]; !ok {
		m.clients[userID] = make(map[chan string]bool)
	}

	// Create a channel for this specific connection
	msgChan := make(chan string)
	m.clients[userID][msgChan] = true
	m.mu.Unlock()

	// Headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Flush immediately to establish connection
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	log.Printf("User %s connected to SSE", userID)

	// Listen for messages or context cancellation (client disconnect)
	for {
		select {
		case msg := <-msgChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			m.mu.Lock()
			delete(m.clients[userID], msgChan)
			if len(m.clients[userID]) == 0 {
				delete(m.clients, userID)
			}
			m.mu.Unlock()
			log.Printf("User %s disconnected", userID)
			return

		}
	}
}

func (m *SSEManager) SendNotification(userID string, payload any) {
	m.mu.RLock()
	userChans, ok := m.clients[userID]
	if !ok {
		m.mu.RUnlock()
		return
	}

	// Keep lock during iteration
	data, err := json.Marshal(payload)
	m.mu.RUnlock()

	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}
	msg := string(data)

	// Send without holding lock
	m.mu.RLock()
	for ch := range userChans {
		select {
		case ch <- msg:
		default:
		}
	}
	m.mu.RUnlock()
}
