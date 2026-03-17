<template>
  <div class="px-10 py-10 max-w-3xl mx-auto anim-fade-up">
    <h1 class="text-2xl font-semibold mb-8" style="font-family: var(--font-heading)">Настройки</h1>

    <!-- Tabs -->
    <div class="flex gap-1 mb-8 border-b border-border-light">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="[
          'px-4 py-2 text-sm font-medium transition-colors border-b-2 -mb-px',
          activeTab === tab.id
            ? 'border-primary text-primary'
            : 'border-transparent text-text-secondary hover:text-text-primary',
        ]"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- General -->
    <div v-if="activeTab === 'general' && ws">
      <div class="space-y-5">
        <div>
          <label class="block text-xs font-medium text-text-secondary uppercase tracking-wide mb-2">Название</label>
          <input
            v-model="wsName"
            class="w-full px-4 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-text-secondary uppercase tracking-wide mb-2">Описание</label>
          <textarea
            v-model="wsDescription"
            rows="3"
            class="w-full px-4 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all resize-none"
          />
        </div>
        <button
          @click="saveWorkspace"
          class="px-4 py-2 bg-primary text-text-inverse rounded-lg text-sm font-medium hover:bg-primary-hover transition"
        >
          Сохранить
        </button>
      </div>
    </div>

    <!-- Members -->
    <div v-if="activeTab === 'members'">
      <!-- Invite -->
      <div class="flex gap-3 mb-8">
        <input
          v-model="inviteEmail"
          placeholder="Email пользователя"
          class="flex-1 px-4 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary text-sm placeholder-text-tertiary focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
        />
        <select
          v-model="inviteRole"
          class="px-3 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary text-sm focus:outline-none focus:border-primary transition"
        >
          <option value="viewer">Viewer</option>
          <option value="editor">Editor</option>
          <option value="admin">Admin</option>
        </select>
        <button
          @click="inviteMember"
          class="px-4 py-2.5 bg-primary text-text-inverse rounded-xl text-sm font-medium hover:bg-primary-hover transition"
        >
          Пригласить
        </button>
      </div>

      <p v-if="inviteError" class="text-error text-sm mb-4">{{ inviteError }}</p>

      <!-- Member list -->
      <div class="space-y-3">
        <div
          v-for="member in members"
          :key="member.id"
          class="flex items-center gap-3 p-4 rounded-2xl bg-bg-surface/60 border border-border-light"
        >
          <div class="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold">
            {{ member.user_id }}
          </div>
          <div class="flex-1">
            <p class="text-sm text-text-primary">User #{{ member.user_id }}</p>
            <p class="text-[11px] text-text-tertiary capitalize">{{ member.role }}</p>
          </div>
          <button
            @click="removeMember(member.user_id)"
            class="text-xs text-text-tertiary hover:text-error transition-colors"
          >
            Удалить
          </button>
        </div>
        <p v-if="!members.length" class="text-sm text-text-tertiary text-center py-4">
          Участников пока нет
        </p>
      </div>
    </div>

    <!-- Tags -->
    <div v-if="activeTab === 'tags'">
      <div class="flex gap-3 mb-6">
        <input
          v-model="newTagName"
          @keydown.enter="createTag"
          placeholder="Новый тег"
          class="flex-1 px-4 py-2.5 rounded-xl border border-border bg-bg-base text-text-primary text-sm placeholder-text-tertiary focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
        />
        <input
          v-model="newTagColor"
          type="color"
          class="w-10 h-10 rounded-lg border border-border cursor-pointer"
        />
        <button
          @click="createTag"
          class="px-4 py-2.5 bg-primary text-text-inverse rounded-xl text-sm font-medium hover:bg-primary-hover transition"
        >
          Создать
        </button>
      </div>
      <div class="flex flex-wrap gap-3">
        <div
          v-for="tag in tagStore.tags"
          :key="tag.id"
          class="group flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-border-light bg-bg-surface/60"
        >
          <span
            class="w-2.5 h-2.5 rounded-full shrink-0"
            :style="{ backgroundColor: tag.color || '#A89888' }"
          ></span>
          <span class="text-sm text-text-primary">{{ tag.name }}</span>
          <button
            @click="deleteTag(tag.id)"
            class="opacity-0 group-hover:opacity-100 text-text-tertiary hover:text-error transition-all ml-1"
          >
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <p v-if="!tagStore.tags.length" class="text-sm text-text-tertiary text-center py-4">
        Тегов пока нет
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useTagStore } from '@/stores/tag'
import { useEncryptionStore } from '@/stores/encryption'
import api from '@/composables/useApi'
import type { WorkspaceMember } from '@/types/workspace'

const route = useRoute()
const workspaceStore = useWorkspaceStore()
const tagStore = useTagStore()
const encryptionStore = useEncryptionStore()

const tabs = [
  { id: 'general', label: 'Общее' },
  { id: 'members', label: 'Участники' },
  { id: 'tags', label: 'Теги' },
]

const activeTab = ref('general')
const ws = ref(workspaceStore.currentWorkspace)
const wsName = ref('')
const wsDescription = ref('')
const members = ref<WorkspaceMember[]>([])
const inviteEmail = ref('')
const inviteRole = ref('viewer')
const inviteError = ref('')
const newTagName = ref('')
const newTagColor = ref('#E8A87C')

async function loadMembers() {
  if (!ws.value) return
  try {
    const { data } = await api.get(`/workspaces/${ws.value.id}/members`)
    members.value = data.data || []
  } catch {
    members.value = []
  }
}

async function inviteMember() {
  if (!ws.value || !inviteEmail.value.trim()) return
  inviteError.value = ''
  try {
    const payload: Record<string, any> = {
      email: inviteEmail.value,
      role: inviteRole.value,
    }

    // If workspace is encrypted, wrap DEK for the invited user
    if (ws.value.is_encrypted) {
      const unlocked = await encryptionStore.ensureUnlocked()
      if (!unlocked) {
        inviteError.value = 'Разблокируйте шифрование в профиле'
        return
      }
      try {
        const { data: pkData } = await api.get('/encryption/public-key-by-email', {
          params: { email: inviteEmail.value },
        })
        const inviteePublicKey = pkData.data?.public_key
        if (inviteePublicKey) {
          const dek = await encryptionStore.getWorkspaceDEK(ws.value.id)
          payload.encrypted_dek = await encryptionStore.wrapDEKForUser(dek, inviteePublicKey)
        }
      } catch (e: any) {
        console.warn('DEK wrapping failed:', e)
        // Proceed without DEK — invitee won't be able to decrypt until re-invited
      }
    }

    await api.post(`/workspaces/${ws.value.id}/members`, payload)
    inviteEmail.value = ''
    await loadMembers()
  } catch (e: any) {
    inviteError.value = e.response?.data?.error || 'Ошибка приглашения'
  }
}

async function removeMember(userId: number) {
  if (!ws.value) return
  await api.delete(`/workspaces/${ws.value.id}/members/${userId}`)
  await loadMembers()
}

async function saveWorkspace() {
  if (!ws.value) return
  await api.patch(`/workspaces/${ws.value.id}`, {
    name: wsName.value,
    description: wsDescription.value,
  })
}

async function createTag() {
  if (!ws.value || !newTagName.value.trim()) return
  await tagStore.createTag(ws.value.id, newTagName.value.trim(), newTagColor.value)
  newTagName.value = ''
}

async function deleteTag(tagId: number) {
  if (!ws.value) return
  await api.delete(`/workspaces/${ws.value.id}/tags/${tagId}`)
  await tagStore.fetchTags(ws.value.id)
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  await workspaceStore.setCurrentBySlug(slug)
  ws.value = workspaceStore.currentWorkspace
  if (ws.value) {
    wsName.value = ws.value.name
    wsDescription.value = ws.value.description || ''
    await loadMembers()
    await tagStore.fetchTags(ws.value.id)
  }
})
</script>
