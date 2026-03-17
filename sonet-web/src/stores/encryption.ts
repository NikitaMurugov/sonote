import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/composables/useApi'
import { useCrypto } from '@/composables/useCrypto'

const SESSION_KEY = 'sonote-enc-pwd'

export const useEncryptionStore = defineStore('encryption', () => {
  const crypto = useCrypto()

  const isSetup = ref(false)
  const isUnlocked = ref(false)
  const loading = ref(false)

  // In-memory only — never persisted to disk
  let privateKey: CryptoKey | null = null
  const workspaceDEKs = new Map<number, CryptoKey>()

  const hasEncryption = computed(() => isSetup.value)

  async function checkSetup() {
    try {
      const { data } = await api.get('/encryption/keys')
      isSetup.value = !!data.data?.has_encryption
    } catch {
      isSetup.value = false
    }
  }

  async function setup(password: string): Promise<string> {
    loading.value = true
    try {
      const salt = crypto.generateSalt()
      const masterKey = await crypto.deriveMasterKey(password, salt)
      const keyPair = await crypto.generateKeyPair()
      const publicKeyB64 = await crypto.exportPublicKey(keyPair.publicKey)
      const { wrapped, iv } = await crypto.wrapPrivateKey(keyPair.privateKey, masterKey)
      const encryptedPrivateKey = `${iv}:${wrapped}`

      // Generate recovery key
      const recoveryBytes = globalThis.crypto.getRandomValues(new Uint8Array(32))
      const recoveryKeyDisplay = crypto.bufToBase64(recoveryBytes.buffer)
      const recoveryMK = await crypto.deriveMasterKey(recoveryKeyDisplay, salt)
      const recoveryWrapped = await crypto.wrapPrivateKey(keyPair.privateKey, recoveryMK)
      const recoveryDEK = `${recoveryWrapped.iv}:${recoveryWrapped.wrapped}`

      await api.post('/encryption/setup', {
        user_salt: salt,
        encrypted_private_key: encryptedPrivateKey,
        public_key: publicKeyB64,
        recovery_dek: recoveryDEK,
      })

      privateKey = keyPair.privateKey
      isSetup.value = true
      isUnlocked.value = true
      sessionStorage.setItem(SESSION_KEY, password)

      // Create workspace DEK for all user's workspaces
      try {
        const { data: wsData } = await api.get('/workspaces')
        const workspaces = wsData.data || []
        for (const ws of workspaces) {
          const dek = await crypto.generateDEK()
          const wrappedDEK = await crypto.wrapDEKWithPublicKey(dek, keyPair.publicKey)
          await api.put(`/workspaces/${ws.id}/members/me/dek`, { encrypted_dek: wrappedDEK })
          workspaceDEKs.set(ws.id, dek)
        }
      } catch {
        // Non-critical
      }

      return recoveryKeyDisplay
    } finally {
      loading.value = false
    }
  }

  async function unlock(password: string) {
    loading.value = true
    try {
      const { data } = await api.get('/encryption/keys')
      const keys = data.data
      if (!keys.user_salt || !keys.encrypted_private_key) {
        throw new Error('Encryption not set up')
      }

      const masterKey = await crypto.deriveMasterKey(password, keys.user_salt)
      const [iv, wrapped] = (keys.encrypted_private_key as string).split(':')
      privateKey = await crypto.unwrapPrivateKey(wrapped, masterKey, iv)
      isUnlocked.value = true
      sessionStorage.setItem(SESSION_KEY, password)
    } finally {
      loading.value = false
    }
  }

  /** Auto-restore from sessionStorage if private key is lost (e.g. page refresh) */
  async function ensureUnlocked(): Promise<boolean> {
    if (isUnlocked.value && privateKey) return true

    const savedPwd = sessionStorage.getItem(SESSION_KEY)
    if (!savedPwd) return false

    try {
      await unlock(savedPwd)
      return true
    } catch {
      sessionStorage.removeItem(SESSION_KEY)
      return false
    }
  }

  function lock() {
    privateKey = null
    workspaceDEKs.clear()
    isUnlocked.value = false
    sessionStorage.removeItem(SESSION_KEY)
  }

  async function getWorkspaceDEK(workspaceId: number, encryptedDEK?: string): Promise<CryptoKey> {
    const cached = workspaceDEKs.get(workspaceId)
    if (cached) return cached

    // Auto-restore if needed
    if (!privateKey) {
      const restored = await ensureUnlocked()
      if (!restored) throw new Error('Encryption locked')
    }

    let dek64 = encryptedDEK
    if (!dek64) {
      const { data } = await api.get(`/workspaces/${workspaceId}/members/me`)
      dek64 = data.data?.encrypted_dek
      if (!dek64) throw new Error('No DEK for this workspace')
    }

    const dek = await crypto.unwrapDEKWithPrivateKey(dek64, privateKey!)
    workspaceDEKs.set(workspaceId, dek)
    return dek
  }

  async function createWorkspaceDEK(publicKeyB64: string): Promise<{ dek: CryptoKey; encryptedDEK: string }> {
    const dek = await crypto.generateDEK()
    const pubKey = await crypto.importPublicKey(publicKeyB64)
    const encryptedDEK = await crypto.wrapDEKWithPublicKey(dek, pubKey)
    return { dek, encryptedDEK }
  }

  async function wrapDEKForUser(dek: CryptoKey, userPublicKeyB64: string): Promise<string> {
    const pubKey = await crypto.importPublicKey(userPublicKeyB64)
    return crypto.wrapDEKWithPublicKey(dek, pubKey)
  }

  async function encryptNote(
    workspaceId: number,
    title: string,
    contentMd: string,
    contentHtml: string,
    contentJson: any,
  ) {
    const dek = await getWorkspaceDEK(workspaceId)

    const titleEnc = await crypto.encryptText(title, dek)
    const contentEnc = await crypto.encryptText(
      JSON.stringify({ md: contentMd, html: contentHtml, json: contentJson }),
      dek,
    )

    return {
      title_encrypted: titleEnc.ciphertext,
      title_iv: titleEnc.iv,
      content_encrypted: contentEnc.ciphertext,
      content_iv: contentEnc.iv,
      is_encrypted: true,
      title: '🔒',
      content_md: '',
      content_html: '',
      content_json: null,
    }
  }

  async function decryptNote(
    workspaceId: number,
    note: {
      title_encrypted?: string | null
      title_iv?: string | null
      content_encrypted?: string | null
      content_iv?: string | null
      is_encrypted: boolean
    },
  ) {
    if (!note.is_encrypted || !note.title_encrypted || !note.content_encrypted) {
      return null
    }

    const dek = await getWorkspaceDEK(workspaceId)

    const title = await crypto.decryptText(note.title_encrypted, note.title_iv!, dek)
    const contentStr = await crypto.decryptText(note.content_encrypted, note.content_iv!, dek)
    const content = JSON.parse(contentStr)

    return {
      title,
      content_md: content.md || '',
      content_html: content.html || '',
      content_json: content.json || null,
    }
  }

  function getPrivateKey(): CryptoKey | null {
    return privateKey
  }

  return {
    isSetup,
    isUnlocked,
    loading,
    hasEncryption,
    checkSetup,
    setup,
    unlock,
    ensureUnlocked,
    lock,
    getWorkspaceDEK,
    createWorkspaceDEK,
    wrapDEKForUser,
    encryptNote,
    decryptNote,
    getPrivateKey,
  }
})
