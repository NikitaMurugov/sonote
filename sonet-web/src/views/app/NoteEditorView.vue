<template>
  <div v-if="note" class="h-full flex flex-col anim-fade-up">
    <!-- Title area -->
    <div class="px-8 pt-8 pb-2">
      <input
        v-model="note.title"
        @input="debouncedSaveTitle"
        class="w-full text-3xl font-semibold bg-transparent border-none outline-none text-text-primary placeholder-text-tertiary/50 tracking-tight"
        style="font-family: var(--font-heading)"
        placeholder="Без названия"
      />
      <div class="flex items-center gap-4 mt-3">
        <span class="text-[11px] text-text-tertiary">{{ note.word_count }} слов</span>
        <span class="w-1 h-1 rounded-full bg-border"></span>
        <span class="text-[11px] text-text-tertiary">{{ formatDate(note.updated_at) }}</span>
        <span class="w-1 h-1 rounded-full bg-border"></span>

        <!-- Pin button -->
        <button
          @click="togglePin"
          :class="['text-[11px] transition-colors', note.is_pinned ? 'text-primary' : 'text-text-tertiary hover:text-primary']"
        >
          {{ note.is_pinned ? 'Закреплено' : 'Закрепить' }}
        </button>

        <div class="flex-1"></div>

        <!-- Save status -->
        <span v-if="saving" class="text-[11px] text-text-tertiary anim-fade-in">Сохранение...</span>
        <span v-else-if="saved" class="text-[11px] text-success anim-fade-in">Сохранено</span>

        <!-- Delete button -->
        <button
          @click="handleDelete"
          class="text-[11px] text-text-tertiary hover:text-error transition-colors"
        >
          Удалить
        </button>
      </div>
      <div class="h-px bg-gradient-to-r from-border via-border-light to-transparent mt-4"></div>
    </div>

    <!-- Editor -->
    <div class="flex-1 px-8 pb-8 overflow-auto">
      <NoteEditor
        :content="initialContent"
        @update="handleEditorUpdate"
      />
    </div>
  </div>

  <!-- Loading -->
  <div v-else class="flex items-center justify-center h-full">
    <div class="text-center anim-fade-in">
      <div class="w-6 h-6 border-2 border-primary/30 border-t-primary rounded-full animate-spin mx-auto mb-3"></div>
      <p class="text-text-tertiary text-sm">Загрузка...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useNoteStore } from '@/stores/note'
import { useDebounce } from '@/composables/useDebounce'
import NoteEditor from '@/components/editor/NoteEditor.vue'

const route = useRoute()
const router = useRouter()
const workspaceStore = useWorkspaceStore()
const noteStore = useNoteStore()

const note = ref(noteStore.currentNote)
const initialContent = ref<any>(null)
const saving = ref(false)
const saved = ref(false)
let savedTimeout: ReturnType<typeof setTimeout>

const debouncedSaveTitle = useDebounce(async () => {
  if (!note.value || !workspaceStore.currentWorkspace) return
  saving.value = true
  await noteStore.updateNote(workspaceStore.currentWorkspace.id, note.value.id, {
    title: note.value.title,
  })
  showSaved()
}, 500)

const debouncedSaveContent = useDebounce(async (updates: any) => {
  if (!note.value || !workspaceStore.currentWorkspace) return
  saving.value = true
  await noteStore.updateNote(workspaceStore.currentWorkspace.id, note.value.id, updates)
  showSaved()
}, 800)

function showSaved() {
  saving.value = false
  saved.value = true
  clearTimeout(savedTimeout)
  savedTimeout = setTimeout(() => { saved.value = false }, 2000)
}

function handleEditorUpdate({ html, json, text }: { html: string; json: any; text: string }) {
  if (!note.value) return
  note.value.word_count = text.split(/\s+/).filter(Boolean).length
  debouncedSaveContent({
    content_html: html,
    content_json: json,
    content_md: text,
  })
}

async function togglePin() {
  if (!note.value || !workspaceStore.currentWorkspace) return
  note.value.is_pinned = !note.value.is_pinned
  await noteStore.updateNote(workspaceStore.currentWorkspace.id, note.value.id, {
    is_pinned: note.value.is_pinned,
  })
}

async function handleDelete() {
  if (!note.value || !workspaceStore.currentWorkspace) return
  await noteStore.deleteNote(workspaceStore.currentWorkspace.id, note.value.id)
  router.push(`/w/${route.params.wsSlug}`)
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  await workspaceStore.setCurrentBySlug(slug)
  if (workspaceStore.currentWorkspace) {
    const noteId = Number(route.params.noteId)
    await noteStore.fetchNote(workspaceStore.currentWorkspace.id, noteId)
    note.value = noteStore.currentNote
    initialContent.value = note.value?.content_json || note.value?.content_html || ''
  }
})

watch(() => route.params.noteId, async (newId) => {
  if (newId && workspaceStore.currentWorkspace) {
    await noteStore.fetchNote(workspaceStore.currentWorkspace.id, Number(newId))
    note.value = noteStore.currentNote
    initialContent.value = note.value?.content_json || note.value?.content_html || ''
  }
})
</script>
