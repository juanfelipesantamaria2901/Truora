<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Investment Recommendations</h1>
        <p class="text-gray-600 mt-1">AI-powered stock recommendations based on market analysis</p>
      </div>
      <div class="flex space-x-3">
        <button
          @click="fetchRecommendations"
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
          @click="generateRecommendations"
          :disabled="stocksStore.loading"
          class="btn-primary flex items-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          <span>Generate New</span>
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Filter Recommendations</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Risk Level Filter -->
        <div>
          <label for="riskLevel" class="block text-sm font-medium text-gray-700 mb-1">
            Risk Level
          </label>
          <select
            id="risk_level"
            v-model="filters.risk_level"
            class="input-field"
            @change="applyFilters"
          >
            <option value="">All Risk Levels</option>
            <option value="Low">Low Risk</option>
            <option value="Medium">Medium Risk</option>
            <option value="High">High Risk</option>
          </select>
        </div>

        <!-- Min Score Filter -->
        <div>
          <label for="minScore" class="block text-sm font-medium text-gray-700 mb-1">
            Min Score
          </label>
          <input
            id="minScore"
            v-model.number="filters.minScore"
            type="number"
            min="0"
            max="10"
            step="0.1"
            placeholder="0.0"
            class="input-field"
            @input="applyFilters"
          />
        </div>

        <!-- Max Score Filter -->
        <div>
          <label for="maxScore" class="block text-sm font-medium text-gray-700 mb-1">
            Max Score
          </label>
          <input
            id="maxScore"
            v-model.number="filters.maxScore"
            type="number"
            min="0"
            max="10"
            step="0.1"
            placeholder="10.0"
            class="input-field"
            @input="applyFilters"
          />
        </div>

        <!-- Sort Options -->
        <div>
          <label for="sortBy" class="block text-sm font-medium text-gray-700 mb-1">
            Sort By
          </label>
          <select
            id="sortBy"
            v-model="filters.sortBy"
            class="input-field"
            @change="applyFilters"
          >
            <option value="score">Score</option>
            <option value="risk">Risk Level</option>
            <option value="return">Expected Return</option>
            <option value="company">Company</option>
          </select>
        </div>
      </div>

      <div class="mt-4 flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <label for="sortOrder" class="text-sm font-medium text-gray-700">
              Order:
            </label>
            <select
              id="sortOrder"
              v-model="filters.sortOrder"
              class="input-field min-w-0 w-auto"
              @change="applyFilters"
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
        </div>
        
        <button
          @click="clearFilters"
          class="btn-secondary text-sm"
        >
          Clear Filters
        </button>
      </div>
    </div>

    <!-- Results Summary -->
    <div class="flex justify-between items-center">
      <p class="text-sm text-gray-600">
        Showing {{ filteredRecommendations.length }} of {{ stocksStore.recommendations.length }} recommendations
      </p>
      <div class="flex items-center space-x-2">
        <label for="pageSize" class="text-sm text-gray-600">Show:</label>
        <select
          id="pageSize"
          v-model="pageSize"
          class="input-field w-auto min-w-0"
        >
          <option value="12">12</option>
          <option value="24">24</option>
          <option value="48">48</option>
        </select>
        <span class="text-sm text-gray-600">per page</span>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="stocksStore.loading" class="flex justify-center py-12">
      <LoadingSpinner size="lg" text="Loading recommendations..." />
    </div>

    <!-- Error State -->
    <div v-else-if="stocksStore.error" class="text-center py-12">
      <div class="text-danger-600 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">Error loading recommendations</h3>
        <p class="text-sm text-gray-600 mt-1">{{ stocksStore.error }}</p>
      </div>
      <button @click="fetchRecommendations" class="btn-primary">
        Try Again
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="stocksStore.recommendations.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">No recommendations available</h3>
        <p class="text-sm text-gray-600 mt-1">Generate recommendations to see investment opportunities</p>
      </div>
      <button @click="generateRecommendations" class="btn-primary">
        Generate Recommendations
      </button>
    </div>

    <!-- No Results -->
    <div v-else-if="filteredRecommendations.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">No recommendations found</h3>
        <p class="text-sm text-gray-600 mt-1">Try adjusting your filter criteria</p>
      </div>
      <button @click="clearFilters" class="btn-secondary">
        Clear Filters
      </button>
    </div>

    <!-- Recommendations Grid -->
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="recommendation in paginatedRecommendations"
          :key="recommendation.id"
          class="card hover:shadow-lg transition-shadow duration-200 cursor-pointer"
          @click="viewRecommendationDetail(recommendation)"
        >
          <!-- Header -->
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ recommendation.stock.ticker }}</h3>
              <p class="text-sm text-gray-600 mt-1">{{ recommendation.stock.company }}</p>
            </div>
            <div class="text-right">
              <span class="text-xs px-2 py-1 rounded-full" :class="getScoreBadgeClass(recommendation.recommendation_score)">
          {{ Number(recommendation.recommendation_score || 0).toFixed(1) }}
        </span>
            </div>
          </div>

          <!-- Metrics -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Target Price:</span>
            <span class="text-sm font-medium">{{ recommendation.target_to || recommendation.stock?.target_to || 'N/A' }}</span>
          </div>
          
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Previous Target:</span>
            <span class="text-sm font-medium text-gray-500">{{ recommendation.target_from || recommendation.stock?.target_from || 'N/A' }}</span>
          </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Risk Level:</span>
              <span class="text-sm font-medium" :class="getRiskColor(recommendation.risk_level)">{{ recommendation.risk_level }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Expected Return:</span>
              <span class="text-sm font-medium text-success-600">{{ Number(recommendation.expected_return || 0).toFixed(1) }}%</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Time Horizon:</span>
              <span class="text-sm font-medium text-gray-900">{{ recommendation.time_horizon || 'medium' }}</span>
            </div>
          </div>

          <!-- Reasoning Preview -->
          <div class="mt-4 pt-4 border-t border-gray-200">
            <p class="text-xs text-gray-600 line-clamp-2">{{ recommendation.reasons || 'No specific reasons provided' }}</p>
          </div>

          <!-- Action Button -->
          <div class="mt-4 flex justify-end">
            <span class="text-xs px-3 py-1 rounded-full" :class="getRecommendationBadgeClass('Buy')">
              Buy
            </span>
          </div>
        </div>
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
import LoadingSpinner from '@/components/UI/LoadingSpinner.vue'
import type { StockRecommendation, SearchFilters } from '@/types'

const stocksStore = useStocksStore()

const filters = ref<SearchFilters>({
  risk_level: '',
  minScore: undefined,
  maxScore: undefined,
  sortBy: 'score',
  sortOrder: 'desc'
})

const pageSize = ref(12)
const currentPage = ref(1)

const filteredRecommendations = computed(() => {
  let result = [...stocksStore.recommendations]

  // Apply risk level filter
  if (filters.value.risk_level) {
    result = result.filter(rec => rec.risk_level === filters.value.risk_level)
  }

  // Apply score filters
  if (filters.value.minScore !== undefined) {
    result = result.filter(rec => (rec.recommendation_score || 0) >= filters.value.minScore!)
  }
  if (filters.value.maxScore !== undefined) {
    result = result.filter(rec => (rec.recommendation_score || 0) <= filters.value.maxScore!)
  }

  // Apply sorting
  if (filters.value.sortBy) {
    result.sort((a, b) => {
      let aValue: any
      let bValue: any

      switch (filters.value.sortBy) {
        case 'score':
          aValue = a.recommendation_score || 0
          bValue = b.recommendation_score || 0
          break
        case 'risk':
          const riskOrder = { 'Low': 1, 'Medium': 2, 'High': 3 }
          aValue = riskOrder[(a.risk_level || 'medium') as keyof typeof riskOrder] || 0
          bValue = riskOrder[(b.risk_level || 'medium') as keyof typeof riskOrder] || 0
          break
        case 'return':
          aValue = a.expected_return || 0
          bValue = b.expected_return || 0
          break
        case 'company':
          aValue = (a.stock?.company || '').toLowerCase()
          bValue = (b.stock?.company || '').toLowerCase()
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
  Math.ceil(filteredRecommendations.value.length / pageSize.value)
)

const paginatedRecommendations = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRecommendations.value.slice(start, end)
})

function applyFilters() {
  currentPage.value = 1 // Reset to first page when filters change
}

function clearFilters() {
  filters.value = {
    risk_level: '',
    minScore: undefined,
    maxScore: undefined,
    sortBy: 'score',
    sortOrder: 'desc'
  }
  currentPage.value = 1
}

function getScoreBadgeClass(score: number): string {
  const numScore = Number(score) || 0
  if (numScore >= 80) return 'bg-success-100 text-success-800'
  if (numScore >= 60) return 'bg-warning-100 text-warning-800'
  return 'bg-danger-100 text-danger-800'
}

function getRiskColor(riskLevel: string): string {
  if (!riskLevel) return 'text-warning-600'
  const lower = riskLevel.toLowerCase()
  if (lower === 'low') return 'text-success-600'
  if (lower === 'medium') return 'text-warning-600'
  return 'text-danger-600'
}

function getRecommendationBadgeClass(recommendation: string): string {
  if (!recommendation) return 'bg-warning-100 text-warning-800'
  const lower = recommendation.toLowerCase()
  if (lower.includes('buy') || lower.includes('strong buy')) {
    return 'bg-success-100 text-success-800'
  } else if (lower.includes('sell')) {
    return 'bg-danger-100 text-danger-800'
  } else {
    return 'bg-warning-100 text-warning-800'
  }
}

function viewRecommendationDetail(recommendation: StockRecommendation) {
  // Navigate to recommendation detail view or show modal
  console.log('View recommendation detail:', recommendation)
}

async function fetchRecommendations() {
  await stocksStore.fetchRecommendations()
}

async function generateRecommendations() {
  await stocksStore.generateRecommendations()
}

// Reset page when page size changes
watch(pageSize, () => {
  currentPage.value = 1
})

// Reset page when filtered results change
watch(filteredRecommendations, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = Math.max(1, totalPages.value)
  }
})

onMounted(() => {
  if (stocksStore.recommendations.length === 0) {
    fetchRecommendations()
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>