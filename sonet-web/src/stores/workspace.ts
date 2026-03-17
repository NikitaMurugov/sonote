import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/composables/useApi'
import type { Workspace } from '@/types/workspace'

const LAST_WS_KEY = 'sonote-last-ws'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces = ref<Workspace[]>([])
  const currentWorkspace = ref<Workspace | null>(null)

  const currentSlug = computed(() =>
    currentWorkspace.value?.slug || localStorage.getItem(LAST_WS_KEY) || '',
  )

  async function fetchWorkspaces() {
    const { data } = await api.get('/workspaces')
    workspaces.value = data.data || []
  }

  async function setCurrentBySlug(slug: string) {
    if (!workspaces.value.length) {
      await fetchWorkspaces()
    }

    // If no slug provided, restore from localStorage
    const effectiveSlug = slug || localStorage.getItem(LAST_WS_KEY) || ''

    currentWorkspace.value = workspaces.value.find((w) => w.slug === effectiveSlug) || null

    // Fallback to first workspace if slug not found
    if (!currentWorkspace.value && workspaces.value.length) {
      currentWorkspace.value = workspaces.value[0]
    }

    // Persist for cross-page navigation
    if (currentWorkspace.value) {
      localStorage.setItem(LAST_WS_KEY, currentWorkspace.value.slug)
    }
  }

  async function createWorkspace(name: string, description?: string) {
    const { data } = await api.post('/workspaces', { name, description })
    workspaces.value.push(data.data)
    return data.data as Workspace
  }

  return {
    workspaces,
    currentWorkspace,
    currentSlug,
    fetchWorkspaces,
    setCurrentBySlug,
    createWorkspace,
  }
})
