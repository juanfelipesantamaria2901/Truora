<template>
  <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
    <h3 class="text-lg font-medium text-gray-900 mb-4">Search & Filter</h3>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Search Input -->
      <div>
        <label for="search" class="block text-sm font-medium text-gray-700 mb-1">
          Search
        </label>
        <input
          id="search"
          v-model="localFilters.search"
          type="text"
          placeholder="Search by ticker or company..."
          class="input-field"
          @input="debouncedUpdate"
        />
      </div>

      <!-- Action Filter -->
      <div>
        <label for="action" class="block text-sm font-medium text-gray-700 mb-1">
          Action
        </label>
        <select
          id="action"
          v-model="localFilters.action"
          class="input-field"
          @change="updateFilters"
        >
          <option value="">All Actions</option>
          <option value="Upgrade">Upgrade</option>
          <option value="Downgrade">Downgrade</option>
          <option value="Initiate">Initiate</option>
          <option value="Maintain">Maintain</option>
          <option value="Reiterate">Reiterate</option>
        </select>
      </div>

      <!-- Rating Filter -->
      <div>
        <label for="rating" class="block text-sm font-medium text-gray-700 mb-1">
          Rating
        </label>
        <select
          id="rating"
          v-model="localFilters.rating"
          class="input-field"
          @change="updateFilters"
        >
          <option value="">All Ratings</option>
          <option value="Buy">Buy</option>
          <option value="Strong Buy">Strong Buy</option>
          <option value="Hold">Hold</option>
          <option value="Sell">Sell</option>
          <option value="Outperform">Outperform</option>
          <option value="Underperform">Underperform</option>
        </select>
      </div>

      <!-- Brokerage Filter -->
      <div>
        <label for="brokerage" class="block text-sm font-medium text-gray-700 mb-1">
          Brokerage
        </label>
        <select
          id="brokerage"
          v-model="localFilters.brokerage"
          class="input-field"
          @change="updateFilters"
        >
          <option value="">All Brokerages</option>
          <option v-for="brokerage in availableBrokerages" :key="brokerage" :value="brokerage">
            {{ brokerage }}
          </option>
        </select>
      </div>
    </div>

    <!-- Sort Options -->
    <div class="mt-4 flex flex-wrap items-center gap-4">
      <div class="flex items-center space-x-2">
        <label for="sortBy" class="text-sm font-medium text-gray-700">
          Sort by:
        </label>
        <select
          id="sortBy"
          v-model="localFilters.sortBy"
          class="input-field min-w-0 w-auto"
          @change="updateFilters"
        >
          <option value="time">Date</option>
          <option value="ticker">Ticker</option>
          <option value="company">Company</option>
          <option value="action">Action</option>
          <option value="rating_to">Rating</option>
        </select>
      </div>

      <div class="flex items-center space-x-2">
        <label for="sortOrder" class="text-sm font-medium text-gray-700">
          Order:
        </label>
        <select
          id="sortOrder"
          v-model="localFilters.sortOrder"
          class="input-field min-w-0 w-auto"
          @change="updateFilters"
        >
          <option value="desc">Descending</option>
          <option value="asc">Ascending</option>
        </select>
      </div>

      <!-- Clear Filters Button -->
      <button
        @click="clearFilters"
        class="btn-secondary text-sm"
      >
        Clear Filters
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { SearchFilters } from '@/types'

interface Props {
  filters: SearchFilters
  availableBrokerages: string[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:filters': [filters: SearchFilters]
}>()

const localFilters = ref<SearchFilters>({ ...props.filters })

// Debounce search input
let debounceTimeout: number | null = null

function debouncedUpdate() {
  if (debounceTimeout) {
    clearTimeout(debounceTimeout)
  }
  debounceTimeout = setTimeout(() => {
    updateFilters()
  }, 300)
}

function updateFilters() {
  emit('update:filters', { ...localFilters.value })
}

function clearFilters() {
  localFilters.value = {
    search: '',
    action: '',
    rating: '',
    brokerage: '',
    sortBy: 'time',
    sortOrder: 'desc'
  }
  updateFilters()
}

// Watch for external filter changes
watch(
  () => props.filters,
  (newFilters) => {
    localFilters.value = { ...newFilters }
  },
  { deep: true }
)

onMounted(() => {
  localFilters.value = { ...props.filters }
})
</script>