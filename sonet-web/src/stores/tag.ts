import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'
import type { Tag } from '@/types/tag'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])

  async function fetchTags(workspaceId: number) {
    const { data } = await api.get(`/workspaces/${workspaceId}/tags`)
    tags.value = data.data || []
  }

  async function createTag(workspaceId: number, name: string, color?: string) {
    const { data } = await api.post(`/workspaces/${workspaceId}/tags`, { name, color })
    tags.value.push(data.data)
    return data.data as Tag
  }

  function clear() {
    tags.value = []
  }

  return { tags, fetchTags, createTag, clear }
})
