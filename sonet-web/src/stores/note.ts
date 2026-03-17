import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'
import { useEncryptionStore } from '@/stores/encryption'
import type { Note } from '@/types/note'

export const useNoteStore = defineStore('note', () => {
  const notes = ref<Note[]>([])
  const currentNote = ref<Note | null>(null)

  async function getEncryption() {
    const enc = useEncryptionStore()
    if (enc.isUnlocked && enc.getPrivateKey()) return enc
    // Try auto-restore from sessionStorage
    const restored = await enc.ensureUnlocked()
    return restored ? enc : null
  }

  async function decryptNoteInPlace(workspaceId: number, note: Note): Promise<Note> {
    if (!note.is_encrypted) return note
    const enc = await getEncryption()
    if (!enc) return note

    try {
      const decrypted = await enc.decryptNote(workspaceId, note)
      if (decrypted) {
        return {
          ...note,
          title: decrypted.title,
          content_md: decrypted.content_md,
          content_html: decrypted.content_html,
          content_json: decrypted.content_json,
        }
      }
    } catch { /* return as-is */ }
    return note
  }

  async function fetchNotes(workspaceId: number, folderId?: number) {
    const params: Record<string, string> = {}
    if (folderId) params.folder_id = String(folderId)
    const { data } = await api.get(`/workspaces/${workspaceId}/notes`, { params })
    const raw: Note[] = data.data || []

    // Decrypt titles for list view
    const enc = await getEncryption()
    if (enc) {
      const decrypted = await Promise.all(
        raw.map(async (n) => {
          if (!n.is_encrypted) return n
          try {
            const d = await enc.decryptNote(workspaceId, n)
            if (d) return { ...n, title: d.title, content_md: d.content_md }
          } catch { /* skip */ }
          return n
        }),
      )
      notes.value = decrypted
    } else {
      notes.value = raw
    }
  }

  async function fetchNote(workspaceId: number, noteId: number) {
    const { data } = await api.get(`/workspaces/${workspaceId}/notes/${noteId}`)
    currentNote.value = await decryptNoteInPlace(workspaceId, data.data)
  }

  async function createNote(workspaceId: number, title: string, folderId?: number) {
    const enc = await getEncryption()

    let payload: Record<string, any> = {
      title,
      folder_id: folderId || null,
      content_md: '',
      content_html: '',
    }

    // Encrypt if possible
    if (enc) {
      try {
        const encrypted = await enc.encryptNote(workspaceId, title, '', '', null)
        payload = { ...payload, ...encrypted }
      } catch { /* fallback to plaintext */ }
    }

    const { data } = await api.post(`/workspaces/${workspaceId}/notes`, payload)
    const note = await decryptNoteInPlace(workspaceId, data.data as Note)
    notes.value.unshift(note)
    return note
  }

  async function updateNote(workspaceId: number, noteId: number, updates: Partial<Note>) {
    const enc = await getEncryption()

    let payload: Record<string, any> = { ...updates }

    // If we have content updates and encryption is available, encrypt
    if (enc && (updates.content_md !== undefined || updates.title !== undefined)) {
      try {
        // Get current note data to merge with updates
        const current = currentNote.value
        const title = updates.title ?? current?.title ?? ''
        const contentMd = updates.content_md ?? current?.content_md ?? ''
        const contentHtml = updates.content_html ?? current?.content_html ?? ''
        const contentJson = updates.content_json ?? current?.content_json ?? null

        const encrypted = await enc.encryptNote(workspaceId, title, contentMd, contentHtml, contentJson)
        payload = {
          ...payload,
          ...encrypted,
        }
      } catch { /* fallback to plaintext */ }
    }

    const { data } = await api.patch(`/workspaces/${workspaceId}/notes/${noteId}`, payload)
    const note = await decryptNoteInPlace(workspaceId, data.data)
    currentNote.value = note
    const idx = notes.value.findIndex((n) => n.id === noteId)
    if (idx !== -1) notes.value[idx] = note
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
