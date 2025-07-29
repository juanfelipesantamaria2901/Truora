# Truora Frontend - Detailed Documentation

## 🎯 Architecture Overview

The Truora frontend is a modern Vue.js 3 application built with TypeScript, featuring a component-based architecture, reactive state management with Pinia, and a responsive design system powered by Tailwind CSS. The application follows Vue 3's Composition API patterns and implements best practices for scalable frontend development.

## 📁 Detailed Project Structure

```
frontend/
├── src/                          # Source code directory
│   ├── components/               # Reusable Vue components
│   │   ├── Layout/               # Layout-specific components
│   │   │   └── AppHeader.vue     # Application header with navigation
│   │   ├── Stocks/               # Stock-specific components
│   │   │   ├── StockCard.vue     # Individual stock display card
│   │   │   └── StockFilters.vue  # Stock filtering controls
│   │   ├── UI/                   # Generic UI components
│   │   │   └── LoadingSpinner.vue # Loading state component
│   │   └── icons/                # SVG icon components
│   │       ├── IconCommunity.vue
│   │       ├── IconDocumentation.vue
│   │       ├── IconEcosystem.vue
│   │       ├── IconSupport.vue
│   │       └── IconTooling.vue
│   ├── views/                  # Page-level components (routes)
│   │   ├── DashboardView.vue     # Main dashboard/home page
│   │   ├── StocksView.vue        # Stock listing and search
│   │   ├── RecommendationsView.vue # AI recommendations
│   │   └── AboutView.vue         # Application information
│   ├── services/               # API service layer
│   │   └── api.ts              # Axios instance and API methods
│   ├── stores/                 # Pinia state management
│   │   ├── stocks.ts           # Stock-related state
│   │   └── counter.ts          # Demo counter store
│   ├── types/                  # TypeScript type definitions
│   │   └── index.ts            # Shared interfaces and types
│   ├── router/                 # Vue Router configuration
│   │   └── index.ts            # Route definitions and guards
│   ├── assets/                 # Static assets
│   │   ├── base.css            # Base CSS styles
│   │   ├── main.css            # Main Tailwind CSS imports
│   │   └── logo.svg            # Application logo
│   ├── main.ts                 # Application entry point
│   └── App.vue                 # Root application component
├── public/                     # Static public files
│   └── favicon.ico             # Application favicon
├── package.json                # Node.js dependencies
├── package-lock.json           # Locked dependency versions
├── vite.config.ts             # Vite build configuration
├── tsconfig.json              # TypeScript configuration
├── tsconfig.app.json          # Application TypeScript config
├── tsconfig.node.json         # Node.js TypeScript config
├── tailwind.config.js         # Tailwind CSS configuration
├── postcss.config.js          # PostCSS configuration
├── eslint.config.ts           # ESLint configuration
├── .prettierrc.json           # Prettier formatting rules
├── .gitignore                 # Git ignore patterns
├── .editorconfig              # Editor configuration
├── .vscode/                   # VS Code workspace settings
│   ├── extensions.json        # Recommended extensions
│   └── settings.json          # Workspace settings
└── README.md                  # Basic frontend documentation
```

## 🛠️ Technology Stack Deep Dive

### Core Technologies
- **Vue.js 3.5.17**: Progressive JavaScript framework
- **TypeScript 5.8**: Type-safe JavaScript
- **Vite 7.0**: Next-generation build tool
- **Pinia 3.0.3**: Intuitive state management
- **Vue Router 4.5.1**: Official routing library

### Styling & UI
- **Tailwind CSS 4.1.11**: Utility-first CSS framework
- **PostCSS 8.5.6**: CSS transformation
- **Autoprefixer**: Vendor prefix management

### Development Tools
- **ESLint 9.29.0**: Code linting and quality
- **Prettier 3.5.3**: Code formatting
- **Vue DevTools**: Browser debugging extension
- **Vue TSC**: TypeScript checking for Vue

## 🎯 Component Architecture

### 1. Root Component (`App.vue`)

The root component provides the application shell:

```vue
<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <AppHeader />
    <RouterView />
  </div>
</template>

<script setup lang="ts">
// Global app configuration
import AppHeader from '@/components/Layout/AppHeader.vue'
</script>
```

### 2. Layout Components (`components/Layout/`)

#### AppHeader.vue

Responsive navigation header with mobile menu support:

```vue
<template>
  <header class="bg-white shadow-sm">
    <nav class="container mx-auto px-4 py-4">
      <!-- Logo and navigation -->
      <div class="flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold text-blue-600">
          Truora
        </RouterLink>
        
        <!-- Desktop navigation -->
        <div class="hidden md:flex space-x-8">
          <RouterLink to="/" class="nav-link">Dashboard</RouterLink>
          <RouterLink to="/stocks" class="nav-link">Stocks</RouterLink>
          <RouterLink to="/recommendations" class="nav-link">Recommendations</RouterLink>
        </div>
        
        <!-- Mobile menu button -->
        <button @click="toggleMobileMenu" class="md:hidden">
          <!-- Hamburger icon -->
        </button>
      </div>
    </nav>
  </header>
</template>
```

### 3. Stock Components (`components/Stocks/`)

#### StockCard.vue

Reusable stock display component:

```vue
<template>
  <div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
    <div class="flex justify-between items-start mb-4">
      <div>
        <h3 class="text-lg font-semibold">{{ stock.symbol }}</h3>
        <p class="text-gray-600">{{ stock.companyName }}</p>
      </div>
      <div class="text-right">
        <p class="text-2xl font-bold">${{ stock.currentPrice }}</p>
        <p class="text-sm" :class="priceChangeColor">
          {{ priceChangePercentage }}%
        </p>
      </div>
    </div>
    
    <div class="grid grid-cols-2 gap-4 text-sm">
      <div>
        <span class="text-gray-500">Market Cap:</span>
        <span class="font-medium">{{ formatMarketCap(stock.marketCap) }}</span>
      </div>
      <div>
        <span class="text-gray-500">P/E Ratio:</span>
        <span class="font-medium">{{ stock.pe }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Stock } from '@/types'
import { computed } from 'vue'

interface Props {
  stock: Stock
}

const props = defineProps<Props>()

const priceChangeColor = computed(() => 
  props.stock.priceChange >= 0 ? 'text-green-600' : 'text-red-600'
)
</script>
```

#### StockFilters.vue

Advanced filtering component:

```vue
<template>
  <div class="bg-white rounded-lg shadow p-4">
    <h3 class="font-semibold mb-4">Filters</h3>
    
    <div class="space-y-4">
      <!-- Search -->
      <div>
        <label class="block text-sm font-medium mb-2">Search</label>
        <input
          v-model="searchQuery"
          @input="updateFilters"
          type="text"
          placeholder="Symbol or company..."
          class="w-full px-3 py-2 border rounded-md"
        />
      </div>
      
      <!-- Sector filter -->
      <div>
        <label class="block text-sm font-medium mb-2">Sector</label>
        <select
          v-model="selectedSector"
          @change="updateFilters"
          class="w-full px-3 py-2 border rounded-md"
        >
          <option value="">All sectors</option>
          <option v-for="sector in sectors" :key="sector" :value="sector">
            {{ sector }}
          </option>
        </select>
      </div>
      
      <!-- Market cap range -->
      <div>
        <label class="block text-sm font-medium mb-2">Market Cap</label>
        <div class="flex space-x-2">
          <input
            v-model.number="minMarketCap"
            @input="updateFilters"
            type="number"
            placeholder="Min"
            class="flex-1 px-3 py-2 border rounded-md"
          />
          <input
            v-model.number="maxMarketCap"
            @input="updateFilters"
            type="number"
            placeholder="Max"
            class="flex-1 px-3 py-2 border rounded-md"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  filterChange: [filters: FilterOptions]
}>()

const searchQuery = ref('')
const selectedSector = ref('')
const minMarketCap = ref<number | null>(null)
const maxMarketCap = ref<number | null>(null)

const updateFilters = () => {
  emit('filterChange', {
    search: searchQuery.value,
    sector: selectedSector.value,
    minMarketCap: minMarketCap.value,
    maxMarketCap: maxMarketCap.value
  })
}
</script>
```

### 4. UI Components (`components/UI/`)

#### LoadingSpinner.vue

Reusable loading state component:

```vue
<template>
  <div class="flex justify-center items-center py-8">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
  </div>
</template>
```

## 🗂️ TypeScript Types (`types/index.ts`)

### Core Type Definitions

```typescript
// Stock data types
export interface Stock {
  id: number
  symbol: string
  companyName: string
  currentPrice: number
  marketCap: number
  pe: number
  dividendYield: number
  week52High: number
  week52Low: number
  volume: number
  sector: string
  industry: string
  createdAt: string
  updatedAt: string
}

// Recommendation types
export interface StockRecommendation {
  id: number
  stockId: number
  stock: Stock
  score: number
  reason: string
  riskLevel: 'low' | 'medium' | 'high'
  timeHorizon: 'short' | 'medium' | 'long'
  createdAt: string
}

// API response types
export interface ApiResponse<T> {
  data: T
  message?: string
  pagination?: PaginationInfo
}

export interface PaginationInfo {
  total: number
  limit: number
  offset: number
  hasMore: boolean
}

// Filter types
export interface FilterOptions {
  search?: string
  sector?: string
  minMarketCap?: number
  maxMarketCap?: number
  limit?: number
  offset?: number
}

// Error handling
export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
}
```

## 🔄 State Management (Pinia Stores)

### Stocks Store (`stores/stocks.ts`)

Centralized stock data management:

```typescript
import { defineStore } from 'pinia'
import type { Stock, StockRecommendation, FilterOptions } from '@/types'
import { stocksApi } from '@/services/api'

interface StocksState {
  stocks: Stock[]
  recommendations: StockRecommendation[]
  loading: boolean
  error: string | null
  filters: FilterOptions
  pagination: {
    total: number
    limit: number
    offset: number
  }
}

export const useStocksStore = defineStore('stocks', {
  state: (): StocksState => ({
    stocks: [],
    recommendations: [],
    loading: false,
    error: null,
    filters: {
      limit: 20,
      offset: 0
    },
    pagination: {
      total: 0,
      limit: 20,
      offset: 0
    }
  }),

  getters: {
    filteredStocks: (state) => {
      // Apply client-side filtering if needed
      return state.stocks
    },
    
    topRecommendations: (state) => {
      return state.recommendations
        .sort((a, b) => b.score - a.score)
        .slice(0, 10)
    }
  },

  actions: {
    async fetchStocks() {
      this.loading = true
      this.error = null
      
      try {
        const response = await stocksApi.getStocks(this.filters)
        this.stocks = response.data.data
        this.pagination = response.data.pagination
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch stocks'
      } finally {
        this.loading = false
      }
    },

    async fetchRecommendations(limit = 10) {
      this.loading = true
      this.error = null
      
      try {
        const response = await stocksApi.getRecommendations({ limit })
        this.recommendations = response.data.data
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch recommendations'
      } finally {
        this.loading = false
      }
    },

    async fetchStockBySymbol(symbol: string): Promise<Stock | null> {
      try {
        const response = await stocksApi.getStock(symbol)
        return response.data.data
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch stock'
        return null
      }
    },

    updateFilters(filters: Partial<FilterOptions>) {
      this.filters = { ...this.filters, ...filters }
      this.fetchStocks()
    },

    clearFilters() {
      this.filters = { limit: 20, offset: 0 }
      this.fetchStocks()
    }
  }
})
```

## 🔌 API Service Layer (`services/api.ts`)

### Axios Configuration

```typescript
import axios from 'axios'
import type { Stock, StockRecommendation, ApiResponse, FilterOptions } from '@/types'

// Create axios instance with default configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor for auth tokens
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('auth_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Stock API methods
export const stocksApi = {
  async getStocks(filters: FilterOptions = {}) {
    const params = new URLSearchParams()
    
    if (filters.limit) params.append('limit', filters.limit.toString())
    if (filters.offset) params.append('offset', filters.offset.toString())
    if (filters.search) params.append('q', filters.search)
    if (filters.sector) params.append('sector', filters.sector)
    
    return api.get<ApiResponse<Stock[]>>(`/stocks?${params}`)
  },

  async getStock(symbol: string) {
    return api.get<ApiResponse<Stock>>(`/stocks/${symbol}`)
  },

  async fetchStocks() {
    return api.post('/stocks/fetch')
  },

  async getRecommendations(filters: { limit?: number } = {}) {
    const params = new URLSearchParams()
    if (filters.limit) params.append('limit', filters.limit.toString())
    
    return api.get<ApiResponse<StockRecommendation[]>>(`/recommendations?${params}`)
  },

  async generateRecommendations() {
    return api.post('/recommendations/generate')
  }
}

// Health check API
export const healthApi = {
  async checkHealth() {
    return api.get('/health')
  }
}
```

## 🛣️ Routing Configuration (`router/index.ts`)

### Vue Router Setup

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: {
      title: 'Dashboard - Truora'
    }
  },
  {
    path: '/stocks',
    name: 'Stocks',
    component: () => import('@/views/StocksView.vue'),
    meta: {
      title: 'Stocks - Truora'
    }
  },
  {
    path: '/stocks/:symbol',
    name: 'StockDetail',
    component: () => import('@/views/StockDetailView.vue'),
    props: true,
    meta: {
      title: 'Stock Details - Truora'
    }
  },
  {
    path: '/recommendations',
    name: 'Recommendations',
    component: () => import('@/views/RecommendationsView.vue'),
    meta: {
      title: 'Recommendations - Truora'
    }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/AboutView.vue'),
    meta: {
      title: 'About - Truora'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue'),
    meta: {
      title: 'Page Not Found - Truora'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // Scroll to top on route change, except for back/forward navigation
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Navigation guards
router.beforeEach((to, from, next) => {
  // Set page title
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
  
  // Authentication guard (future)
  // if (to.meta.requiresAuth && !isAuthenticated()) {
  //   next('/login')
  // } else {
  //   next()
  // }
  
  next()
})

export default router
```

## 🎨 Styling System

### Tailwind CSS Configuration

#### `tailwind.config.js`
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          900: '#1e3a8a',
        },
        secondary: {
          50: '#f8fafc',
          100: '#f1f5f9',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          900: '#0f172a',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
      }
    },
  },
  plugins: [],
}
```

#### Global Styles (`assets/main.css`)
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  html {
    font-family: 'Inter', system-ui, sans-serif;
  }
  
  body {
    @apply bg-gray-50 text-gray-900;
  }
}

@layer components {
  .btn-primary {
    @apply bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors;
  }
  
  .btn-secondary {
    @apply bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium py-2 px-4 rounded-md transition-colors;
  }
  
  .card {
    @apply bg-white rounded-lg shadow-sm border border-gray-200;
  }
}

@layer utilities {
  .text-gradient {
    @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
  }
}
```

## 🧪 Testing Strategy

### Component Testing

#### Unit Tests with Vue Test Utils
```typescript
// Example: StockCard.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StockCard from '@/components/Stocks/StockCard.vue'

describe('StockCard', () => {
  it('renders stock information correctly', () => {
    const stock = {
      symbol: 'AAPL',
      companyName: 'Apple Inc.',
      currentPrice: 150.00,
      marketCap: 2500000000000,
      pe: 25.5
    }
    
    const wrapper = mount(StockCard, {
      props: { stock }
    })
    
    expect(wrapper.text()).toContain('AAPL')
    expect(wrapper.text()).toContain('$150.00')
  })
})
```

#### End-to-End Tests with Cypress
```typescript
// Example: stocks.cy.ts
describe('Stocks Page', () => {
  it('displays stock list', () => {
    cy.visit('/stocks')
    cy.contains('Stocks').should('be.visible')
    cy.get('[data-testid="stock-card"]').should('have.length.greaterThan', 0)
  })
  
  it('filters stocks by search', () => {
    cy.visit('/stocks')
    cy.get('[data-testid="search-input"]').type('AAPL')
    cy.get('[data-testid="stock-card"]').should('contain', 'AAPL')
  })
})
```

## 🚀 Performance Optimization

### Code Splitting
```typescript
// Lazy loading routes
const StocksView = () => import('@/views/StocksView.vue')
const RecommendationsView = () => import('@/views/RecommendationsView.vue')
```

### Image Optimization
```vue
<!-- Optimized images with lazy loading -->
<img 
  :src="stock.logo" 
  :alt="stock.companyName"
  loading="lazy"
  class="w-12 h-12 rounded-full"
/>
```

### Bundle Analysis
```bash
# Analyze bundle size
npm run build
npm run preview
```

## 🔒 Security Best Practices

### Content Security Policy (CSP)
```html
<meta 
  http-equiv="Content-Security-Policy" 
  content="default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
/>
```

### Input Validation
```typescript
// Sanitize user input
import DOMPurify from 'dompurify'

const sanitizeInput = (input: string): string => {
  return DOMPurify.sanitize(input, { ALLOWED_TAGS: [] })
}
```

## 🚀 Deployment

### Build Configuration

#### Vite Build Settings (`vite.config.ts`)
```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    target: 'es2015',
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          ui: ['@headlessui/vue', '@heroicons/vue']
        }
      }
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
```

### Production Deployment

#### Vercel Configuration (`vercel.json`)
```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "devCommand": "npm run dev",
  "installCommand": "npm install"
}
```

#### Netlify Configuration (`netlify.toml`)
```toml
[build]
  command = "npm run build"
  publish = "dist"

[build.environment]
  NODE_VERSION = "18"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

## 📚 Additional Resources

### Learning Resources
- [Vue.js 3 Documentation](https://vuejs.org/guide/introduction.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Pinia Documentation](https://pinia.vuejs.org/introduction.html)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Vue Router Guide](https://router.vuejs.org/guide/)

### Development Tools
- **Vue DevTools**: Browser extension for debugging
- **TypeScript Vue Plugin**: VS Code extension
- **ESLint**: Code quality and consistency
- **Prettier**: Code formatting

---

**Next Steps**: Continue to the [Backend Detailed Documentation](../backend/DETAILED_README.md) for complete system understanding.