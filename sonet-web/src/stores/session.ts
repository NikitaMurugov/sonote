import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/composables/useApi'

export interface Session {
  id: number
  user_id: number
  device_name: string | null
  device_type: string
  os: string | null
  browser: string | null
  ip_address: string | null
  location: string | null
  last_active_at: string
  created_at: string
  is_current: boolean
}

export const useSessionStore = defineStore('session', () => {
  const sessions = ref<Session[]>([])
  const loading = ref(false)

  async function fetchSessions() {
    loading.value = true
    try {
      const { data } = await api.get('/sessions')
      sessions.value = data.data || []
    } catch {
      sessions.value = []
    } finally {
      loading.value = false
    }
  }

  async function revokeSession(sessionId: number) {
    await api.delete(`/sessions/${sessionId}`)
    sessions.value = sessions.value.filter((s) => s.id !== sessionId)
  }

  async function revokeAllOther(currentSessionId: number) {
    await api.post('/sessions/revoke-all', { current_session_id: currentSessionId })
    sessions.value = sessions.value.filter((s) => s.id === currentSessionId)
  }

  return {
    sessions,
    loading,
    fetchSessions,
    revokeSession,
    revokeAllOther,
  }
})
