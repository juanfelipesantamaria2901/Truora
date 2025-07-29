export interface Stock {
  id: number
  ticker: string
  company: string
  target_from?: string
  target_to?: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
  last_updated: string
}

export interface StockRecommendation {
  id: number
  stock_id: number
  stock: Stock
  recommendation_score: number
  risk_level: string
  expected_return: number
  time_horizon: string
  reasons: string
  analyst_sentiment: string
  upgrade_count: number
  downgrade_count: number
  created_at: string
  updated_at: string
  target_to?: string
  target_from?: string
}

export interface PaginationInfo {
  limit: number
  offset: number
  total: number
}

export interface ApiResponse<T> {
  data: T
  pagination?: PaginationInfo
  count?: number
}

export interface SearchFilters {
  search?: string
  query?: string
  action?: string
  rating?: string
  brokerage?: string
  risk_level?: string
  minScore?: number
  maxScore?: number
  sortBy?: 'time' | 'ticker' | 'company' | 'action' | 'rating_to' | 'score' | 'risk' | 'return'
  sortOrder?: 'asc' | 'desc'
}