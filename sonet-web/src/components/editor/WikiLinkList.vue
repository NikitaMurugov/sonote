<template>
  <div
    v-if="items.length"
    class="bg-bg-surface border border-border rounded-xl shadow-lg overflow-hidden min-w-[200px] max-w-[320px]"
  >
    <div class="px-3 py-1.5 border-b border-border-light">
      <p class="text-[10px] text-text-tertiary uppercase tracking-wide font-semibold">Ссылка на заметку</p>
    </div>
    <div class="max-h-48 overflow-auto py-1">
      <button
        v-for="(item, index) in items"
        :key="item.id"
        @click="selectItem(index)"
        :class="[
          'w-full text-left px-3 py-2 text-sm transition-colors',
          index === selectedIndex
            ? 'bg-bg-active text-text-primary'
            : 'text-text-secondary hover:bg-bg-hover',
        ]"
      >
        <span class="truncate block">{{ item.title || 'Без названия' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  items: Array<{ id: number; title: string }>
  command: (item: any) => void
}>()

const selectedIndex = ref(0)

watch(
  () => props.items,
  () => {
    selectedIndex.value = 0
  }
)

function selectItem(index: number) {
  const item = props.items[index]
  if (item) {
    props.command(item)
  }
}

defineExpose({
  onKeyDown(event: KeyboardEvent) {
    if (event.key === 'ArrowUp') {
      selectedIndex.value = (selectedIndex.value + props.items.length - 1) % props.items.length
      return true
    }
    if (event.key === 'ArrowDown') {
      selectedIndex.value = (selectedIndex.value + 1) % props.items.length
      return true
    }
    if (event.key === 'Enter') {
      selectItem(selectedIndex.value)
      return true
    }
    return false
  },
})
</script>
