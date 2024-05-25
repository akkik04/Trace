import React from 'react';

function App() {
  const handleClick = (message) => {
    fetch('http://localhost:5000/log', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message, timestamp: new Date().toISOString() })
    });
  };

  return (
    <div>
      <button onClick={() => handleClick('Button 1 clicked')}>Button 1</button>
      <button onClick={() => handleClick('Button 2 clicked')}>Button 2</button>
    </div>
  );
}

export default App;

