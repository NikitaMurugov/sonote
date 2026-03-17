<template>
  <header class="h-11 border-b border-border-light/60 flex items-center justify-between px-4 bg-bg-base/80 backdrop-blur-sm shrink-0">
    <div class="flex items-center gap-3">
      <button
        @click="uiStore.toggleSidebar"
        class="text-text-tertiary hover:text-text-primary transition-colors duration-150 p-1 -ml-1 rounded-md hover:bg-bg-hover"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>

      <!-- Breadcrumb -->
      <div class="flex items-center gap-1.5 text-xs text-text-tertiary">
        <span>{{ workspaceStore.currentWorkspace?.name }}</span>
        <template v-if="currentFolderName">
          <span class="opacity-40">/</span>
          <span>{{ currentFolderName }}</span>
        </template>
        <template v-if="noteStore.currentNote">
          <span class="opacity-40">/</span>
          <span class="text-text-secondary">{{ noteStore.currentNote.title || 'Без названия' }}</span>
        </template>
      </div>
    </div>

    <div class="flex items-center gap-1">
      <!-- Search -->
      <router-link
        :to="`/w/${$route.params.wsSlug}/search`"
        class="p-2 rounded-lg text-text-tertiary hover:text-text-primary hover:bg-bg-hover transition-all duration-150"
        title="Поиск"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </router-link>

      <!-- Theme -->
      <button
        @click="toggleTheme"
        class="p-2 rounded-lg text-text-tertiary hover:text-text-primary hover:bg-bg-hover transition-all duration-150"
        title="Сменить тему"
      >
        <svg v-if="theme === 'light'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
      </button>

      <!-- Avatar -->
      <router-link
        to="/profile"
        class="ml-1 w-7 h-7 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold hover:bg-primary/20 transition-colors duration-150"
        title="Профиль"
      >
        {{ authStore.user?.display_name?.charAt(0)?.toUpperCase() || '?' }}
      </router-link>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUiStore } from '@/stores/ui'
import { useWorkspaceStore } from '@/stores/workspace'
import { useAuthStore } from '@/stores/auth'
import { useNoteStore } from '@/stores/note'
import { useFolderStore } from '@/stores/folder'
import { useTheme } from '@/composables/useTheme'

const uiStore = useUiStore()
const workspaceStore = useWorkspaceStore()
const authStore = useAuthStore()
const noteStore = useNoteStore()
const folderStore = useFolderStore()
const { theme, toggleTheme } = useTheme()

const currentFolderName = computed(() => {
  if (folderStore.selectedFolderId === null) return null
  const folder = folderStore.folders.find((f) => f.id === folderStore.selectedFolderId)
  return folder?.name || null
})
</script>
