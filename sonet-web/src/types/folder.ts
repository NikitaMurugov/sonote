export interface Folder {
  id: number
  workspace_id: number
  parent_id: number | null
  name: string
  sort_order: number
  created_at: string
  updated_at: string
}
