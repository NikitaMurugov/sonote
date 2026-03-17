const PBKDF2_ITERATIONS = 600_000
const AES_KEY_LENGTH = 256
const RSA_MODULUS_LENGTH = 2048

function bufToBase64(buf: ArrayBuffer): string {
  return btoa(String.fromCharCode(...new Uint8Array(buf)))
}

function base64ToBuf(b64: string): ArrayBuffer {
  const bin = atob(b64)
  const buf = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; i++) buf[i] = bin.charCodeAt(i)
  return buf.buffer
}

function generateSalt(): string {
  return bufToBase64(crypto.getRandomValues(new Uint8Array(32)).buffer)
}

function generateIV(): string {
  return bufToBase64(crypto.getRandomValues(new Uint8Array(12)).buffer)
}

async function deriveMasterKey(password: string, salt: string): Promise<CryptoKey> {
  const enc = new TextEncoder()
  const keyMaterial = await crypto.subtle.importKey(
    'raw',
    enc.encode(password),
    'PBKDF2',
    false,
    ['deriveKey'],
  )
  return crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt: base64ToBuf(salt),
      iterations: PBKDF2_ITERATIONS,
      hash: 'SHA-256',
    },
    keyMaterial,
    { name: 'AES-GCM', length: AES_KEY_LENGTH },
    true,
    ['wrapKey', 'unwrapKey'],
  )
}

async function generateKeyPair(): Promise<CryptoKeyPair> {
  return crypto.subtle.generateKey(
    {
      name: 'RSA-OAEP',
      modulusLength: RSA_MODULUS_LENGTH,
      publicExponent: new Uint8Array([1, 0, 1]),
      hash: 'SHA-256',
    },
    true,
    ['wrapKey', 'unwrapKey'],
  )
}

async function exportPublicKey(key: CryptoKey): Promise<string> {
  const exported = await crypto.subtle.exportKey('spki', key)
  return bufToBase64(exported)
}

async function importPublicKey(b64: string): Promise<CryptoKey> {
  return crypto.subtle.importKey(
    'spki',
    base64ToBuf(b64),
    { name: 'RSA-OAEP', hash: 'SHA-256' },
    true,
    ['wrapKey'],
  )
}

async function wrapPrivateKey(
  privateKey: CryptoKey,
  masterKey: CryptoKey,
): Promise<{ wrapped: string; iv: string }> {
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const wrapped = await crypto.subtle.wrapKey('pkcs8', privateKey, masterKey, {
    name: 'AES-GCM',
    iv,
  })
  return { wrapped: bufToBase64(wrapped), iv: bufToBase64(iv.buffer) }
}

async function unwrapPrivateKey(
  wrappedB64: string,
  masterKey: CryptoKey,
  ivB64: string,
): Promise<CryptoKey> {
  return crypto.subtle.unwrapKey(
    'pkcs8',
    base64ToBuf(wrappedB64),
    masterKey,
    { name: 'AES-GCM', iv: base64ToBuf(ivB64) },
    { name: 'RSA-OAEP', hash: 'SHA-256' },
    true,
    ['unwrapKey'],
  )
}

async function generateDEK(): Promise<CryptoKey> {
  return crypto.subtle.generateKey(
    { name: 'AES-GCM', length: AES_KEY_LENGTH },
    true,
    ['encrypt', 'decrypt'],
  )
}

async function wrapDEKWithPublicKey(dek: CryptoKey, publicKey: CryptoKey): Promise<string> {
  const wrapped = await crypto.subtle.wrapKey('raw', dek, publicKey, { name: 'RSA-OAEP' })
  return bufToBase64(wrapped)
}

async function unwrapDEKWithPrivateKey(
  wrappedB64: string,
  privateKey: CryptoKey,
): Promise<CryptoKey> {
  return crypto.subtle.unwrapKey(
    'raw',
    base64ToBuf(wrappedB64),
    privateKey,
    { name: 'RSA-OAEP' },
    { name: 'AES-GCM', length: AES_KEY_LENGTH },
    true,
    ['encrypt', 'decrypt'],
  )
}

async function encryptText(text: string, dek: CryptoKey): Promise<{ ciphertext: string; iv: string }> {
  const enc = new TextEncoder()
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const encrypted = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv },
    dek,
    enc.encode(text),
  )
  return { ciphertext: bufToBase64(encrypted), iv: bufToBase64(iv.buffer) }
}

async function decryptText(ciphertextB64: string, ivB64: string, dek: CryptoKey): Promise<string> {
  const dec = new TextDecoder()
  const decrypted = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: base64ToBuf(ivB64) },
    dek,
    base64ToBuf(ciphertextB64),
  )
  return dec.decode(decrypted)
}

export function useCrypto() {
  return {
    generateSalt,
    generateIV,
    deriveMasterKey,
    generateKeyPair,
    exportPublicKey,
    importPublicKey,
    wrapPrivateKey,
    unwrapPrivateKey,
    generateDEK,
    wrapDEKWithPublicKey,
    unwrapDEKWithPrivateKey,
    encryptText,
    decryptText,
    bufToBase64,
    base64ToBuf,
  }
}
