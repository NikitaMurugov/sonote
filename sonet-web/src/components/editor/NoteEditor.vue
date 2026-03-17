<template>
  <div class="sonote-editor">
    <editor-content :editor="editor" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { WikiLink, createWikiLinkSuggestion } from './WikiLinkExtension'
import { useDebounce } from '@/composables/useDebounce'
import api from '@/composables/useApi'
import { useWorkspaceStore } from '@/stores/workspace'

const props = defineProps<{
  content: any
}>()

const emit = defineEmits<{
  update: [payload: { html: string; json: any; text: string }]
}>()

const workspaceStore = useWorkspaceStore()

async function fetchNotes() {
  if (!workspaceStore.currentWorkspace) return []
  try {
    const { data } = await api.get(`/workspaces/${workspaceStore.currentWorkspace.id}/notes`)
    return (data.data || []).map((n: any) => ({ id: n.id, title: n.title }))
  } catch {
    return []
  }
}

const debouncedEmit = useDebounce((editor: any) => {
  emit('update', {
    html: editor.getHTML(),
    json: editor.getJSON(),
    text: editor.getText(),
  })
}, 500)

const editor = useEditor({
  content: props.content,
  extensions: [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] },
    }),
    Link.configure({ openOnClick: false }),
    Placeholder.configure({ placeholder: 'Начните писать... Введите [[ для ссылки на заметку' }),
    WikiLink.configure({
      suggestion: createWikiLinkSuggestion(fetchNotes),
    }),
  ],
  editorProps: {
    attributes: {
      class: 'sonote-prose',
    },
  },
  onUpdate: ({ editor }) => {
    debouncedEmit(editor)
  },
})

watch(
  () => props.content,
  (newContent) => {
    if (editor.value && newContent !== editor.value.getHTML()) {
      editor.value.commands.setContent(newContent, false)
    }
  }
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style>
.sonote-prose {
  min-height: 400px;
  font-family: var(--font-body);
  font-size: 15.5px;
  line-height: 1.75;
  color: var(--color-text-primary);
  outline: none;
  caret-color: var(--color-primary);
}

.sonote-prose > * + * {
  margin-top: 0.75em;
}

.sonote-prose p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: var(--color-text-tertiary);
  pointer-events: none;
  height: 0;
  font-style: italic;
}

.sonote-prose h1 {
  font-family: var(--font-heading);
  font-size: 2em;
  font-weight: 600;
  line-height: 1.3;
  margin: 1.2em 0 0.5em;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}
.sonote-prose h2 {
  font-family: var(--font-heading);
  font-size: 1.5em;
  font-weight: 600;
  line-height: 1.35;
  margin: 1em 0 0.4em;
  color: var(--color-text-primary);
}
.sonote-prose h3 {
  font-family: var(--font-heading);
  font-size: 1.2em;
  font-weight: 600;
  line-height: 1.4;
  margin: 0.8em 0 0.3em;
  color: var(--color-text-primary);
}

.sonote-prose a {
  color: var(--color-link);
  text-decoration: underline;
  text-decoration-color: var(--color-link);
  text-underline-offset: 3px;
  text-decoration-thickness: 1px;
  transition: color 0.15s, text-decoration-color 0.15s;
}
.sonote-prose a:hover {
  color: var(--color-link-hover);
  text-decoration-color: var(--color-link-hover);
}

/* Wiki link styling */
.sonote-prose a.wiki-link {
  color: var(--color-primary);
  text-decoration: none;
  background: var(--color-primary-light);
  padding: 0.1em 0.35em;
  border-radius: 4px;
  font-size: 0.95em;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.sonote-prose a.wiki-link:hover {
  background: var(--color-primary);
  color: var(--color-text-inverse);
}

.sonote-prose code {
  font-family: var(--font-mono);
  font-size: 0.88em;
  background: var(--color-bg-hover);
  padding: 0.15em 0.45em;
  border-radius: 5px;
  color: var(--color-accent);
}

.sonote-prose pre {
  background: var(--color-bg-sidebar);
  border: 1px solid var(--color-border-light);
  padding: 1.2em 1.4em;
  border-radius: 10px;
  overflow-x: auto;
  font-size: 0.88em;
  line-height: 1.6;
}
.sonote-prose pre code {
  background: none;
  padding: 0;
  color: var(--color-text-primary);
  font-size: 1em;
}

.sonote-prose blockquote {
  border-left: 3px solid var(--color-primary);
  padding-left: 1.2em;
  color: var(--color-text-secondary);
  font-style: italic;
  margin: 1.2em 0;
}

.sonote-prose ul, .sonote-prose ol {
  padding-left: 1.6em;
  margin: 0.6em 0;
}
.sonote-prose li { margin: 0.25em 0; }
.sonote-prose li::marker { color: var(--color-text-tertiary); }

.sonote-prose hr {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--color-border), transparent);
  margin: 2em 0;
}

.sonote-prose strong { font-weight: 600; }
.sonote-prose em { font-style: italic; }
</style>
