<template>
  <button
    :title="title"
    :disabled="disabled"
    @click="$emit('click')"
    :class="[
      'w-7 h-7 flex items-center justify-center rounded-md transition-all duration-100',
      active
        ? 'bg-primary/15 text-primary'
        : disabled
          ? 'text-text-tertiary/40 cursor-not-allowed'
          : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
    ]"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  active?: boolean
  disabled?: boolean
}>()

defineEmits<{
  click: []
}>()
</script>
