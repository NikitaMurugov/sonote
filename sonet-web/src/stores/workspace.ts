import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'
import type { Workspace } from '@/types/workspace'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces = ref<Workspace[]>([])
  const currentWorkspace = ref<Workspace | null>(null)

  async function fetchWorkspaces() {
    const { data } = await api.get('/workspaces')
    workspaces.value = data.data || []
  }

  async function setCurrentBySlug(slug: string) {
    if (!workspaces.value.length) {
      await fetchWorkspaces()
    }
    currentWorkspace.value = workspaces.value.find((w) => w.slug === slug) || null
  }

  async function createWorkspace(name: string, description?: string) {
    const { data } = await api.post('/workspaces', { name, description })
    workspaces.value.push(data.data)
    return data.data as Workspace
  }

  return {
    workspaces,
    currentWorkspace,
    fetchWorkspaces,
    setCurrentBySlug,
    createWorkspace,
  }
})
