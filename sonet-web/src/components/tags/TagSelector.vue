<template>
  <div class="relative" ref="container">
    <button
      @click="open = !open"
      class="flex items-center gap-1 px-2 py-1 rounded-lg text-xs text-text-tertiary hover:text-text-secondary hover:bg-bg-hover transition-all"
    >
      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z" />
      </svg>
      Теги
    </button>

    <div
      v-if="open"
      class="absolute top-full left-0 mt-1 w-56 bg-bg-surface border border-border rounded-xl shadow-lg z-50 p-2 anim-fade-up"
    >
      <!-- Create new tag -->
      <div class="flex gap-1 mb-2">
        <input
          v-model="newTagName"
          @keydown.enter="createTag"
          placeholder="Новый тег..."
          class="flex-1 px-2.5 py-1.5 text-xs rounded-lg border border-border bg-bg-base text-text-primary placeholder-text-tertiary focus:outline-none focus:border-primary transition"
        />
        <button
          v-if="newTagName.trim()"
          @click="createTag"
          class="px-2 py-1.5 text-xs bg-primary text-text-inverse rounded-lg hover:bg-primary-hover transition"
        >
          +
        </button>
      </div>

      <!-- Tag list -->
      <div class="max-h-40 overflow-auto space-y-0.5">
        <button
          v-for="tag in tagStore.tags"
          :key="tag.id"
          @click="toggleTag(tag.id)"
          :class="[
            'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs transition-all',
            selectedTagIds.includes(tag.id)
              ? 'bg-bg-active text-text-primary'
              : 'text-text-secondary hover:bg-bg-hover',
          ]"
        >
          <span
            class="w-2 h-2 rounded-full shrink-0"
            :style="{ backgroundColor: tag.color || '#A89888' }"
          ></span>
          <span class="flex-1 text-left truncate">{{ tag.name }}</span>
          <svg
            v-if="selectedTagIds.includes(tag.id)"
            class="w-3 h-3 text-primary shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </button>
      </div>

      <p v-if="!tagStore.tags.length" class="text-center text-[11px] text-text-tertiary py-2">
        Нет тегов
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useTagStore } from '@/stores/tag'

const props = defineProps<{
  workspaceId: number
  noteId: number
  selectedTagIds: number[]
}>()

const emit = defineEmits<{
  toggle: [tagId: number]
  created: []
}>()

const tagStore = useTagStore()
const open = ref(false)
const newTagName = ref('')
const container = ref<HTMLElement>()

async function createTag() {
  if (!newTagName.value.trim()) return
  await tagStore.createTag(props.workspaceId, newTagName.value.trim())
  newTagName.value = ''
  emit('created')
}

function toggleTag(tagId: number) {
  emit('toggle', tagId)
}

function handleClickOutside(e: MouseEvent) {
  if (container.value && !container.value.contains(e.target as Node)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>
