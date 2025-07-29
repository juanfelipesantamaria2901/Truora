<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <h1 class="text-2xl font-bold text-primary-600">Truora Stocks</h1>
          </div>
        </div>

        <!-- Navigation -->
        <nav class="hidden md:flex space-x-8">
          <RouterLink
            to="/"
            class="text-gray-600 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/' }"
          >
            Dashboard
          </RouterLink>
          <RouterLink
            to="/stocks"
            class="text-gray-600 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/stocks' }"
          >
            Stocks
          </RouterLink>
          <RouterLink
            to="/recommendations"
            class="text-gray-600 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/recommendations' }"
          >
            Recommendations
          </RouterLink>
        </nav>

        <!-- Mobile menu button -->
        <div class="md:hidden">
          <button
            @click="toggleMobileMenu"
            class="text-gray-600 hover:text-primary-600 focus:outline-none focus:text-primary-600"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                v-if="!showMobileMenu"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
              <path
                v-else
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </div>

      <!-- Mobile menu -->
      <div v-if="showMobileMenu" class="md:hidden">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
          <RouterLink
            to="/"
            @click="showMobileMenu = false"
            class="text-gray-600 hover:text-primary-600 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/' }"
          >
            Dashboard
          </RouterLink>
          <RouterLink
            to="/stocks"
            @click="showMobileMenu = false"
            class="text-gray-600 hover:text-primary-600 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/stocks' }"
          >
            Stocks
          </RouterLink>
          <RouterLink
            to="/recommendations"
            @click="showMobileMenu = false"
            class="text-gray-600 hover:text-primary-600 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'text-primary-600 bg-primary-50': $route.path === '/recommendations' }"
          >
            Recommendations
          </RouterLink>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'

const showMobileMenu = ref(false)

function toggleMobileMenu() {
  showMobileMenu.value = !showMobileMenu.value
}
</script>