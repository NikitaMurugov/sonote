<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base relative overflow-hidden">
    <!-- Atmospheric blurs -->
    <div class="absolute top-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full bg-primary-light opacity-30 blur-3xl"></div>
    <div class="absolute bottom-[-15%] left-[-5%] w-[400px] h-[400px] rounded-full bg-accent-soft opacity-20 blur-3xl"></div>
    <div class="absolute top-0 left-1/2 -translate-x-1/2 w-px h-32 bg-gradient-to-b from-transparent via-border to-transparent"></div>

    <div class="w-full max-w-[420px] px-6 relative z-10">
      <!-- Logo -->
      <div class="text-center mb-10 anim-fade-up">
        <h1 class="text-5xl font-semibold tracking-tight text-text-primary" style="font-family: var(--font-heading)">
          Sonote
        </h1>
        <div class="mt-3 flex items-center justify-center gap-3">
          <span class="h-px w-8 bg-border"></span>
          <p class="text-text-tertiary text-xs tracking-[0.2em] uppercase">ваши мысли, связанные воедино</p>
          <span class="h-px w-8 bg-border"></span>
        </div>
      </div>

      <!-- Card -->
      <div class="bg-bg-surface/80 backdrop-blur-sm rounded-2xl border border-border p-8 shadow-[0_4px_24px_rgba(45,32,24,0.06)] anim-fade-up d2">
        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Email</label>
            <input
              v-model="email"
              type="email"
              required
              autocomplete="email"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Пароль</label>
            <input
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="••••••••"
            />
          </div>

          <p v-if="error" class="text-error text-sm bg-error/8 px-3 py-2 rounded-lg">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-primary text-text-inverse rounded-xl text-sm font-semibold tracking-wide hover:bg-primary-hover active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-[0_2px_12px_rgba(184,98,27,0.25)]"
          >
            {{ loading ? 'Входим...' : 'Войти' }}
          </button>
        </form>
      </div>

      <p class="text-center text-sm text-text-tertiary mt-6 anim-fade-up d4">
        Нет аккаунта?
        <router-link to="/register" class="text-primary hover:text-primary-hover font-medium transition-colors underline underline-offset-4 decoration-primary/30 hover:decoration-primary">
          Создать
        </router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const router = useRouter()
const authStore = useAuthStore()

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await authStore.loginAndSetup(email.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Неверный email или пароль'
  } finally {
    loading.value = false
  }
}
</script>
