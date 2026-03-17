<template>
  <div class="px-10 py-10 max-w-5xl mx-auto">
    <!-- Greeting -->
    <div class="mb-10 anim-fade-up">
      <h1 class="text-3xl font-semibold text-text-primary tracking-tight" style="font-family: var(--font-heading)">
        {{ greeting }}, {{ authStore.user?.display_name || 'друг' }}
      </h1>
      <p class="text-text-tertiary text-sm mt-1.5 tracking-wide">{{ formattedDate }}</p>
    </div>

    <!-- Quick action -->
    <div class="mb-12 anim-fade-up d2">
      <button
        @click="handleNewNote"
        class="group flex items-center gap-2.5 px-5 py-3 bg-primary text-text-inverse rounded-xl text-sm font-semibold hover:bg-primary-hover active:scale-[0.98] transition-all duration-200 shadow-[0_2px_16px_rgba(184,98,27,0.2)]"
      >
        <svg class="w-4 h-4 transition-transform group-hover:rotate-90 duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Новая заметка
      </button>
    </div>

    <!-- Pinned -->
    <section v-if="pinnedNotes.length" class="mb-12 anim-fade-up d3">
      <div class="flex items-center gap-3 mb-4">
        <h2 class="text-[10px] font-semibold text-text-tertiary uppercase tracking-[0.15em]">Закреплённые</h2>
        <div class="flex-1 h-px bg-border-light"></div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="(note, i) in pinnedNotes"
          :key="note.id"
          @click="openNote(note)"
          :style="{ animationDelay: `${i * 60 + 200}ms` }"
          class="group p-5 bg-bg-surface border border-border rounded-2xl cursor-pointer hover:border-primary/30 hover:shadow-[0_4px_20px_rgba(45,32,24,0.08)] transition-all duration-250 anim-fade-up"
        >
          <div class="flex items-start justify-between mb-2">
            <h3 class="font-semibold text-text-primary text-[15px] truncate flex-1" style="font-family: var(--font-heading); font-size: 18px;">{{ note.title }}</h3>
            <span class="text-primary/40 ml-2 text-xs">&#9733;</span>
          </div>
          <p class="text-sm text-text-secondary leading-relaxed line-clamp-2">{{ note.content_md.slice(0, 120) }}</p>
          <p class="text-[11px] text-text-tertiary mt-3">{{ formatDate(note.updated_at) }}</p>
        </div>
      </div>
    </section>

    <!-- Recent -->
    <section class="anim-fade-up d4">
      <div class="flex items-center gap-3 mb-4">
        <h2 class="text-[10px] font-semibold text-text-tertiary uppercase tracking-[0.15em]">Недавние</h2>
        <div class="flex-1 h-px bg-border-light"></div>
      </div>

      <div v-if="recentNotes.length" class="space-y-3">
        <div
          v-for="(note, i) in recentNotes"
          :key="note.id"
          @click="openNote(note)"
          :style="{ animationDelay: `${i * 40 + 300}ms` }"
          class="group flex items-start gap-4 p-4 bg-bg-surface/60 border border-transparent rounded-2xl cursor-pointer hover:bg-bg-surface hover:border-border hover:shadow-[0_2px_12px_rgba(45,32,24,0.05)] transition-all duration-200 anim-fade-up"
        >
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-text-primary text-[15px] mb-0.5 group-hover:text-primary transition-colors" style="font-family: var(--font-heading); font-size: 17px;">
              {{ note.title }}
            </h3>
            <p class="text-sm text-text-secondary/80 line-clamp-1">{{ note.content_md.slice(0, 200) }}</p>
          </div>
          <span class="text-[11px] text-text-tertiary whitespace-nowrap shrink-0 mt-0.5">{{ formatDate(note.updated_at) }}</span>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="text-center py-16">
        <div class="text-5xl mb-4 opacity-20">&#9998;</div>
        <p class="text-text-tertiary text-sm mb-1">Заметок пока нет</p>
        <p class="text-text-tertiary/60 text-xs">Нажмите "Новая заметка" чтобы начать</p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useWorkspaceStore } from '@/stores/workspace'
import { useNoteStore } from '@/stores/note'
import type { Note } from '@/types/note'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const workspaceStore = useWorkspaceStore()
const noteStore = useNoteStore()

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return 'Доброй ночи'
  if (hour < 12) return 'Доброе утро'
  if (hour < 18) return 'Добрый день'
  return 'Добрый вечер'
})

const formattedDate = computed(() =>
  new Date().toLocaleDateString('ru-RU', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
)

const pinnedNotes = computed(() => noteStore.notes.filter((n) => n.is_pinned))
const recentNotes = computed(() => noteStore.notes.filter((n) => !n.is_pinned).slice(0, 10))

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function openNote(note: Note) {
  router.push(`/w/${route.params.wsSlug}/note/${note.id}`)
}

async function handleNewNote() {
  if (!workspaceStore.currentWorkspace) return
  const note = await noteStore.createNote(workspaceStore.currentWorkspace.id, 'Без названия')
  router.push(`/w/${route.params.wsSlug}/note/${note.id}`)
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  await workspaceStore.setCurrentBySlug(slug)
  if (workspaceStore.currentWorkspace) {
    // Note store handles decryption transparently in fetchNotes
    await noteStore.fetchNotes(workspaceStore.currentWorkspace.id)
  }
})
</script>
