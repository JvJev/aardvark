// src/components/TopicList.tsx
import React, { useState } from 'react';

interface TopicListProps {
  topics: string[];
  currentTopic: string;
  connected: boolean;
  onSelectTopic: (topic: string) => void;
  onRemoveTopic: (topic: string) => void;
}

const TopicList: React.FC<TopicListProps> = ({ 
  topics, 
  currentTopic, 
  connected,
  onSelectTopic,
  onRemoveTopic 
}) => {
  const [newTopic, setNewTopic] = useState('');
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newTopic.trim()) {
      onSelectTopic(newTopic);
      setNewTopic('');
    }
  };
  
  return (
    <div className="topic-list">
      <h2>Topic (URL endpoint)</h2>
      
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={newTopic}
          onChange={(e) => setNewTopic(e.target.value)}
          placeholder="Enter endpoint name"
        />
        <button type="submit">Join</button>
      </form>
      
      <ul>
        {topics.map(topic => (
          <li key={topic} className="topic-item">
            <div className="topic-name">
              {topic}
              {currentTopic === topic && (
                <span className="connection-indicator">
                  {connected ? " (Connected)" : " (timeout)"}
                </span>
              )}
            </div>
            <div className="topic-actions">
              {(currentTopic !== topic || !connected) && (
                <button 
                  className="connect-btn"
                  onClick={() => onSelectTopic(topic)}
                >
                  Listen
                </button>
              )}
              <button 
                className="remove-topic-btn"
                onClick={() => onRemoveTopic(topic)}
              >
                ×
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TopicList;