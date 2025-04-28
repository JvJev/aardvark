package services

import (
	"sync"
	"sync/atomic"
	"time"

	"infocenter-service/models"
)

// InfocenterService manages topics and messages
type InfocenterService struct {
	topics              map[string]*models.Topic
	rwLock              sync.RWMutex
	messageIDSequence   int64
	maxConnectionTime   time.Duration
}

// NewInfocenterService creates a new instance of InfocenterService
func NewInfocenterService(maxConnectionTime time.Duration) *InfocenterService {
	return &InfocenterService{
		topics:              make(map[string]*models.Topic),
		maxConnectionTime:   maxConnectionTime,
	}
}

// GetOrCreateTopic gets or creates a topic
func (service *InfocenterService) GetOrCreateTopic(topicName string) *models.Topic {
	// Try read lock first for better performance
	service.rwLock.RLock()
	topic, topicExists := service.topics[topicName]
	service.rwLock.RUnlock()

	if topicExists {
		return topic
	}

	// Not found, acquire write lock and create
	service.rwLock.Lock()
	defer service.rwLock.Unlock()
	
	// Double-check after acquiring write lock
	if topic, topicExists = service.topics[topicName]; topicExists {
		return topic
	}
	
	// Create new topic
	topic = &models.Topic{
		Name:        topicName,
		Subscribers: make(map[chan *models.Message]time.Time),
		Messages:    make([]*models.Message, 0),
	}
	service.topics[topicName] = topic
	return topic
}

// CleanupTopic removes a topic if it has no subscribers
func (service *InfocenterService) CleanupTopic(topicName string) {
	service.rwLock.Lock()
	defer service.rwLock.Unlock()
	
	topic, topicExists := service.topics[topicName]
	if !topicExists {
		return
	}
	
	topic.RWLock.RLock()
	hasNoSubscribers := len(topic.Subscribers) == 0
	topic.RWLock.RUnlock()
	
	if hasNoSubscribers {
		delete(service.topics, topicName)
	}
}

// AddMessage adds a message to a topic
func (service *InfocenterService) AddMessage(topicName, messageContent string) *models.Message {
	// Create message with next ID
	messageID := atomic.AddInt64(&service.messageIDSequence, 1)
	newMessage := &models.Message{
		ID:      messageID,
		Content: messageContent,
		Topic:   topicName,
	}
	
	topic := service.GetOrCreateTopic(topicName)
	
	// Store and broadcast the message
	topic.RWLock.Lock()
	topic.Messages = append(topic.Messages, newMessage)
	activeSubscribers := make([]chan *models.Message, 0, len(topic.Subscribers))
	for messageChannel := range topic.Subscribers {
		activeSubscribers = append(activeSubscribers, messageChannel)
	}
	topic.RWLock.Unlock()
	
	// Broadcast to all subscribers (without holding the lock)
	for _, messageChannel := range activeSubscribers {
		select {
		case messageChannel <- newMessage:
		default: // If channel is full (100 messages) skip this subscriber
		}
	}
	
	return newMessage
}

// Subscribe adds a subscriber to a topic and returns a channel for messages
func (service *InfocenterService) Subscribe(topicName string) (chan *models.Message, func()) {
	topic := service.GetOrCreateTopic(topicName)
	messageChannel := make(chan *models.Message, 100)
	
	// Adding subscriber
	topic.RWLock.Lock()
	topic.Subscribers[messageChannel] = time.Now()
	topic.RWLock.Unlock()
	
	// Return cleanup function
	cleanupFunction := func() {
		topic.RWLock.Lock()
		delete(topic.Subscribers, messageChannel)
		topic.RWLock.Unlock()
		close(messageChannel)
		service.CleanupTopic(topicName)
	}
	
	return messageChannel, cleanupFunction
}

// GetTopicMessages returns all messages for a topic
func (service *InfocenterService) GetTopicMessages(topicName string) []*models.Message {
	topic := service.GetOrCreateTopic(topicName)
	
	topic.RWLock.RLock()
	defer topic.RWLock.RUnlock()
	
	// Make a copy to avoid race conditions
	topicMessages := make([]*models.Message, len(topic.Messages))
	copy(topicMessages, topic.Messages)
	
	return topicMessages
}

// GetMaxConnectionDuration returns the maximum connection duration
func (service *InfocenterService) GetMaxConnectionDuration() time.Duration {
	return service.maxConnectionTime
}

// GetMessageIDSeq returns the current message ID sequence
func (service *InfocenterService) GetMessageIDSeq() int64 {
	return atomic.LoadInt64(&service.messageIDSequence)
}