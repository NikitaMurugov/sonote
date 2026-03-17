export interface Note {
  id: number
  workspace_id: number
  folder_id: number | null
  title: string
  slug: string
  content_md: string
  content_html: string
  content_json: any
  author_id: number
  is_pinned: boolean
  is_archived: boolean
  word_count: number
  created_at: string
  updated_at: string
}
