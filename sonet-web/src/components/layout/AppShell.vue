<template>
  <div class="h-screen flex bg-bg-base overflow-hidden">
    <Transition name="sidebar">
      <Sidebar v-show="!uiStore.sidebarCollapsed || !isMobile" />
    </Transition>
    <main class="flex-1 flex flex-col overflow-hidden relative">
      <TopBar />
      <div class="flex-1 overflow-auto">
        <router-view v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </router-view>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'
import Sidebar from './Sidebar.vue'
import TopBar from './TopBar.vue'

const authStore = useAuthStore()
const uiStore = useUiStore()
const isMobile = ref(false)

onMounted(async () => {
  if (authStore.isAuthenticated && !authStore.user) {
    await authStore.fetchUser()
  }
  isMobile.value = window.innerWidth < 768
})
</script>

<style scoped>
.page-enter-active { animation: fadeUp 0.35s cubic-bezier(0.22, 1, 0.36, 1) both; }
.page-leave-active { animation: fadeOut 0.15s ease both; }
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
@keyframes fadeOut {
  from { opacity: 1; }
  to { opacity: 0; }
}

.sidebar-enter-active { transition: transform 0.25s cubic-bezier(0.22, 1, 0.36, 1), opacity 0.25s ease; }
.sidebar-leave-active { transition: transform 0.2s ease, opacity 0.2s ease; }
.sidebar-enter-from { transform: translateX(-100%); opacity: 0; }
.sidebar-leave-to { transform: translateX(-100%); opacity: 0; }
</style>
