package utils

import (
	"fmt"
	"net/http"
	"time"

	"infocenter-service/models"
)

// WriteSSEMessage writes a message in SSE format
func WriteSSEMessage(responseWriter http.ResponseWriter, flusher http.Flusher, message *models.Message) {
	fmt.Fprintf(responseWriter, "id: %d\nevent: msg\ndata: %s\n\n", message.ID, message.Content)
	flusher.Flush()
}

// WriteSSETimeout writes a timeout event in SSE format
func WriteSSETimeout(responseWriter http.ResponseWriter, flusher http.Flusher, messageID int64, duration time.Duration) {
	// Format as required for timeout events
	fmt.Fprintf(responseWriter, "id: %d\nevent: timeout\ndata: 30s\n\n", messageID)
	flusher.Flush()
}