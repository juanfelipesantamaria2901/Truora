import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Stock, StockRecommendation, SearchFilters } from '@/types'
import { apiService } from '@/services/api'

export const useStocksStore = defineStore('stocks', () => {
  // State
  const stocks = ref<Stock[]>([])
  const recommendations = ref<StockRecommendation[]>([])
  const selectedStock = ref<Stock | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')
  const filters = ref<SearchFilters>({})
  const pagination = ref({
    limit: 20,
    offset: 0,
    total: 0
  })

  // Getters
  const filteredRecommendations = computed(() => {
    let filtered = [...recommendations.value]

    if (filters.value.risk_level) {
        filtered = filtered.filter(rec => rec.risk_level === filters.value.risk_level)
      }

    if (filters.value.minScore !== undefined) {
      filtered = filtered.filter(rec => (rec.recommendation_score || 0) >= filters.value.minScore!)
    }

    if (filters.value.maxScore !== undefined) {
      filtered = filtered.filter(rec => (rec.recommendation_score || 0) <= filters.value.maxScore!)
    }

    if (filters.value.query) {
        const query = filters.value.query.toLowerCase()
        filtered = filtered.filter(rec => 
          (rec.stock?.company || '').toLowerCase().includes(query) ||
          (rec.stock?.ticker || '').toLowerCase().includes(query)
        )
      }

    // Sort
    if (filters.value.sortBy) {
      filtered.sort((a, b) => {
        let aValue: any, bValue: any
        
        switch (filters.value.sortBy) {
          case 'score':
            aValue = a.recommendation_score || 0
            bValue = b.recommendation_score || 0
            break
          case 'return':
            aValue = a.expected_return || 0
            bValue = (b.stock?.company || '').toLowerCase()
            break
          case 'risk':
            aValue = (a.risk_level || '').toLowerCase()
            bValue = (b.risk_level || '').toLowerCase()
            break
          default:
            return 0
        }

        if (typeof aValue === 'string') {
          aValue = aValue.toLowerCase()
          bValue = bValue.toLowerCase()
        }

        const result = aValue < bValue ? -1 : aValue > bValue ? 1 : 0
        return filters.value.sortOrder === 'desc' ? -result : result
      })
    }

    return filtered
  })

  const topRecommendations = computed(() => {
    return recommendations.value
      .slice()
      .sort((a, b) => (b.recommendation_score || 0) - (a.recommendation_score || 0))
      .slice(0, 5)
  })

  // Actions
  async function fetchStocks(params?: { limit?: number; offset?: number; q?: string }) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiService.getStocks(params)
      stocks.value = response.data
      if (response.pagination) {
        pagination.value = response.pagination
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stocks'
    } finally {
      loading.value = false
    }
  }

  async function fetchStockByTicker(ticker: string) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiService.getStockByTicker(ticker)
      selectedStock.value = response.data
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stock'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchRecommendations(limit?: number) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiService.getRecommendations(limit)
      recommendations.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recommendations'
    } finally {
      loading.value = false
    }
  }

  async function generateRecommendations() {
    loading.value = true
    error.value = null
    
    try {
      await apiService.generateRecommendations()
      // Refresh recommendations after generation
      await fetchRecommendations()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to generate recommendations'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchAndStoreStocks() {
    loading.value = true
    error.value = null
    
    try {
      await apiService.fetchStocks()
      // Refresh stocks after fetching
      await fetchStocks()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch and store stocks'
      throw err
    } finally {
      loading.value = false
    }
  }

  function updateFilters(newFilters: Partial<SearchFilters>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearFilters() {
    filters.value = {}
    searchQuery.value = ''
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
    filters.value.query = query
  }

  return {
    // State
    stocks,
    recommendations,
    selectedStock,
    loading,
    error,
    searchQuery,
    filters,
    pagination,
    // Getters
    filteredRecommendations,
    topRecommendations,
    // Actions
    fetchStocks,
    fetchStockByTicker,
    fetchRecommendations,
    generateRecommendations,
    fetchAndStoreStocks,
    updateFilters,
    clearFilters,
    setSearchQuery
  }
})