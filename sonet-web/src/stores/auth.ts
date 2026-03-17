import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/composables/useApi'
import type { User } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('sonote-access-token'))
  const refreshToken = ref<string | null>(localStorage.getItem('sonote-refresh-token'))

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(email: string, password: string) {
    const { data } = await api.post('/auth/login', { email, password })
    setAuth(data.data)
  }

  async function register(email: string, username: string, password: string, displayName: string) {
    const { data } = await api.post('/auth/register', {
      email,
      username,
      password,
      display_name: displayName,
    })
    setAuth(data.data)

    // Auto-setup encryption using the same password
    const { useEncryptionStore } = await import('@/stores/encryption')
    const encStore = useEncryptionStore()
    const recoveryKey = await encStore.setup(password)
    return recoveryKey
  }

  async function refreshTokens() {
    if (!refreshToken.value) throw new Error('No refresh token')
    const { data } = await api.post('/auth/refresh', {
      refresh_token: refreshToken.value,
    })
    accessToken.value = data.data.access_token
    refreshToken.value = data.data.refresh_token
    localStorage.setItem('sonote-access-token', data.data.access_token)
    localStorage.setItem('sonote-refresh-token', data.data.refresh_token)
  }

  async function fetchUser() {
    const { data } = await api.get('/users/me')
    user.value = data.data
  }

  function setAuth(authData: { user: User; tokens: { access_token: string; refresh_token: string } }) {
    user.value = authData.user
    accessToken.value = authData.tokens.access_token
    refreshToken.value = authData.tokens.refresh_token
    localStorage.setItem('sonote-access-token', authData.tokens.access_token)
    localStorage.setItem('sonote-refresh-token', authData.tokens.refresh_token)
  }

  function logout() {
    if (refreshToken.value) {
      api.post('/auth/logout', { refresh_token: refreshToken.value }).catch(() => {})
    }
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('sonote-access-token')
    localStorage.removeItem('sonote-refresh-token')
  }

  async function loginAndSetup(email: string, password: string) {
    await login(email, password)
    // Try to auto-unlock encryption if set up
    try {
      const { useEncryptionStore } = await import('@/stores/encryption')
      const encStore = useEncryptionStore()
      await encStore.checkSetup()
      if (encStore.isSetup) {
        await encStore.unlock(password)
      }
    } catch (e) {
      console.warn('Auto-unlock encryption failed:', e)
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    login,
    loginAndSetup,
    register,
    refreshTokens,
    fetchUser,
    logout,
  }
})
