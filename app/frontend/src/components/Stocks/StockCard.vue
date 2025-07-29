<template>
  <div class="card hover:shadow-lg transition-shadow duration-200 cursor-pointer" @click="$emit('click', stock)">
    <div class="flex justify-between items-start mb-4">
      <div>
        <h3 class="text-lg font-semibold text-gray-900">{{ stock.ticker }}</h3>
        <p class="text-sm text-gray-600 mt-1">{{ stock.company }}</p>
      </div>
      <div class="text-right">
        <span class="text-xs text-gray-500">{{ formatDate(stock.time) }}</span>
      </div>
    </div>

    <div class="space-y-3">
      <!-- Rating Information -->
      <div class="flex justify-between items-center">
        <span class="text-sm text-gray-600">Rating:</span>
        <div class="flex items-center space-x-2">
          <span class="text-sm font-medium text-gray-700">{{ stock.rating_from }}</span>
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
          </svg>
          <span class="text-sm font-medium" :class="getRatingColor(stock.rating_to)">{{ stock.rating_to }}</span>
        </div>
      </div>

      <!-- Target Price Range -->
      <div v-if="stock.target_from && stock.target_to" class="flex justify-between items-center">
        <span class="text-sm text-gray-600">Target:</span>
        <span class="text-sm font-medium text-gray-900">
          {{ stock.target_from }} - {{ stock.target_to }}
        </span>
      </div>

      <!-- Action -->
      <div class="flex justify-between items-center">
        <span class="text-sm text-gray-600">Action:</span>
        <span class="text-sm font-medium" :class="getActionColor(stock.action)">{{ stock.action }}</span>
      </div>

      <!-- Brokerage -->
      <div class="flex justify-between items-center">
        <span class="text-sm text-gray-600">Brokerage:</span>
        <span class="text-sm font-medium text-gray-700">{{ stock.brokerage }}</span>
      </div>
    </div>

    <!-- Action Badge -->
    <div class="mt-4 flex justify-end">
      <span :class="getActionBadgeClass(stock.action)">
        {{ stock.action }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Stock } from '@/types'

interface Props {
  stock: Stock
}

defineProps<Props>()
defineEmits<{
  click: [stock: Stock]
}>()

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

function getRatingColor(rating: string): string {
  if (!rating) return 'text-gray-600'
  const lowerRating = rating.toLowerCase()
  if (lowerRating.includes('buy') || lowerRating.includes('outperform')) {
    return 'text-success-600'
  } else if (lowerRating.includes('sell') || lowerRating.includes('underperform')) {
    return 'text-danger-600'
  } else {
    return 'text-warning-600'
  }
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
  if (!action) return 'badge-gray'
  const lowerAction = action.toLowerCase()
  if (lowerAction.includes('upgrade')) {
    return 'badge-success'
  } else if (lowerAction.includes('downgrade')) {
    return 'badge-danger'
  } else {
    return 'badge-warning'
  }
}
</script>