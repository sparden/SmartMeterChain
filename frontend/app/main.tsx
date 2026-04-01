import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { getToken } from './lib/api'
import './styles.css'

import RootLayout from './routes/__root'
import Dashboard from './routes/index'
import LoginPage from './routes/login'
import MetersPage from './routes/meters'
import BillingPage from './routes/billing'
import DisputesPage from './routes/disputes'
import TariffsPage from './routes/tariffs'
import AlertsPage from './routes/alerts'
import ConsumersPage from './routes/consumers'
import ExplorerPage from './routes/explorer'

const queryClient = new QueryClient()

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = getToken()
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<ProtectedRoute><RootLayout /></ProtectedRoute>}>
            <Route index element={<Dashboard />} />
            <Route path="meters" element={<MetersPage />} />
            <Route path="billing" element={<BillingPage />} />
            <Route path="disputes" element={<DisputesPage />} />
            <Route path="tariffs" element={<TariffsPage />} />
            <Route path="alerts" element={<AlertsPage />} />
            <Route path="consumers" element={<ConsumersPage />} />
            <Route path="explorer" element={<ExplorerPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  </React.StrictMode>
)
