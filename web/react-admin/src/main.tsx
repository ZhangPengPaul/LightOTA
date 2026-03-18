import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import './index.css'
import App from './App.tsx'
import Login from './pages/Login.tsx'
import Tenants from './pages/Tenants.tsx'
import Products from './pages/Products.tsx'
import Firmwares from './pages/Firmwares.tsx'
import UpgradeTasks from './pages/UpgradeTasks.tsx'
import { setApiKey } from './api/client'

const savedApiKey = localStorage.getItem('api_key')
if (savedApiKey) {
  setApiKey(savedApiKey)
}

const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/',
    element: <App />,
    children: [
      {
        path: 'tenants',
        element: <Tenants />,
      },
      {
        path: 'products',
        element: <Products />,
      },
      {
        path: 'firmwares',
        element: <Firmwares />,
      },
      {
        path: 'tasks',
        element: <UpgradeTasks />,
      },
      {
        index: true,
        element: <div style={{ padding: '20px' }}>
          <h2>Welcome to LightOTA Admin</h2>
          <p>Select a menu item from the sidebar to get started.</p>
        </div>,
      },
    ],
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
