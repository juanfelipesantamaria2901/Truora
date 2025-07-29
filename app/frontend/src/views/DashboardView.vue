<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Investment Dashboard</h1>
        <p class="text-gray-600 mt-1">Your personalized stock recommendations and market insights</p>
      </div>
      <button
        @click="refreshData"
        :disabled="isLoading"
        class="btn-primary flex items-center space-x-2"
      >
        <svg
          class="w-4 h-4"
          :class="{ 'animate-spin': isLoading }"
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
        <span>{{ isLoading ? 'Refreshing...' : 'Refresh' }}</span>
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-primary-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Total Stocks</p>
            <p class="text-2xl font-semibold text-gray-900">{{ stocksStore.stocks.length }}</p>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-success-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Buy Recommendations</p>
            <p class="text-2xl font-semibold text-gray-900">{{ buyRecommendations }}</p>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-warning-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-warning-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Hold Recommendations</p>
            <p class="text-2xl font-semibold text-gray-900">{{ holdRecommendations }}</p>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-danger-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Sell Recommendations</p>
            <p class="text-2xl font-semibold text-gray-900">{{ sellRecommendations }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Top Recommendations -->
    <div class="card">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-semibold text-gray-900">Top Investment Recommendations</h2>
        <RouterLink to="/recommendations" class="btn-secondary text-sm">
          View All
        </RouterLink>
      </div>

      <div v-if="stocksStore.loading" class="flex justify-center py-8">
        <LoadingSpinner text="Loading recommendations..." />
      </div>

      <div v-else-if="stocksStore.error" class="text-center py-8">
        <div class="text-danger-600 mb-2">
          <svg class="w-12 h-12 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-lg font-medium">Error loading recommendations</p>
          <p class="text-sm text-gray-600 mt-1">{{ stocksStore.error }}</p>
        </div>
        <button @click="refreshData" class="btn-primary mt-4">
          Try Again
        </button>
      </div>

      <div v-else-if="topRecommendations.length === 0" class="text-center py-8">
        <div class="text-gray-400 mb-4">
          <svg class="w-12 h-12 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <p class="text-lg font-medium text-gray-900">No recommendations available</p>
          <p class="text-sm text-gray-600 mt-1">Generate recommendations to see investment opportunities</p>
        </div>
        <button @click="generateRecommendations" class="btn-primary">
          Generate Recommendations
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="recommendation in topRecommendations"
          :key="recommendation.id"
          class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
          @click="viewRecommendation(recommendation)"
        >
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-semibold text-gray-900">{{ recommendation.stock.ticker }}</h3>
              <p class="text-sm text-gray-600">{{ recommendation.stock.company }}</p>
            </div>
            <span class="text-xs px-2 py-1 rounded-full" :class="getScoreBadgeClass(recommendation.recommendation_score || 0)">
              {{ (recommendation.recommendation_score || 0).toFixed(1) }}
            </span>
          </div>
          
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">Target Price:</span>
              <span class="font-medium">{{ recommendation.target_to || recommendation.stock?.target_to || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Previous Target:</span>
              <span class="text-gray-500">{{ recommendation.target_from || recommendation.stock?.target_from || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Risk Level:</span>
              <span :class="getRiskColor(recommendation.risk_level || 'medium')">{{ recommendation.risk_level || 'medium' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Expected Return:</span>
              <span class="text-success-600 font-medium">{{ (recommendation.expected_return || 0).toFixed(1) }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="card">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-semibold text-gray-900">Recent Stock Updates</h2>
        <RouterLink to="/stocks" class="btn-secondary text-sm">
          View All Stocks
        </RouterLink>
      </div>

      <div v-if="recentStocks.length === 0" class="text-center py-8 text-gray-500">
        <p>No recent stock updates available</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="stock in recentStocks"
          :key="stock.id"
          class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <div class="flex items-center space-x-4">
            <div>
              <h4 class="font-medium text-gray-900">{{ stock.ticker }}</h4>
              <p class="text-sm text-gray-600">{{ stock.company }}</p>
              <p class="text-xs text-gray-500">Target: {{ stock.target_to }} (from {{ stock.target_from }})</p>
            </div>
          </div>
          
          <div class="flex items-center space-x-4">
            <div class="text-right">
              <p class="text-sm font-medium" :class="getActionColor(stock.action)">{{ stock.action }}</p>
              <p class="text-xs text-gray-500">{{ formatDate(stock.time) }}</p>
            </div>
            <span :class="getActionBadgeClass(stock.action)">
                {{ stock.rating_to }}
              </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useStocksStore } from '@/stores/stocks'
import LoadingSpinner from '@/components/UI/LoadingSpinner.vue'
import type { StockRecommendation } from '@/types'

const stocksStore = useStocksStore()

const isLoading = computed(() => stocksStore.loading)

const topRecommendations = computed(() => 
  stocksStore.topRecommendations.slice(0, 6)
)

const recentStocks = computed(() => 
  stocksStore.stocks.slice(0, 5)
)

const buyRecommendations = computed(() => 
  stocksStore.recommendations.filter(rec => 
    (rec.recommendation_score || 0) >= 75 && 
    (rec.analyst_sentiment === 'bullish' || rec.risk_level === 'low')
  ).length
)

const holdRecommendations = computed(() => 
  stocksStore.recommendations.filter(rec => 
    (rec.recommendation_score || 0) >= 50 && 
    (rec.recommendation_score || 0) < 75 &&
    (rec.analyst_sentiment === 'neutral' || rec.risk_level === 'medium')
  ).length
)

const sellRecommendations = computed(() => 
  stocksStore.recommendations.filter(rec => 
    (rec.recommendation_score || 0) < 50 || 
    (rec.analyst_sentiment === 'bearish' || rec.risk_level === 'high')
  ).length
)

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric'
  })
}

function getActionColor(action: string): string {
  if (!action) return 'text-gray-700'
  const lowerAction = action.toLowerCase()
  if (lowerAction.includes('upgrade')) {
    return 'text-success-600'
  } else if (lowerAction.includes('downgrade')) {
    return 'text-danger-600'
  } else {
    return 'text-gray-700'
  }
}

function getActionBadgeClass(action: string): string {
  if (!action) return 'badge-warning'
  const lowerAction = action.toLowerCase()
  if (lowerAction.includes('upgrade')) {
    return 'badge-success'
  } else if (lowerAction.includes('downgrade')) {
    return 'badge-danger'
  } else {
    return 'badge-warning'
  }
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

function viewRecommendation(recommendation: StockRecommendation) {
  // Navigate to recommendation detail or show modal
  console.log('View recommendation:', recommendation)
}

async function refreshData() {
  await Promise.all([
    stocksStore.fetchStocks(),
    stocksStore.fetchRecommendations()
  ])
}

async function generateRecommendations() {
  await stocksStore.generateRecommendations()
}

onMounted(() => {
  refreshData()
})
</script>