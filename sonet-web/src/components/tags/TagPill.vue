<template>
  <span
    :class="['inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-medium transition-colors', clickable ? 'cursor-pointer hover:opacity-80' : '']"
    :style="pillStyle"
    @click="$emit('click')"
  >
    {{ tag.name }}
    <button
      v-if="removable"
      @click.stop="$emit('remove')"
      class="ml-0.5 opacity-60 hover:opacity-100 transition-opacity"
    >
      <svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tag } from '@/types/tag'

const props = defineProps<{
  tag: Tag
  removable?: boolean
  clickable?: boolean
}>()

defineEmits<{
  click: []
  remove: []
}>()

const defaultColors = ['#E8A87C', '#C75B39', '#D4A843', '#5B8C5A', '#6B8FAD', '#9B7DB8', '#C4786F', '#8B6C4F']

const pillStyle = computed(() => {
  const color = props.tag.color || defaultColors[props.tag.id % defaultColors.length]
  return {
    backgroundColor: `${color}18`,
    color: color,
  }
})
</script>
