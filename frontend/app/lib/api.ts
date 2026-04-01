const API_BASE = import.meta.env.VITE_API_URL || '/api/v1'

let authToken: string | null = null

export function setToken(token: string) {
  authToken = token
  if (typeof window !== 'undefined') {
    localStorage.setItem('smc_token', token)
  }
}

export function getToken(): string | null {
  if (authToken) return authToken
  if (typeof window !== 'undefined') {
    authToken = localStorage.getItem('smc_token')
  }
  return authToken
}

export function clearToken() {
  authToken = null
  if (typeof window !== 'undefined') {
    localStorage.removeItem('smc_token')
    localStorage.removeItem('smc_user')
  }
}

export function getUser(): any {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('smc_user')
    return stored ? JSON.parse(stored) : null
  }
  return null
}

async function fetchAPI<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

  if (res.status === 401) {
    clearToken()
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
    throw new Error('Unauthorized')
  }

  const data = await res.json()
  if (!data.success) {
    throw new Error(data.error || 'API Error')
  }
  return data as T
}

export const api = {
  // Auth
  login: (username: string, password: string) =>
    fetchAPI<any>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  // Dashboard
  getStats: () => fetchAPI<any>('/dashboard/stats'),
  getConsumptionTrend: (days = 30) => fetchAPI<any>(`/dashboard/consumption?days=${days}`),
  getAlerts: (page = 1) => fetchAPI<any>(`/dashboard/alerts?page=${page}`),

  // Meters
  getMeters: (page = 1, perPage = 20) => fetchAPI<any>(`/meters?page=${page}&per_page=${perPage}`),
  getMeter: (id: string) => fetchAPI<any>(`/meters/${id}`),
  getReadings: (meterId: string, page = 1) => fetchAPI<any>(`/meters/${meterId}/readings?page=${page}`),
  registerMeter: (data: { meter_id: string; consumer_id: string; location: string; meter_type: string }) =>
    fetchAPI<any>('/meters', { method: 'POST', body: JSON.stringify(data) }),

  // Billing
  getBills: (consumerId?: string, page = 1) => {
    const params = consumerId ? `?consumer_id=${consumerId}&page=${page}` : `?page=${page}`
    return fetchAPI<any>(`/bills${params}`)
  },
  generateBill: (data: { meter_id: string; period_start: string; period_end: string }) =>
    fetchAPI<any>('/bills/generate', { method: 'POST', body: JSON.stringify(data) }),
  payBill: (billId: string) =>
    fetchAPI<any>(`/bills/${billId}/pay`, { method: 'POST' }),
  verifyBill: (billId: string) => fetchAPI<any>(`/verify/bill/${billId}`),

  // Disputes
  getDisputes: (page = 1) => fetchAPI<any>(`/disputes?page=${page}`),
  fileDispute: (billId: string, reason: string) =>
    fetchAPI<any>('/disputes', {
      method: 'POST',
      body: JSON.stringify({ bill_id: billId, reason }),
    }),
  resolveDispute: (disputeId: string, resolution: string) =>
    fetchAPI<any>(`/disputes/${disputeId}/resolve`, {
      method: 'POST',
      body: JSON.stringify({ resolution }),
    }),

  // Consumers
  getConsumers: (page = 1) => fetchAPI<any>(`/consumers?page=${page}`),

  // Tariffs
  getTariffs: (category?: string) => {
    const params = category ? `?category=${category}` : ''
    return fetchAPI<any>(`/tariffs${params}`)
  },

  // Alerts
  acknowledgeAlert: (alertId: number) =>
    fetchAPI<any>(`/alerts/${alertId}/acknowledge`, { method: 'POST' }),

  // Profile
  getProfile: () => fetchAPI<any>('/profile'),
}
