import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import TopicList from './components/TopicList';
import MessageList from './components/MessageList';
import MessageInput from './components/MessageInput';

interface Message {
  id: number;
  content: string;
  topic: string;
  event?: string;
}

const App: React.FC = () => {
  // Initialize with empty array, no localStorage
  const [topics, setTopics] = useState<string[]>([]);
  
  const [currentTopic, setCurrentTopic] = useState<string>('');
  // Store messages by topic
  const [messagesByTopic, setMessagesByTopic] = useState<Record<string, Message[]>>({});
  const [connected, setConnected] = useState<boolean>(false);
  const eventSourceRef = useRef<EventSource | null>(null);
  
  // Create a function to get current messages
  const getCurrentMessages = () => {
    return currentTopic ? messagesByTopic[currentTopic] || [] : [];
  };

  // Join or connect to a topic
  const joinTopic = (topic: string) => {
    if (topic.trim() === '') return;
    
    // Close existing connection
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      setConnected(false);
    }
    
    setCurrentTopic(topic);
    
    // Add to topics list if not already there
    if (!topics.includes(topic)) {
      setTopics([...topics, topic]);
      // Initialize empty message array for new topic
      setMessagesByTopic(prev => ({
        ...prev,
        [topic]: []
      }));
    }
    
    console.log(`Connecting to topic: ${topic}`);
    
    // Define connect function
    const connectEventSource = () => {
      try {
        console.log(`Starting EventSource connection to ${topic}...`);
        const eventSource = new EventSource(`http://localhost:8080/infocenter/${topic}`);
        eventSourceRef.current = eventSource;
        
        let openReceived = false;
        
        eventSource.onopen = () => {
          console.log(`Successfully connected to topic: ${topic}`);
          openReceived = true;
          setConnected(true);
        };
        
        setTimeout(() => {
          if (!openReceived && eventSourceRef.current === eventSource) {
            console.log("No onopen event received, but connection seems to be working");
            setConnected(true);
          }
        }, 1000);
        
        eventSource.addEventListener('msg', (event) => {
          console.log('Received message:', event);
          const content = event.data;
          const id = event.lastEventId;
          
          if (!openReceived) {
            setConnected(true);
            openReceived = true;
          }
          
          const newMessage: Message = {
            id: parseInt(id || '0'),
            content,
            topic,
            event: 'msg'
          };
          
          // Add message to the specific topic
          setMessagesByTopic(prev => ({
            ...prev,
            [topic]: [...(prev[topic] || []), newMessage]
          }));
        });
        
        eventSource.addEventListener('timeout', (event) => {
          console.log(`Connection timeout after ${event.data}`);
          
          const timeoutMessage: Message = {
            id: parseInt(event.lastEventId || '0'),
            content: event.data,
            topic,
            event: 'timeout'
          };
          
          // Add timeout message to the messages
          setMessagesByTopic(prev => ({
            ...prev,
            [topic]: [...(prev[topic] || []), timeoutMessage]
          }));
          
          setConnected(false);
          eventSource.close();
        });
        
        eventSource.onerror = (error) => {
          console.error('EventSource connection error', error);
          if (openReceived) {
            setConnected(false);
          }
          eventSource.close();
        };
      } catch (error) {
        console.error('Error creating EventSource:', error);
        setConnected(false);
      }
    };
    
    // Start the connection
    connectEventSource();
  };

  // Function to clear messages for a topic
  const clearMessages = (topic: string) => {
    setMessagesByTopic(prev => ({
      ...prev,
      [topic]: []
    }));
  };

  // Function to remove a topic
  const removeTopic = (topicToRemove: string) => {
    setTopics(topics.filter(topic => topic !== topicToRemove));
    
    // If removing the current topic, disconnect
    if (currentTopic === topicToRemove) {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
      setCurrentTopic('');
      setConnected(false);
    }
    
    // Also remove messages for this topic
    setMessagesByTopic(prev => {
      const newState = {...prev};
      delete newState[topicToRemove];
      return newState;
    });
  };
  
  // Send a message to the current topic
  const sendMessage = async (content: string) => {
    if (!currentTopic || content.trim() === '') return;
    
    try {
      console.log(`Sending message to ${currentTopic}: ${content}`);
      const response = await fetch(`http://localhost:8080/infocenter/${currentTopic}`, {
        method: 'POST',
        body: content,
        headers: {
          'Content-Type': 'text/plain'
        }
      });
      
      if (!response.ok) {
        console.error('Failed to send message:', response.status);
      } else {
        console.log('Message sent successfully');
      }
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };
  
  // Clean up EventSource on component unmount
  useEffect(() => {
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, []);
  
  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Infocenter Chat</h1>
      </header>
      
      <div className="main-content">
        <div className="sidebar">
          <TopicList 
            topics={topics} 
            currentTopic={currentTopic}
            connected={connected}
            onSelectTopic={joinTopic}
            onRemoveTopic={removeTopic}
          />
        </div>
        
        <div className="chat-area">
          <MessageList 
            messages={getCurrentMessages()} 
            topicName={currentTopic}
            connected={connected}
          />
          
          <MessageInput 
            onSendMessage={sendMessage} 
            disabled={!connected} 
          />
        </div>
      </div>
    </div>
  );
};

export default App;