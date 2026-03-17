<template>
  <div class="flex items-center justify-center h-full">
    <p class="text-text-secondary">Загрузка...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'

const router = useRouter()
const workspaceStore = useWorkspaceStore()

onMounted(async () => {
  await workspaceStore.fetchWorkspaces()
  if (workspaceStore.workspaces.length > 0) {
    const personal = workspaceStore.workspaces.find((w) => w.is_personal)
    const target = personal || workspaceStore.workspaces[0]
    router.replace(`/w/${target.slug}`)
  }
})
</script>
