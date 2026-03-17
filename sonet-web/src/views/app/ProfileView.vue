<template>
  <div class="p-8 max-w-lg mx-auto anim-fade-up">
    <h1 class="text-2xl font-semibold mb-8" style="font-family: var(--font-heading)">Профиль</h1>

    <div v-if="authStore.user" class="space-y-6">
      <!-- Avatar -->
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-primary/10 text-primary flex items-center justify-center text-2xl font-bold" style="font-family: var(--font-heading)">
          {{ authStore.user.display_name?.charAt(0)?.toUpperCase() || '?' }}
        </div>
        <div>
          <p class="font-semibold text-text-primary">{{ authStore.user.display_name }}</p>
          <p class="text-sm text-text-secondary">@{{ authStore.user.username }}</p>
          <p class="text-xs text-text-tertiary">{{ authStore.user.email }}</p>
        </div>
      </div>

      <!-- Theme -->
      <div class="p-4 bg-bg-surface rounded-xl border border-border">
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

      <!-- Logout -->
      <button
        @click="handleLogout"
        class="w-full py-2.5 rounded-xl border border-error/30 text-error text-sm font-medium hover:bg-error/5 transition-colors"
      >
        Выйти
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'

const router = useRouter()
const authStore = useAuthStore()
const { theme, toggleTheme } = useTheme()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
