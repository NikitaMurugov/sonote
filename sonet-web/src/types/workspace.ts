export interface Workspace {
  id: number
  name: string
  slug: string
  description: string | null
  owner_id: number
  is_personal: boolean
  icon: string | null
  created_at: string
  updated_at: string
}

export interface WorkspaceMember {
  id: number
  workspace_id: number
  user_id: number
  role: 'viewer' | 'editor' | 'admin'
  invited_by: number | null
  joined_at: string
}
