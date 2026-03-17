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

      <!-- Meta bar -->
      <div class="flex items-center gap-3 mt-3 flex-wrap">
        <span class="text-[11px] text-text-tertiary">{{ note.word_count }} слов</span>
        <span class="w-1 h-1 rounded-full bg-border"></span>
        <span class="text-[11px] text-text-tertiary">{{ formatDate(note.updated_at) }}</span>
        <span class="w-1 h-1 rounded-full bg-border"></span>

        <button
          @click="togglePin"
          :class="['text-[11px] transition-colors', note.is_pinned ? 'text-primary' : 'text-text-tertiary hover:text-primary']"
        >
          {{ note.is_pinned ? '&#9733; Закреплено' : '&#9734; Закрепить' }}
        </button>

        <!-- Tags -->
        <span class="w-1 h-1 rounded-full bg-border"></span>
        <div class="flex items-center gap-1 flex-wrap">
          <TagPill
            v-for="tag in noteTags"
            :key="tag.id"
            :tag="tag"
            removable
            @remove="detachTag(tag.id)"
          />
          <TagSelector
            v-if="workspaceStore.currentWorkspace"
            :workspace-id="workspaceStore.currentWorkspace.id"
            :note-id="note.id"
            :selected-tag-ids="noteTags.map(t => t.id)"
            @toggle="handleTagToggle"
            @created="loadTags"
          />
        </div>

        <div class="flex-1"></div>

        <span v-if="saving" class="text-[11px] text-text-tertiary anim-fade-in">Сохранение...</span>
        <span v-else-if="saved" class="text-[11px] text-success anim-fade-in">Сохранено</span>

        <button
          @click="handleDelete"
          class="text-[11px] text-text-tertiary hover:text-error transition-colors"
        >
          Удалить
        </button>
      </div>
      <div class="h-px bg-gradient-to-r from-border via-border-light to-transparent mt-4"></div>
    </div>

    <!-- Editor + Backlinks -->
    <div class="flex-1 px-8 pb-8 overflow-auto">
      <NoteEditor
        :content="initialContent"
        @update="handleEditorUpdate"
      />

      <BacklinksPanel
        v-if="workspaceStore.currentWorkspace"
        :workspace-id="workspaceStore.currentWorkspace.id"
        :note-id="note.id"
        @navigate="navigateToNote"
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
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useNoteStore } from '@/stores/note'
import { useTagStore } from '@/stores/tag'
import { useDebounce } from '@/composables/useDebounce'
import api from '@/composables/useApi'
import NoteEditor from '@/components/editor/NoteEditor.vue'
import BacklinksPanel from '@/components/editor/BacklinksPanel.vue'
import TagPill from '@/components/tags/TagPill.vue'
import TagSelector from '@/components/tags/TagSelector.vue'
import type { Tag } from '@/types/tag'

const route = useRoute()
const router = useRouter()
const workspaceStore = useWorkspaceStore()
const noteStore = useNoteStore()
const tagStore = useTagStore()

const note = ref(noteStore.currentNote)
const initialContent = ref<any>(null)
const noteTags = ref<Tag[]>([])
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

async function loadNoteTags() {
  if (!note.value || !workspaceStore.currentWorkspace) return
  try {
    const { data } = await api.get(`/workspaces/${workspaceStore.currentWorkspace.id}/notes/${note.value.id}/tags`)
    noteTags.value = data.data || []
  } catch {
    noteTags.value = []
  }
}

async function loadTags() {
  if (!workspaceStore.currentWorkspace) return
  await tagStore.fetchTags(workspaceStore.currentWorkspace.id)
}

async function handleTagToggle(tagId: number) {
  if (!note.value || !workspaceStore.currentWorkspace) return
  const wsId = workspaceStore.currentWorkspace.id
  const isAttached = noteTags.value.some(t => t.id === tagId)

  if (isAttached) {
    await api.delete(`/workspaces/${wsId}/notes/${note.value.id}/tags/${tagId}`)
  } else {
    await api.post(`/workspaces/${wsId}/notes/${note.value.id}/tags`, { tag_ids: [tagId] })
  }
  await loadNoteTags()
}

async function detachTag(tagId: number) {
  if (!note.value || !workspaceStore.currentWorkspace) return
  await api.delete(`/workspaces/${workspaceStore.currentWorkspace.id}/notes/${note.value.id}/tags/${tagId}`)
  await loadNoteTags()
}

function navigateToNote(noteId: number) {
  router.push(`/w/${route.params.wsSlug}/note/${noteId}`)
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function loadNote() {
  if (!workspaceStore.currentWorkspace) return
  const noteId = Number(route.params.noteId)
  await noteStore.fetchNote(workspaceStore.currentWorkspace.id, noteId)
  note.value = noteStore.currentNote
  initialContent.value = note.value?.content_json || note.value?.content_html || ''
  await loadNoteTags()
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  await workspaceStore.setCurrentBySlug(slug)
  if (workspaceStore.currentWorkspace) {
    await tagStore.fetchTags(workspaceStore.currentWorkspace.id)
    await loadNote()
  }
})

watch(() => route.params.noteId, async (newId) => {
  if (newId && workspaceStore.currentWorkspace) {
    await loadNote()
  }
})
</script>
