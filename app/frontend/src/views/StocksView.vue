<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Stocks</h1>
        <p class="text-gray-600 mt-1">Browse and analyze stock recommendations</p>
      </div>
      <div class="flex space-x-3">
        <button
          @click="fetchStocks"
          :disabled="stocksStore.loading"
          class="btn-secondary flex items-center space-x-2"
        >
          <svg
            class="w-4 h-4"
            :class="{ 'animate-spin': stocksStore.loading }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
          <span>{{ stocksStore.loading ? 'Loading...' : 'Refresh' }}</span>
        </button>
        <button
          @click="fetchAndStoreStocks"
          :disabled="stocksStore.loading"
          class="btn-primary flex items-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span>Fetch New Data</span>
        </button>
      </div>
    </div>

    <!-- Filters -->
    <StockFilters
      :filters="filters"
      :available-brokerages="availableBrokerages"
      @update:filters="updateFilters"
    />

    <!-- Results Summary -->
    <div class="flex justify-between items-center">
      <p class="text-sm text-gray-600">
        Showing {{ filteredStocks.length }} of {{ stocksStore.stocks.length }} stocks
      </p>
      <div class="flex items-center space-x-2">
        <label for="pageSize" class="text-sm text-gray-600">Show:</label>
        <select
          id="pageSize"
          v-model="pageSize"
          class="input-field w-auto min-w-0"
        >
          <option value="10">10</option>
          <option value="25">25</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
        <span class="text-sm text-gray-600">per page</span>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="stocksStore.loading" class="flex justify-center py-12">
      <LoadingSpinner size="lg" text="Loading stocks..." />
    </div>

    <!-- Error State -->
    <div v-else-if="stocksStore.error" class="text-center py-12">
      <div class="text-danger-600 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">Error loading stocks</h3>
        <p class="text-sm text-gray-600 mt-1">{{ stocksStore.error }}</p>
      </div>
      <button @click="fetchStocks" class="btn-primary">
        Try Again
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="stocksStore.stocks.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">No stocks available</h3>
        <p class="text-sm text-gray-600 mt-1">Fetch stock data to get started</p>
      </div>
      <button @click="fetchAndStoreStocks" class="btn-primary">
        Fetch Stock Data
      </button>
    </div>

    <!-- No Results -->
    <div v-else-if="filteredStocks.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">No stocks found</h3>
        <p class="text-sm text-gray-600 mt-1">Try adjusting your search criteria</p>
      </div>
      <button @click="clearFilters" class="btn-secondary">
        Clear Filters
      </button>
    </div>

    <!-- Stocks Grid -->
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <StockCard
          v-for="stock in paginatedStocks"
          :key="stock.id"
          :stock="stock"
          @click="viewStockDetail"
        />
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="mt-8 flex justify-center">
        <nav class="flex items-center space-x-2">
          <button
            @click="currentPage = 1"
            :disabled="currentPage === 1"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            First
          </button>
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          
          <span class="px-4 py-2 text-sm text-gray-700">
            Page {{ currentPage }} of {{ totalPages }}
          </span>
          
          <button
            @click="currentPage++"
            :disabled="currentPage === totalPages"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
          <button
            @click="currentPage = totalPages"
            :disabled="currentPage === totalPages"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Last
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useStocksStore } from '@/stores/stocks'
import StockCard from '@/components/Stocks/StockCard.vue'
import StockFilters from '@/components/Stocks/StockFilters.vue'
import LoadingSpinner from '@/components/UI/LoadingSpinner.vue'
import type { Stock, SearchFilters } from '@/types'

const stocksStore = useStocksStore()

const filters = ref<SearchFilters>({
  search: '',
  action: '',
  rating: '',
  brokerage: '',
  sortBy: 'time',
  sortOrder: 'desc'
})

const pageSize = ref(25)
const currentPage = ref(1)

const availableBrokerages = computed(() => {
  const brokerages = new Set(stocksStore.stocks.map(stock => stock.brokerage))
  return Array.from(brokerages).sort()
})

const filteredStocks = computed(() => {
  let result = [...stocksStore.stocks]

  // Apply search filter
  if (filters.value.search) {
    const searchTerm = filters.value.search.toLowerCase()
    result = result.filter(stock => 
      stock.ticker?.toLowerCase().includes(searchTerm) ||
      stock.company?.toLowerCase().includes(searchTerm)
    )
  }

  // Apply action filter
  if (filters.value.action) {
    const actionFilter = filters.value.action.toLowerCase()
    result = result.filter(stock => 
      stock.action?.toLowerCase().includes(actionFilter)
    )
  }

  // Apply rating filter
  if (filters.value.rating) {
    const ratingFilter = filters.value.rating.toLowerCase()
    result = result.filter(stock => 
      stock.rating_to?.toLowerCase().includes(ratingFilter)
    )
  }

  // Apply brokerage filter
  if (filters.value.brokerage) {
    result = result.filter(stock => stock.brokerage === filters.value.brokerage)
  }

  // Apply sorting
  if (filters.value.sortBy) {
    result.sort((a, b) => {
      let aValue: any
      let bValue: any

      switch (filters.value.sortBy) {
        case 'time':
          aValue = new Date(a.time)
          bValue = new Date(b.time)
          break
        case 'ticker':
          aValue = a.ticker
          bValue = b.ticker
          break
        case 'company':
          aValue = a.company
          bValue = b.company
          break
        case 'action':
          aValue = a.action
          bValue = b.action
          break
        case 'rating_to':
          aValue = a.rating_to
          bValue = b.rating_to
          break
        default:
          return 0
      }

      if (aValue < bValue) {
        return filters.value.sortOrder === 'asc' ? -1 : 1
      }
      if (aValue > bValue) {
        return filters.value.sortOrder === 'asc' ? 1 : -1
      }
      return 0
    })
  }

  return result
})

const totalPages = computed(() => 
  Math.ceil(filteredStocks.value.length / pageSize.value)
)

const paginatedStocks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredStocks.value.slice(start, end)
})

function updateFilters(newFilters: SearchFilters) {
  filters.value = { ...newFilters }
  currentPage.value = 1 // Reset to first page when filters change
}

function clearFilters() {
  filters.value = {
    search: '',
    action: '',
    rating: '',
    brokerage: '',
    sortBy: 'time',
    sortOrder: 'desc'
  }
  currentPage.value = 1
}

function viewStockDetail(stock: Stock) {
  // Navigate to stock detail view or show modal
  console.log('View stock detail:', stock)
}

async function fetchStocks() {
  await stocksStore.fetchStocks()
}

async function fetchAndStoreStocks() {
  await stocksStore.fetchAndStoreStocks()
}

// Reset page when page size changes
watch(pageSize, () => {
  currentPage.value = 1
})

// Reset page when filtered results change
watch(filteredStocks, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = Math.max(1, totalPages.value)
  }
})

onMounted(() => {
  if (stocksStore.stocks.length === 0) {
    fetchStocks()
  }
})
</script>