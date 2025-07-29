<template>
  <div class="flex items-center justify-center" :class="containerClass">
    <div class="animate-spin rounded-full border-4 border-gray-200" :class="spinnerClass">
      <div class="rounded-full" :class="innerClass"></div>
    </div>
    <span v-if="text" class="ml-3 text-gray-600" :class="textClass">{{ text }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
  text?: string
  color?: 'primary' | 'gray'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'primary'
})

const containerClass = computed(() => {
  const classes = ['loading-spinner']
  if (props.size === 'sm') classes.push('p-2')
  else if (props.size === 'lg') classes.push('p-8')
  else classes.push('p-4')
  return classes.join(' ')
})

const spinnerClass = computed(() => {
  const classes = []
  
  // Size classes
  if (props.size === 'sm') {
    classes.push('h-4 w-4 border-2')
  } else if (props.size === 'lg') {
    classes.push('h-12 w-12 border-4')
  } else {
    classes.push('h-8 w-8 border-4')
  }
  
  // Color classes
  if (props.color === 'primary') {
    classes.push('border-t-primary-600')
  } else {
    classes.push('border-t-gray-600')
  }
  
  return classes.join(' ')
})

const innerClass = computed(() => {
  const classes = []
  
  if (props.size === 'sm') {
    classes.push('h-3 w-3')
  } else if (props.size === 'lg') {
    classes.push('h-11 w-11')
  } else {
    classes.push('h-7 w-7')
  }
  
  return classes.join(' ')
})

const textClass = computed(() => {
  if (props.size === 'sm') return 'text-sm'
  if (props.size === 'lg') return 'text-lg'
  return 'text-base'
})
</script>