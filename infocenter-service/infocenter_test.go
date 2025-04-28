package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	
	"infocenter-service/handlers"
	"infocenter-service/services"
)

// Test message IDs and content
func TestMessaging(t *testing.T) {
	infoService := services.NewInfocenterService(30 * time.Second)
	
	firstMessage := infoService.AddMessage("test-topic", "Message 1")
	secondMessage := infoService.AddMessage("test-topic", "Message 2")
	
	if firstMessage.ID != 1 || secondMessage.ID != 2 {
		t.Errorf("Message IDs not incrementing correctly: got %d and %d", firstMessage.ID, secondMessage.ID)
	}
	
	if firstMessage.Content != "Message 1" || firstMessage.Topic != "test-topic" {
		t.Errorf("Message content or topic incorrect")
	}
	
	infoService.AddMessage("topic1", "First")
	infoService.AddMessage("topic2", "Second")
	infoService.AddMessage("topic1", "Third")
	
	if len(infoService.GetTopicMessages("topic1")) != 2 || len(infoService.GetTopicMessages("topic2")) != 1 {
		t.Errorf("Topic message counts incorrect")
	}
}

// Test Last-Event-ID handling
func TestLastEventID(t *testing.T) {
	infoService := services.NewInfocenterService(30 * time.Second)
	infoService.AddMessage("test-topic", "Message 1")
	infoService.AddMessage("test-topic", "Message 2")
	
	request, _ := http.NewRequest("GET", "/infocenter/test-topic", nil)
	request.Header.Set("Last-Event-ID", "1")
	responseRecorder := httptest.NewRecorder()
	
	router := chi.NewRouter()
	handler := handlers.NewHandlers(infoService)
	router.Get("/infocenter/{topic}", handler.ReceiveMessages)
	
	go func() {
		router.ServeHTTP(responseRecorder, request)
	}()
	
	time.Sleep(100 * time.Millisecond)
	
	if responseRecorder.Code != http.StatusOK || strings.Contains(responseRecorder.Body.String(), "data: Message 1") {
		t.Errorf("Last-Event-ID not properly handled")
	}
}