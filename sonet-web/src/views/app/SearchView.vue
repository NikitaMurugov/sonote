<template>
  <div class="px-10 py-10 max-w-3xl mx-auto anim-fade-up">
    <h1 class="text-2xl font-semibold mb-8" style="font-family: var(--font-heading)">Поиск</h1>

    <div class="relative mb-8">
      <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-text-tertiary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        v-model="query"
        @input="debouncedSearch"
        autofocus
        placeholder="Поиск по заметкам..."
        class="w-full pl-11 pr-4 py-3 rounded-xl border border-border bg-bg-surface text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
      />
    </div>

    <div v-if="results.length" class="space-y-3">
      <div
        v-for="note in results"
        :key="note.id"
        @click="openNote(note.id)"
        class="p-4 bg-bg-surface/60 border border-transparent rounded-2xl cursor-pointer hover:bg-bg-surface hover:border-border hover:shadow-[0_2px_12px_rgba(45,32,24,0.05)] transition-all duration-200"
      >
        <h3 class="font-medium text-text-primary text-[15px] mb-0.5" style="font-family: var(--font-heading); font-size: 17px;">
          {{ note.title }}
        </h3>
        <p class="text-sm text-text-secondary/80 line-clamp-2">{{ note.content_md.slice(0, 200) }}</p>
      </div>
    </div>

    <p v-else-if="query && !loading" class="text-text-tertiary text-sm text-center py-8">
      Ничего не найдено
    </p>

    <p v-else-if="!query" class="text-text-tertiary text-sm text-center py-8">
      Введите запрос для поиска
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useDebounce } from '@/composables/useDebounce'
import api from '@/composables/useApi'
import type { Note } from '@/types/note'

const route = useRoute()
const router = useRouter()
const workspaceStore = useWorkspaceStore()

const query = ref('')
const results = ref<Note[]>([])
const loading = ref(false)

const debouncedSearch = useDebounce(async () => {
  if (!query.value.trim() || !workspaceStore.currentWorkspace) {
    results.value = []
    return
  }
  loading.value = true
  try {
    const { data } = await api.get(`/workspaces/${workspaceStore.currentWorkspace.id}/search`, {
      params: { q: query.value },
    })
    results.value = data.data || []
  } catch {
    results.value = []
  } finally {
    loading.value = false
  }
}, 300)

function openNote(noteId: number) {
  router.push(`/w/${route.params.wsSlug}/note/${noteId}`)
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  if (slug) {
    await workspaceStore.setCurrentBySlug(slug)
  }
})
</script>
