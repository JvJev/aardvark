package models

import (
	"sync"
	"time"
)

// Message represents a message sent to a topic
type Message struct {
	ID      int64
	Content string
	Topic   string
}

// Topic represents a topic with subscribers and stored messages
type Topic struct {
	Name        string
	Subscribers map[chan *Message]time.Time
	Messages    []*Message
	RWLock      sync.RWMutex
}