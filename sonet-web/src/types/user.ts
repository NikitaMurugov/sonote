export interface User {
  id: number
  email: string
  username: string
  display_name: string
  avatar_url: string | null
  email_verified: boolean
  user_salt?: string | null
  encrypted_private_key?: string | null
  public_key?: string | null
  created_at: string
  updated_at: string
}
