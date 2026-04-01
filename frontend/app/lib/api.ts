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

  // Billing
  getBills: (consumerId?: string, page = 1) => {
    const params = consumerId ? `?consumer_id=${consumerId}&page=${page}` : `?page=${page}`
    return fetchAPI<any>(`/bills${params}`)
  },
  verifyBill: (billId: string) => fetchAPI<any>(`/verify/bill/${billId}`),

  // Disputes
  getDisputes: (page = 1) => fetchAPI<any>(`/disputes?page=${page}`),
  fileDispute: (billId: string, reason: string) =>
    fetchAPI<any>('/disputes', {
      method: 'POST',
      body: JSON.stringify({ bill_id: billId, reason }),
    }),

  // Tariffs
  getTariffs: (category?: string) => {
    const params = category ? `?category=${category}` : ''
    return fetchAPI<any>(`/tariffs${params}`)
  },

  // Profile
  getProfile: () => fetchAPI<any>('/profile'),
}
