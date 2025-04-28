import React, { useEffect, useRef } from 'react';

interface Message {
  id: number;
  content: string;
  topic: string;
  event?: string;
}

interface MessageListProps {
  messages: Message[];
  topicName: string;
  connected: boolean;
}

const MessageList: React.FC<MessageListProps> = ({ 
  messages, 
  topicName, 
  connected, 
}) => {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);
  
  return (
    <div className="message-container">
      <div className="message-header">
        <h2>
          {topicName ? (
            <>
              Topic: {topicName} 
              <span className={`connection-status ${connected ? 'connected' : 'disconnected'}`}>
                ({connected ? 'Connected' : 'timeout'})
              </span>
            </>
          ) : (
            'Select a topic to start chatting'
          )}
        </h2>
      </div>
      
      <div className="message-list">
        {messages.length === 0 ? (
          <div className="empty-state">No messages yet</div>
        ) : (
          messages.map(message => (
            <div key={message.id} className="message">
              <pre className="message-format">
id: {message.id}
<br></br>
event: {message.event || "msg"}
<br></br>
message: {message.content}
              </pre>
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>
    </div>
  );
};

export default MessageList;