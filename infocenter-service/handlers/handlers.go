package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"strconv"

	"github.com/go-chi/chi/v5"

	"infocenter-service/services"
	"infocenter-service/utils"
)

// Handlers contains the HTTP request handlers
type Handlers struct {
	infocenterService *services.InfocenterService
}

// NewHandlers creates a new instance of Handlers
func NewHandlers(infocenterService *services.InfocenterService) *Handlers {
	return &Handlers{
		infocenterService: infocenterService,
	}
}

// ReceiveMessages handles subscribing to a topic
func (handler *Handlers) ReceiveMessages(responseWriter http.ResponseWriter, request *http.Request) {
	topicName := chi.URLParam(request, "topic")
	
	// Set SSE headers
	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")
	
	flusher, isFlusherSupported := responseWriter.(http.Flusher)
	if !isFlusherSupported {
		http.Error(responseWriter, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()
	
	// Check for Last-Event-ID header
	lastEventIDHeader := request.Header.Get("Last-Event-ID")
	var lastEventIDValue int64 = 0
	hasLastEventID := false
	
	if lastEventIDHeader != "" {
		// Parse the Last-Event-ID value
		parsedID, parseError := strconv.ParseInt(lastEventIDHeader, 10, 64)
		if parseError == nil {
			lastEventIDValue = parsedID
			hasLastEventID = true
		}
	}
	
	// Create subscriber and get existing messages
	messageChannel, cleanupFunction := handler.infocenterService.Subscribe(topicName)
	defer cleanupFunction()
	
	// Send existing messages, but filter out already seen ones
	for _, topicMessage := range handler.infocenterService.GetTopicMessages(topicName) {
		// Skip messages that the client has already seen
		if hasLastEventID && topicMessage.ID <= lastEventIDValue {
			continue
		}
		utils.WriteSSEMessage(responseWriter, flusher, topicMessage)
	}
	
	// Set up connection timeout that resets with new messages
	maxConnectionDuration := handler.infocenterService.GetMaxConnectionDuration()
	lastActivityTime := time.Now()
	timeoutChecker := time.NewTicker(1 * time.Second) // Check more frequently
	defer timeoutChecker.Stop()
	
	// Handle messages and timeouts
	for {
		select {
		case <-request.Context().Done():
			return // Client disconnected
			
		case <-timeoutChecker.C:
			// Check if timeout has been reached since last activity
			elapsedTime := time.Since(lastActivityTime)
			if elapsedTime >= maxConnectionDuration {
				// Send timeout event before disconnecting - DIRECTLY WRITE THE EVENT
				currentMessageID := handler.infocenterService.GetMessageIDSeq()
				fmt.Fprintf(responseWriter, "id: %d\nevent: timeout\ndata: 30s\n\n", currentMessageID)
				flusher.Flush()
				return // Disconnect after sending timeout
			}
			
		case receivedMessage := <-messageChannel:
			if receivedMessage == nil {
				return // Channel closed
			}
			utils.WriteSSEMessage(responseWriter, flusher, receivedMessage)
			lastActivityTime = time.Now() // Reset timeout countdown on new message
		}
	}
}

// SendMessage handles sending messages to a topic
func (handler *Handlers) SendMessage(responseWriter http.ResponseWriter, request *http.Request) {
	topicName := chi.URLParam(request, "topic")
	
	// Read message content
	messageBodyBytes, readError := io.ReadAll(request.Body)
	if readError != nil {
		http.Error(responseWriter, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Add message and return success
	handler.infocenterService.AddMessage(topicName, string(messageBodyBytes))
	responseWriter.WriteHeader(http.StatusNoContent)
}