<template>
  <aside class="h-screen w-60 bg-bg-sidebar border-r border-border flex flex-col shrink-0 anim-slide-r">
    <!-- Workspace Picker -->
    <div class="px-4 py-4 border-b border-border-light">
      <button
        @click="showWsPicker = !showWsPicker"
        class="w-full flex items-center gap-2.5 group"
      >
        <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary text-sm font-bold shrink-0" style="font-family: var(--font-heading)">
          {{ workspaceStore.currentWorkspace?.name?.charAt(0)?.toUpperCase() || 'S' }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-text-primary truncate">
            {{ workspaceStore.currentWorkspace?.name || 'Sonote' }}
          </p>
          <p class="text-[10px] text-text-tertiary tracking-wide uppercase">workspace</p>
        </div>
        <svg
          class="w-3.5 h-3.5 text-text-tertiary transition-transform duration-200 shrink-0"
          :class="showWsPicker ? 'rotate-180' : ''"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      <!-- Dropdown -->
      <div v-if="showWsPicker" class="mt-2 space-y-0.5">
        <button
          v-for="ws in workspaceStore.workspaces"
          :key="ws.id"
          @click="switchWorkspace(ws)"
          :class="[
            'w-full text-left px-3 py-2 rounded-lg text-sm transition-all duration-150 truncate',
            ws.id === workspaceStore.currentWorkspace?.id
              ? 'bg-bg-active text-text-primary font-medium'
              : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
          ]"
        >
          {{ ws.name }}
        </button>
      </div>
    </div>

    <!-- New Note Button -->
    <div class="px-3 pt-3">
      <button
        @click="handleNewNote"
        class="group w-full flex items-center gap-2 py-2.5 px-3 bg-primary text-text-inverse rounded-xl text-sm font-medium hover:bg-primary-hover active:scale-[0.98] transition-all duration-200 shadow-[0_1px_8px_rgba(184,98,27,0.2)]"
      >
        <svg class="w-4 h-4 transition-transform group-hover:rotate-90 duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Новая заметка
      </button>
    </div>

    <!-- Folders -->
    <div class="flex-1 overflow-auto px-2 pt-4">
      <div class="flex items-center justify-between px-3 mb-2">
        <p class="text-[10px] font-semibold text-text-tertiary uppercase tracking-[0.15em]">Папки</p>
        <button
          @click="showNewFolder = !showNewFolder"
          class="text-text-tertiary hover:text-primary transition-colors p-0.5 rounded"
          title="Новая папка"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      <!-- New folder input -->
      <div v-if="showNewFolder" class="px-2 mb-2">
        <input
          ref="newFolderInput"
          v-model="newFolderName"
          @keydown.enter="createFolder"
          @keydown.escape="showNewFolder = false"
          placeholder="Название папки"
          class="w-full px-3 py-1.5 text-sm rounded-lg border border-border bg-bg-base text-text-primary placeholder-text-tertiary focus:outline-none focus:border-primary transition"
        />
      </div>

      <button
        @click="selectFolder(null)"
        :class="[
          'w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm cursor-pointer transition-all duration-150',
          folderStore.selectedFolderId === null
            ? 'bg-bg-active text-text-primary font-medium'
            : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
        ]"
      >
        <svg class="w-3.5 h-3.5 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        <span class="flex-1 text-left">Все заметки</span>
        <span class="text-[10px] text-text-tertiary">{{ noteStore.notes.length }}</span>
      </button>

      <button
        v-for="folder in rootFolders"
        :key="folder.id"
        @click="selectFolder(folder.id)"
        :class="[
          'w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm cursor-pointer transition-all duration-150 group',
          folderStore.selectedFolderId === folder.id
            ? 'bg-bg-active text-text-primary font-medium shadow-[inset_3px_0_0_var(--color-primary)]'
            : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
        ]"
      >
        <svg class="w-3.5 h-3.5 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        <span class="flex-1 text-left truncate">{{ folder.name }}</span>
        <button
          @click.stop="deleteFolder(folder.id)"
          class="opacity-0 group-hover:opacity-100 text-text-tertiary hover:text-error transition-all p-0.5 rounded"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </button>

      <!-- Note list in sidebar -->
      <div class="mt-4 border-t border-border-light pt-3">
        <p class="text-[10px] font-semibold text-text-tertiary uppercase tracking-[0.15em] px-3 mb-2">Заметки</p>
        <div v-if="noteStore.notes.length" class="space-y-0.5">
          <button
            v-for="note in noteStore.notes.slice(0, 20)"
            :key="note.id"
            @click="openNote(note.id)"
            :class="[
              'w-full text-left px-3 py-1.5 rounded-lg text-sm truncate transition-all duration-150',
              noteStore.currentNote?.id === note.id
                ? 'bg-bg-active text-text-primary font-medium'
                : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
            ]"
          >
            {{ note.title || 'Без названия' }}
          </button>
        </div>
        <p v-else class="px-3 text-xs text-text-tertiary">Пусто</p>
      </div>
    </div>

    <!-- Bottom -->
    <div class="px-2 pb-3 pt-2 border-t border-border-light space-y-0.5">
      <router-link
        :to="`/w/${workspaceStore.currentSlug}/graph`"
        class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-text-secondary hover:bg-bg-hover hover:text-text-primary transition-all duration-150"
      >
        <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        Граф
      </router-link>
      <router-link
        :to="`/w/${workspaceStore.currentSlug}/search`"
        class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-text-secondary hover:bg-bg-hover hover:text-text-primary transition-all duration-150"
      >
        <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        Поиск
      </router-link>
      <router-link
        :to="`/w/${workspaceStore.currentSlug}/settings`"
        class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-text-secondary hover:bg-bg-hover hover:text-text-primary transition-all duration-150"
      >
        <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        Настройки
      </router-link>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUiStore } from '@/stores/ui'
import { useWorkspaceStore } from '@/stores/workspace'
import { useFolderStore } from '@/stores/folder'
import { useNoteStore } from '@/stores/note'
import type { Workspace } from '@/types/workspace'

const route = useRoute()
const router = useRouter()
const uiStore = useUiStore()
const workspaceStore = useWorkspaceStore()
const folderStore = useFolderStore()
const noteStore = useNoteStore()

const showWsPicker = ref(false)
const showNewFolder = ref(false)
const newFolderName = ref('')
const newFolderInput = ref<HTMLInputElement>()

const rootFolders = computed(() => folderStore.folders.filter((f) => !f.parent_id))

watch(showNewFolder, async (v) => {
  if (v) {
    await nextTick()
    newFolderInput.value?.focus()
  }
})

function switchWorkspace(ws: Workspace) {
  showWsPicker.value = false
  router.push(`/w/${ws.slug}`)
}

async function createFolder() {
  if (!newFolderName.value.trim() || !workspaceStore.currentWorkspace) return
  await folderStore.createFolder(workspaceStore.currentWorkspace.id, newFolderName.value.trim())
  newFolderName.value = ''
  showNewFolder.value = false
}

async function deleteFolder(folderId: number) {
  if (!workspaceStore.currentWorkspace) return
  await folderStore.deleteFolder(workspaceStore.currentWorkspace.id, folderId)
  if (folderStore.selectedFolderId === folderId) {
    selectFolder(null)
  }
}

function selectFolder(folderId: number | null) {
  folderStore.selectFolder(folderId)
  if (workspaceStore.currentWorkspace) {
    noteStore.fetchNotes(workspaceStore.currentWorkspace.id, folderId ?? undefined)
  }
}

function openNote(noteId: number) {
  router.push(`/w/${workspaceStore.currentSlug}/note/${noteId}`)
}

async function handleNewNote() {
  if (!workspaceStore.currentWorkspace) return
  const note = await noteStore.createNote(
    workspaceStore.currentWorkspace.id,
    'Без названия',
    folderStore.selectedFolderId ?? undefined
  )
  router.push(`/w/${workspaceStore.currentSlug}/note/${note.id}`)
}

watch(
  () => route.params.wsSlug,
  async (slug) => {
    if (slug) {
      await workspaceStore.setCurrentBySlug(slug as string)
    } else if (!workspaceStore.currentWorkspace) {
      // No slug in URL (e.g. /profile) — restore from localStorage
      await workspaceStore.setCurrentBySlug('')
    }
    if (workspaceStore.currentWorkspace) {
      await folderStore.fetchFolders(workspaceStore.currentWorkspace.id)
      await noteStore.fetchNotes(workspaceStore.currentWorkspace.id)
    }
  },
  { immediate: true }
)
</script>
