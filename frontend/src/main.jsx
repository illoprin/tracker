import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@/style/output.css'
import '@/style/App.css'
import App from '@/App.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
