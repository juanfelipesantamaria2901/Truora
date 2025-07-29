import type { Stock, StockRecommendation, ApiResponse } from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1'

class ApiService {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    })

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  // Stock endpoints
  async getStocks(params?: {
    limit?: number
    offset?: number
    q?: string
  }): Promise<ApiResponse<Stock[]>> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    if (params?.offset) searchParams.append('offset', params.offset.toString())
    if (params?.q) searchParams.append('q', params.q)

    const query = searchParams.toString()
    return this.request(`/stocks${query ? `?${query}` : ''}`)
  }

  async getStockByTicker(ticker: string): Promise<ApiResponse<Stock>> {
    return this.request(`/stocks/${ticker}`)
  }

  async fetchStocks(): Promise<{ message: string }> {
    return this.request('/stocks/fetch', {
      method: 'POST',
    })
  }

  // Recommendation endpoints
  async getRecommendations(limit?: number): Promise<ApiResponse<StockRecommendation[]>> {
    const query = limit ? `?limit=${limit}` : ''
    return this.request(`/recommendations${query}`)
  }

  async generateRecommendations(): Promise<{ message: string }> {
    return this.request('/recommendations/generate', {
      method: 'POST',
    })
  }

  // Health check
  async healthCheck(): Promise<{ status: string; service: string }> {
    return this.request('/health')
  }
}

export const apiService = new ApiService()
export default apiService
