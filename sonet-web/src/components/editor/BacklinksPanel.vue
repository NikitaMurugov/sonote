<template>
  <div v-if="backlinks.length" class="mt-8 pt-6 border-t border-border-light">
    <div class="flex items-center gap-2 mb-3">
      <svg class="w-3.5 h-3.5 text-text-tertiary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
      </svg>
      <p class="text-[10px] font-semibold text-text-tertiary uppercase tracking-[0.15em]">
        Ссылаются сюда ({{ backlinks.length }})
      </p>
    </div>
    <div class="space-y-1.5">
      <button
        v-for="link in backlinks"
        :key="link.id"
        @click="$emit('navigate', link.source_note_id)"
        class="w-full text-left p-3 rounded-lg bg-bg-surface/50 border border-transparent hover:border-border hover:bg-bg-surface transition-all duration-150 group"
      >
        <p class="text-sm text-text-primary group-hover:text-primary transition-colors font-medium">
          Заметка #{{ link.source_note_id }}
        </p>
        <p v-if="link.context_snippet" class="text-xs text-text-tertiary mt-0.5 line-clamp-2">
          {{ link.context_snippet }}
        </p>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import api from '@/composables/useApi'

const props = defineProps<{
  workspaceId: number
  noteId: number
}>()

defineEmits<{
  navigate: [noteId: number]
}>()

interface NoteLink {
  id: number
  source_note_id: number
  target_note_id: number
  context_snippet: string | null
}

const backlinks = ref<NoteLink[]>([])

async function fetchBacklinks() {
  try {
    const { data } = await api.get(`/workspaces/${props.workspaceId}/notes/${props.noteId}/backlinks`)
    backlinks.value = data.data || []
  } catch {
    backlinks.value = []
  }
}

onMounted(fetchBacklinks)
watch(() => props.noteId, fetchBacklinks)
</script>
