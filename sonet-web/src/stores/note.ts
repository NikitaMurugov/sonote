import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'
import type { Note } from '@/types/note'

export const useNoteStore = defineStore('note', () => {
  const notes = ref<Note[]>([])
  const currentNote = ref<Note | null>(null)

  async function fetchNotes(workspaceId: number, folderId?: number) {
    const params: Record<string, string> = {}
    if (folderId) params.folder_id = String(folderId)
    const { data } = await api.get(`/workspaces/${workspaceId}/notes`, { params })
    notes.value = data.data || []
  }

  async function fetchNote(workspaceId: number, noteId: number) {
    const { data } = await api.get(`/workspaces/${workspaceId}/notes/${noteId}`)
    currentNote.value = data.data
  }

  async function createNote(workspaceId: number, title: string, folderId?: number) {
    const { data } = await api.post(`/workspaces/${workspaceId}/notes`, {
      title,
      folder_id: folderId || null,
      content_md: '',
      content_html: '',
    })
    const note = data.data as Note
    notes.value.unshift(note)
    return note
  }

  async function updateNote(workspaceId: number, noteId: number, updates: Partial<Note>) {
    const { data } = await api.patch(`/workspaces/${workspaceId}/notes/${noteId}`, updates)
    currentNote.value = data.data
    const idx = notes.value.findIndex((n) => n.id === noteId)
    if (idx !== -1) notes.value[idx] = data.data
  }

  async function deleteNote(workspaceId: number, noteId: number) {
    await api.delete(`/workspaces/${workspaceId}/notes/${noteId}`)
    notes.value = notes.value.filter((n) => n.id !== noteId)
    if (currentNote.value?.id === noteId) currentNote.value = null
  }

  function clear() {
    notes.value = []
    currentNote.value = null
  }

  return {
    notes,
    currentNote,
    fetchNotes,
    fetchNote,
    createNote,
    updateNote,
    deleteNote,
    clear,
  }
})
