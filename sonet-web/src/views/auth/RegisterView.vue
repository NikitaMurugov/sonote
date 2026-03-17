<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base relative overflow-hidden">
    <div class="absolute top-[-15%] left-[-10%] w-[500px] h-[500px] rounded-full bg-accent-soft opacity-20 blur-3xl"></div>
    <div class="absolute bottom-[-20%] right-[-8%] w-[600px] h-[600px] rounded-full bg-primary-light opacity-25 blur-3xl"></div>

    <div class="w-full max-w-[420px] px-6 relative z-10">
      <div class="text-center mb-10 anim-fade-up">
        <h1 class="text-5xl font-semibold tracking-tight text-text-primary" style="font-family: var(--font-heading)">
          Sonote
        </h1>
        <div class="mt-3 flex items-center justify-center gap-3">
          <span class="h-px w-8 bg-border"></span>
          <p class="text-text-tertiary text-xs tracking-[0.2em] uppercase">создайте аккаунт</p>
          <span class="h-px w-8 bg-border"></span>
        </div>
      </div>

      <div class="bg-bg-surface/80 backdrop-blur-sm rounded-2xl border border-border p-8 shadow-[0_4px_24px_rgba(45,32,24,0.06)] anim-fade-up d2">
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Email</label>
            <input v-model="email" type="email" required autocomplete="email"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="you@example.com" />
          </div>
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Имя пользователя</label>
            <input v-model="username" type="text" required autocomplete="username"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="username" />
          </div>
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Отображаемое имя</label>
            <input v-model="displayName" type="text" autocomplete="name"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="Как вас называть" />
          </div>
          <div>
            <label class="block text-xs font-medium text-text-secondary tracking-wide uppercase mb-2">Пароль</label>
            <input v-model="password" type="password" required autocomplete="new-password"
              class="w-full px-4 py-3 rounded-xl border border-border bg-bg-base/60 text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
              placeholder="Минимум 8 символов" />
          </div>

          <p v-if="error" class="text-error text-sm bg-error/8 px-3 py-2 rounded-lg">{{ error }}</p>

          <button type="submit" :disabled="loading"
            class="w-full py-3 bg-primary text-text-inverse rounded-xl text-sm font-semibold tracking-wide hover:bg-primary-hover active:scale-[0.98] transition-all duration-200 disabled:opacity-50 shadow-[0_2px_12px_rgba(184,98,27,0.25)]">
            {{ loading ? 'Создаём...' : 'Создать аккаунт' }}
          </button>
        </form>
      </div>

      <p class="text-center text-sm text-text-tertiary mt-6 anim-fade-up d4">
        Уже есть аккаунт?
        <router-link to="/login" class="text-primary hover:text-primary-hover font-medium transition-colors underline underline-offset-4 decoration-primary/30 hover:decoration-primary">
          Войти
        </router-link>
      </p>
    </div>

    <!-- Recovery Key Modal -->
    <Teleport to="body">
      <div
        v-if="recoveryKey"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
      >
        <div class="bg-bg-base border border-border rounded-2xl p-6 w-full max-w-md mx-4 shadow-xl anim-fade-up">
          <h2 class="text-lg font-semibold text-text-primary mb-1" style="font-family: var(--font-heading)">
            Ключ восстановления
          </h2>
          <p class="text-xs text-text-tertiary mb-4">
            Ваш аккаунт защищён сквозным шифрованием. Сохраните этот ключ в безопасном месте —
            он понадобится, если вы забудете пароль. Ключ показывается только один раз!
          </p>
          <div class="p-3 bg-warning/10 border border-warning/30 rounded-xl mb-4">
            <code class="text-sm text-text-primary break-all select-all">{{ recoveryKey }}</code>
          </div>
          <button
            @click="copyRecoveryKey"
            class="w-full py-2 rounded-xl border border-border text-sm text-text-secondary hover:bg-bg-hover transition mb-2"
          >
            {{ copied ? 'Скопировано!' : 'Скопировать' }}
          </button>
          <button
            @click="finishRegistration"
            class="w-full py-2.5 rounded-xl bg-primary text-text-inverse text-sm font-semibold hover:bg-primary-hover transition"
          >
            Я сохранил ключ
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const email = ref('')
const username = ref('')
const displayName = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const recoveryKey = ref('')
const copied = ref(false)
const router = useRouter()
const authStore = useAuthStore()

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    const key = await authStore.register(email.value, username.value, password.value, displayName.value || username.value)
    recoveryKey.value = key || ''
    // If no recovery key returned (shouldn't happen), go straight to app
    if (!recoveryKey.value) {
      router.push('/')
    }
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}

async function copyRecoveryKey() {
  await navigator.clipboard.writeText(recoveryKey.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

function finishRegistration() {
  recoveryKey.value = ''
  router.push('/')
}
</script>
