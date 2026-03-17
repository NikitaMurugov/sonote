import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'
import type { Folder } from '@/types/folder'

export const useFolderStore = defineStore('folder', () => {
  const folders = ref<Folder[]>([])
  const selectedFolderId = ref<number | null>(null)

  async function fetchFolders(workspaceId: number) {
    const { data } = await api.get(`/workspaces/${workspaceId}/folders`)
    folders.value = data.data || []
  }

  async function createFolder(workspaceId: number, name: string, parentId?: number) {
    const { data } = await api.post(`/workspaces/${workspaceId}/folders`, {
      name,
      parent_id: parentId || null,
    })
    folders.value.push(data.data)
    return data.data as Folder
  }

  async function deleteFolder(workspaceId: number, folderId: number) {
    await api.delete(`/workspaces/${workspaceId}/folders/${folderId}`)
    folders.value = folders.value.filter((f) => f.id !== folderId)
  }

  function selectFolder(folderId: number | null) {
    selectedFolderId.value = folderId
  }

  function clear() {
    folders.value = []
    selectedFolderId.value = null
  }

  return {
    folders,
    selectedFolderId,
    fetchFolders,
    createFolder,
    deleteFolder,
    selectFolder,
    clear,
  }
})
