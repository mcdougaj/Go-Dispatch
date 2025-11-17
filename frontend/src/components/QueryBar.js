import React, { useState } from 'react';
import './QueryBar.css';

function QueryBar({ onQuery, queryResponse }) {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);

  const exampleQueries = [
    "Where are all my vehicles?",
    "Which drivers are available?",
    "Show me empty trailers",
    "Find vehicle TO6076",
    "What's the status of driver 8024360?",
  ];

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    await onQuery(query);
    setLoading(false);
  };

  const handleExampleClick = (example) => {
    setQuery(example);
  };

  return (
    <div className="query-bar-container">
      <form className="query-bar" onSubmit={handleSubmit}>
        <input
          type="text"
          className="query-input"
          placeholder="Ask me anything about your fleet... (e.g., 'Where are all my vehicles?')"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          disabled={loading}
        />
        <button
          type="submit"
          className="query-button"
          disabled={loading || !query.trim()}
        >
          {loading ? 'Thinking...' : 'Ask AI'}
        </button>
      </form>

      <div className="example-queries">
        {exampleQueries.map((example, index) => (
          <button
            key={index}
            className="example-chip"
            onClick={() => handleExampleClick(example)}
            disabled={loading}
          >
            {example}
          </button>
        ))}
      </div>

      {queryResponse && (
        <div className="query-response">
          <div className="response-header">
            <strong>AI Response:</strong>
            <span className="response-meta">
              {queryResponse.response_time_ms}ms • {queryResponse.model}
            </span>
          </div>
          <div className="response-content">
            {queryResponse.response}
          </div>
          {queryResponse.input_tokens && (
            <div className="response-footer">
              Tokens: {queryResponse.input_tokens} in / {queryResponse.output_tokens} out
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default QueryBar;
