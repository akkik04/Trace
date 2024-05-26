import React, { useState } from 'react';
import './App.css';

function App() {
  const [logMessage, setLogMessage] = useState('');

  const generateRandomLog = () => {
    const ip = `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`;
    const userIdentifier = '-';
    const userId = 'user-' + Math.floor(Math.random() * 1000);
    const time = new Date().toISOString();
    const method = ['GET', 'POST', 'PUT', 'DELETE'][Math.floor(Math.random() * 4)];
    const url = ['/home', '/login', '/signup', '/logout'][Math.floor(Math.random() * 4)];
    const protocol = 'HTTP/1.1';
    const status = [200, 404, 500, 403][Math.floor(Math.random() * 4)];
    const size = Math.floor(Math.random() * 10000);

    const logMessage = `${ip} ${userIdentifier} ${userId} [${time}] "${method} ${url} ${protocol}" ${status} ${size}`;

    return logMessage;
  };

  const sendLog = async () => {
    const logMessage = generateRandomLog();
    setLogMessage(logMessage);

    try {
      const response = await fetch('http://localhost:3000/logs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message: logMessage, timestamp: new Date().toISOString() }),
      });

      if (response.ok) {
        console.log('Log sent successfully');
      } else {
        console.error('Failed to send log');
      }
    } catch (error) {
      console.error('Error sending log:', error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Log Generator</h1>
        <button className="log-button" onClick={sendLog}>Generate and Send Log</button>
        {logMessage && (
          <div className="log-message">
            <h2>Generated Random Log:</h2>
            <p>{logMessage}</p>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
