<template>
  <div class="sonote-editor">
    <EditorToolbar v-if="editor" :editor="editor" />
    <editor-content :editor="editor" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Underline from '@tiptap/extension-underline'
import Highlight from '@tiptap/extension-highlight'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import TextStyle from '@tiptap/extension-text-style'
import Color from '@tiptap/extension-color'
import Image from '@tiptap/extension-image'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Table from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import { WikiLink, createWikiLinkSuggestion } from './WikiLinkExtension'
import EditorToolbar from './EditorToolbar.vue'
import { useDebounce } from '@/composables/useDebounce'
import api from '@/composables/useApi'
import { useWorkspaceStore } from '@/stores/workspace'

const lowlight = createLowlight(common)

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
      codeBlock: false, // replaced by CodeBlockLowlight
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: { class: 'sonote-link' },
    }),
    Placeholder.configure({
      placeholder: 'Начните писать... Введите [[ для ссылки на заметку',
    }),
    Underline,
    Highlight.configure({ multicolor: true }),
    Typography,
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Subscript,
    Superscript,
    TextStyle,
    Color,
    Image.configure({
      inline: false,
      allowBase64: true,
    }),
    TaskList,
    TaskItem.configure({
      nested: true,
    }),
    Table.configure({
      resizable: true,
    }),
    TableRow,
    TableCell,
    TableHeader,
    CodeBlockLowlight.configure({
      lowlight,
    }),
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

/* Placeholder */
.sonote-prose p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: var(--color-text-tertiary);
  pointer-events: none;
  height: 0;
  font-style: italic;
}

/* Headings */
.sonote-prose h1 {
  font-family: var(--font-heading);
  font-size: 2em;
  font-weight: 600;
  line-height: 1.3;
  margin: 1.2em 0 0.5em;
  letter-spacing: -0.01em;
}
.sonote-prose h2 {
  font-family: var(--font-heading);
  font-size: 1.5em;
  font-weight: 600;
  line-height: 1.35;
  margin: 1em 0 0.4em;
}
.sonote-prose h3 {
  font-family: var(--font-heading);
  font-size: 1.2em;
  font-weight: 600;
  line-height: 1.4;
  margin: 0.8em 0 0.3em;
}

/* Links */
.sonote-prose a, .sonote-prose .sonote-link {
  color: var(--color-link);
  text-decoration: underline;
  text-underline-offset: 3px;
  text-decoration-thickness: 1px;
  transition: color 0.15s;
}
.sonote-prose a:hover {
  color: var(--color-link-hover);
}

/* Wiki links */
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

/* Inline code */
.sonote-prose code {
  font-family: var(--font-mono);
  font-size: 0.88em;
  background: var(--color-bg-hover);
  padding: 0.15em 0.45em;
  border-radius: 5px;
  color: var(--color-accent);
}

/* Code blocks */
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

/* Blockquote */
.sonote-prose blockquote {
  border-left: 3px solid var(--color-primary);
  padding-left: 1.2em;
  color: var(--color-text-secondary);
  font-style: italic;
  margin: 1.2em 0;
}

/* Lists */
.sonote-prose ul, .sonote-prose ol {
  padding-left: 1.6em;
  margin: 0.6em 0;
}
.sonote-prose li { margin: 0.25em 0; }
.sonote-prose li::marker { color: var(--color-text-tertiary); }

/* Task list (checklist) */
.sonote-prose ul[data-type="taskList"] {
  padding-left: 0;
  list-style: none;
}
.sonote-prose ul[data-type="taskList"] li {
  display: flex;
  align-items: flex-start;
  gap: 0.5em;
  margin: 0.4em 0;
}
.sonote-prose ul[data-type="taskList"] li > label {
  flex-shrink: 0;
  margin-top: 0.15em;
}
.sonote-prose ul[data-type="taskList"] li > label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
  cursor: pointer;
  border-radius: 3px;
}
.sonote-prose ul[data-type="taskList"] li > div {
  flex: 1;
}
.sonote-prose ul[data-type="taskList"] li[data-checked="true"] > div {
  text-decoration: line-through;
  color: var(--color-text-tertiary);
}

/* Highlight */
.sonote-prose mark {
  background-color: var(--color-primary-light);
  padding: 0.1em 0.2em;
  border-radius: 3px;
}

/* Horizontal rule */
.sonote-prose hr {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--color-border), transparent);
  margin: 2em 0;
}

/* Images */
.sonote-prose img {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 1em 0;
}
.sonote-prose img.ProseMirror-selectednode {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

/* Tables */
.sonote-prose table {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
  overflow: hidden;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}
.sonote-prose th,
.sonote-prose td {
  border: 1px solid var(--color-border-light);
  padding: 0.6em 0.8em;
  text-align: left;
  font-size: 0.92em;
  vertical-align: top;
  min-width: 80px;
}
.sonote-prose th {
  background: var(--color-bg-sidebar);
  font-weight: 600;
  font-size: 0.85em;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--color-text-secondary);
}
.sonote-prose td {
  background: var(--color-bg-surface);
}
.sonote-prose .selectedCell {
  background: var(--color-primary-light) !important;
}

/* Text alignment */
.sonote-prose .has-text-align-center { text-align: center; }
.sonote-prose .has-text-align-right { text-align: right; }
.sonote-prose .has-text-align-justify { text-align: justify; }

/* Subscript / Superscript */
.sonote-prose sub { font-size: 0.75em; }
.sonote-prose sup { font-size: 0.75em; }

/* Underline */
.sonote-prose u { text-underline-offset: 3px; }

/* Strong / em */
.sonote-prose strong { font-weight: 600; }
.sonote-prose em { font-style: italic; }

/* Syntax highlighting (lowlight) */
.sonote-prose pre .hljs-keyword { color: var(--color-accent); }
.sonote-prose pre .hljs-string { color: var(--color-success); }
.sonote-prose pre .hljs-comment { color: var(--color-text-tertiary); font-style: italic; }
.sonote-prose pre .hljs-number { color: var(--color-warning); }
.sonote-prose pre .hljs-function { color: var(--color-primary); }
.sonote-prose pre .hljs-title { color: var(--color-primary); }
.sonote-prose pre .hljs-built_in { color: var(--color-accent); }
.sonote-prose pre .hljs-type { color: var(--color-warning); }
.sonote-prose pre .hljs-attr { color: var(--color-link); }
.sonote-prose pre .hljs-variable { color: var(--color-accent-soft); }
.sonote-prose pre .hljs-literal { color: var(--color-error); }
</style>
