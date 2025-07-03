import App from '@/App';
import React from 'react';
import ReactDOM from 'react-dom/client'; // For React 18+

const rootElement = document.getElementById('root');

if (rootElement) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );
}
