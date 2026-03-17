<template>
  <div class="px-10 py-10 max-w-xl mx-auto anim-fade-up">
    <h1 class="text-2xl font-semibold mb-10" style="font-family: var(--font-heading)">Профиль</h1>

    <div v-if="authStore.user" class="flex flex-col gap-6">
      <!-- Avatar -->
      <div class="flex items-center gap-5 pb-2">
        <div class="w-16 h-16 rounded-full bg-primary/10 text-primary flex items-center justify-center text-2xl font-bold shrink-0" style="font-family: var(--font-heading)">
          {{ authStore.user.display_name?.charAt(0)?.toUpperCase() || '?' }}
        </div>
        <div>
          <p class="font-semibold text-text-primary text-lg">{{ authStore.user.display_name }}</p>
          <p class="text-sm text-text-secondary">@{{ authStore.user.username }}</p>
          <p class="text-xs text-text-tertiary mt-0.5">{{ authStore.user.email }}</p>
        </div>
      </div>

      <div class="h-px bg-gradient-to-r from-border via-border-light to-transparent"></div>

      <!-- Encryption -->
      <div class="p-5 bg-bg-surface rounded-2xl border border-border space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-text-primary">E2E шифрование</p>
            <p class="text-xs text-text-tertiary">
              {{ encryptionStore.isUnlocked ? 'Активно' : 'Заблокировано' }}
            </p>
          </div>
          <span
            :class="[
              'w-2.5 h-2.5 rounded-full',
              encryptionStore.isUnlocked ? 'bg-success' : 'bg-warning',
            ]"
          ></span>
        </div>

        <!-- Unlock -->
        <div v-if="!encryptionStore.isUnlocked">
          <div class="flex gap-2">
            <input
              v-model="unlockPassword"
              type="password"
              placeholder="Пароль аккаунта"
              class="flex-1 px-3 py-2 rounded-lg border border-border bg-bg-base text-text-primary text-sm focus:outline-none focus:border-primary transition"
              @keydown.enter="handleUnlock"
            />
            <button
              @click="handleUnlock"
              :disabled="encryptionStore.loading"
              class="px-4 py-2 rounded-lg bg-primary text-text-inverse text-sm font-medium hover:bg-primary-hover transition disabled:opacity-50"
            >
              Разблокировать
            </button>
          </div>
          <p v-if="unlockError" class="text-error text-xs mt-1">{{ unlockError }}</p>
        </div>

        <!-- Unlocked -->
        <div v-else>
          <button
            @click="encryptionStore.lock()"
            class="w-full py-2 rounded-lg border border-border text-text-secondary text-sm hover:bg-bg-hover transition"
          >
            Заблокировать
          </button>
        </div>
      </div>

      <!-- Sessions -->
      <div class="p-5 bg-bg-surface rounded-2xl border border-border space-y-4">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-text-primary">Активные сессии</p>
          <button
            v-if="sessionStore.sessions.length > 1"
            @click="handleRevokeAll"
            class="text-xs text-error hover:underline"
          >
            Завершить все
          </button>
        </div>

        <div v-if="sessionStore.loading" class="text-xs text-text-tertiary text-center py-2">
          Загрузка...
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="session in sessionStore.sessions"
            :key="session.id"
            class="flex items-center gap-3 p-2 rounded-lg"
            :class="session.is_current ? 'bg-primary/5' : ''"
          >
            <div class="flex-1 min-w-0">
              <p class="text-sm text-text-primary truncate">
                {{ session.browser || 'Unknown' }} / {{ session.os || 'Unknown' }}
                <span v-if="session.is_current" class="text-xs text-primary ml-1">(текущая)</span>
              </p>
              <p class="text-[11px] text-text-tertiary">
                {{ session.ip_address || '-' }} &middot; {{ formatDate(session.last_active_at) }}
              </p>
            </div>
            <button
              v-if="!session.is_current"
              @click="sessionStore.revokeSession(session.id)"
              class="text-xs text-text-tertiary hover:text-error transition-colors shrink-0"
            >
              Завершить
            </button>
          </div>
        </div>
      </div>

      <!-- Theme -->
      <div class="p-5 bg-bg-surface rounded-2xl border border-border">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-text-primary">Тема оформления</p>
            <p class="text-xs text-text-tertiary">{{ theme === 'light' ? 'Светлая (Parchment)' : 'Тёмная (Inkwell)' }}</p>
          </div>
          <button
            @click="toggleTheme"
            class="px-4 py-2 rounded-lg bg-bg-hover text-sm text-text-secondary hover:text-text-primary transition-colors"
          >
            {{ theme === 'light' ? 'Тёмная' : 'Светлая' }}
          </button>
        </div>
      </div>

      <div class="h-px bg-gradient-to-r from-border via-border-light to-transparent"></div>

      <!-- Logout -->
      <button
        @click="handleLogout"
        class="w-full py-3 rounded-2xl border border-error/30 text-error text-sm font-medium hover:bg-error/5 transition-colors"
      >
        Выйти
      </button>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useEncryptionStore } from '@/stores/encryption'
import { useSessionStore } from '@/stores/session'
import { useTheme } from '@/composables/useTheme'

const router = useRouter()
const authStore = useAuthStore()
const encryptionStore = useEncryptionStore()
const sessionStore = useSessionStore()
const { theme, toggleTheme } = useTheme()

const unlockPassword = ref('')
const unlockError = ref('')

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function handleUnlock() {
  unlockError.value = ''
  try {
    await encryptionStore.unlock(unlockPassword.value)
    unlockPassword.value = ''
  } catch {
    unlockError.value = 'Неверный пароль'
  }
}

async function handleRevokeAll() {
  const current = sessionStore.sessions.find((s) => s.is_current)
  if (current) {
    await sessionStore.revokeAllOther(current.id)
  }
}

function handleLogout() {
  encryptionStore.lock()
  authStore.logout()
  router.push('/login')
}

onMounted(async () => {
  await encryptionStore.checkSetup()
  await encryptionStore.ensureUnlocked()
  await sessionStore.fetchSessions()
})
</script>
