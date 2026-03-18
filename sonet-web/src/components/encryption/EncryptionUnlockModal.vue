<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="encryptionStore.showUnlockModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
        @click.self="handleCancel"
      >
        <div class="bg-bg-card rounded-2xl shadow-xl p-6 w-full max-w-sm mx-4 border border-border-light">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center">
              <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <div>
              <h3 class="text-base font-semibold text-text-primary" style="font-family: var(--font-heading)">
                Разблокировать шифрование
              </h3>
              <p class="text-xs text-text-tertiary">Введите пароль аккаунта</p>
            </div>
          </div>

          <form @submit.prevent="handleSubmit">
            <input
              ref="passwordInput"
              v-model="password"
              type="password"
              placeholder="Пароль"
              autocomplete="current-password"
              class="w-full px-4 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary placeholder-text-tertiary text-sm focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition"
              :disabled="submitting"
            />

            <p v-if="error" class="mt-2 text-xs text-error">{{ error }}</p>

            <div class="flex gap-2 mt-4">
              <button
                type="button"
                @click="handleCancel"
                class="flex-1 py-2.5 rounded-xl text-sm font-medium text-text-secondary border border-border hover:bg-bg-hover transition"
                :disabled="submitting"
              >
                Позже
              </button>
              <button
                type="submit"
                class="flex-1 py-2.5 rounded-xl text-sm font-medium text-text-inverse bg-primary hover:bg-primary-hover transition disabled:opacity-50"
                :disabled="submitting || !password"
              >
                {{ submitting ? 'Проверяю...' : 'Разблокировать' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useEncryptionStore } from '@/stores/encryption'

const encryptionStore = useEncryptionStore()

const password = ref('')
const error = ref('')
const submitting = ref(false)
const passwordInput = ref<HTMLInputElement>()

watch(() => encryptionStore.showUnlockModal, async (show) => {
  if (show) {
    password.value = ''
    error.value = ''
    submitting.value = false
    await nextTick()
    passwordInput.value?.focus()
  }
})

async function handleSubmit() {
  if (!password.value || submitting.value) return
  submitting.value = true
  error.value = ''
  try {
    await encryptionStore.resolveUnlock(password.value)
  } catch {
    error.value = 'Неверный пароль'
    submitting.value = false
  }
}

function handleCancel() {
  encryptionStore.cancelUnlock()
}
</script>

<style scoped>
.modal-enter-active { transition: opacity 0.2s ease; }
.modal-enter-active > div { transition: transform 0.2s cubic-bezier(0.22, 1, 0.36, 1), opacity 0.2s ease; }
.modal-leave-active { transition: opacity 0.15s ease; }
.modal-leave-active > div { transition: transform 0.15s ease, opacity 0.15s ease; }
.modal-enter-from { opacity: 0; }
.modal-enter-from > div { transform: scale(0.95); opacity: 0; }
.modal-leave-to { opacity: 0; }
.modal-leave-to > div { transform: scale(0.95); opacity: 0; }
</style>
