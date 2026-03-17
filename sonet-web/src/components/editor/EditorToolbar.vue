<template>
  <div class="flex items-center gap-0.5 flex-wrap py-2 px-1 border-b border-border-light/60 mb-4 sticky top-0 bg-bg-base/95 backdrop-blur-sm z-10">
    <!-- Text style group -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="Жирный (Ctrl+B)"
        :active="editor.isActive('bold')"
        @click="editor.chain().focus().toggleBold().run()"
      >
        <span class="font-bold text-xs">B</span>
      </ToolbarButton>
      <ToolbarButton
        title="Курсив (Ctrl+I)"
        :active="editor.isActive('italic')"
        @click="editor.chain().focus().toggleItalic().run()"
      >
        <span class="italic text-xs">I</span>
      </ToolbarButton>
      <ToolbarButton
        title="Подчёркнутый (Ctrl+U)"
        :active="editor.isActive('underline')"
        @click="editor.chain().focus().toggleUnderline().run()"
      >
        <span class="underline text-xs">U</span>
      </ToolbarButton>
      <ToolbarButton
        title="Зачёркнутый"
        :active="editor.isActive('strike')"
        @click="editor.chain().focus().toggleStrike().run()"
      >
        <span class="line-through text-xs">S</span>
      </ToolbarButton>
      <ToolbarButton
        title="Выделение"
        :active="editor.isActive('highlight')"
        @click="editor.chain().focus().toggleHighlight().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Код"
        :active="editor.isActive('code')"
        @click="editor.chain().focus().toggleCode().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
        </svg>
      </ToolbarButton>
    </div>

    <!-- Headings -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="Заголовок 1"
        :active="editor.isActive('heading', { level: 1 })"
        @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      >
        <span class="text-[10px] font-bold">H1</span>
      </ToolbarButton>
      <ToolbarButton
        title="Заголовок 2"
        :active="editor.isActive('heading', { level: 2 })"
        @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      >
        <span class="text-[10px] font-bold">H2</span>
      </ToolbarButton>
      <ToolbarButton
        title="Заголовок 3"
        :active="editor.isActive('heading', { level: 3 })"
        @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
      >
        <span class="text-[10px] font-bold">H3</span>
      </ToolbarButton>
      <ToolbarButton
        title="Параграф"
        :active="editor.isActive('paragraph') && !editor.isActive('heading')"
        @click="editor.chain().focus().setParagraph().run()"
      >
        <span class="text-[10px] font-medium">P</span>
      </ToolbarButton>
    </div>

    <!-- Lists -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="Маркированный список"
        :active="editor.isActive('bulletList')"
        @click="editor.chain().focus().toggleBulletList().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16" />
          <circle cx="1" cy="6" r="1" fill="currentColor" />
          <circle cx="1" cy="12" r="1" fill="currentColor" />
          <circle cx="1" cy="18" r="1" fill="currentColor" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Нумерованный список"
        :active="editor.isActive('orderedList')"
        @click="editor.chain().focus().toggleOrderedList().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 6h13M8 12h13M8 18h13" />
          <text x="1" y="8" font-size="7" fill="currentColor" font-family="sans-serif">1</text>
          <text x="1" y="14" font-size="7" fill="currentColor" font-family="sans-serif">2</text>
          <text x="1" y="20" font-size="7" fill="currentColor" font-family="sans-serif">3</text>
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Чеклист"
        :active="editor.isActive('taskList')"
        @click="editor.chain().focus().toggleTaskList().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      </ToolbarButton>
    </div>

    <!-- Block elements -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="Цитата"
        :active="editor.isActive('blockquote')"
        @click="editor.chain().focus().toggleBlockquote().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
          <path d="M4.583 17.321C3.553 16.227 3 15 3 13.011c0-3.5 2.457-6.637 6.03-8.188l.893 1.378c-3.335 1.804-3.987 4.145-4.247 5.621.537-.278 1.24-.375 1.929-.311C9.591 11.68 11 13.175 11 15c0 1.933-1.567 3.5-3.5 3.5-1.102 0-2.152-.494-2.917-1.179zM14.583 17.321C13.553 16.227 13 15 13 13.011c0-3.5 2.457-6.637 6.03-8.188l.893 1.378c-3.335 1.804-3.987 4.145-4.247 5.621.537-.278 1.24-.375 1.929-.311C19.591 11.68 21 13.175 21 15c0 1.933-1.567 3.5-3.5 3.5-1.102 0-2.152-.494-2.917-1.179z"/>
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Блок кода"
        :active="editor.isActive('codeBlock')"
        @click="editor.chain().focus().toggleCodeBlock().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Разделитель"
        @click="editor.chain().focus().setHorizontalRule().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-width="1.5" d="M3 12h18" />
        </svg>
      </ToolbarButton>
    </div>

    <!-- Alignment -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="По левому краю"
        :active="editor.isActive({ textAlign: 'left' })"
        @click="editor.chain().focus().setTextAlign('left').run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-width="1.5" d="M3 6h18M3 10h12M3 14h18M3 18h12" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="По центру"
        :active="editor.isActive({ textAlign: 'center' })"
        @click="editor.chain().focus().setTextAlign('center').run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-width="1.5" d="M3 6h18M6 10h12M3 14h18M6 18h12" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="По правому краю"
        :active="editor.isActive({ textAlign: 'right' })"
        @click="editor.chain().focus().setTextAlign('right').run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-width="1.5" d="M3 6h18M9 10h12M3 14h18M9 18h12" />
        </svg>
      </ToolbarButton>
    </div>

    <!-- Table + Image -->
    <div class="flex items-center gap-0.5 pr-2 mr-1 border-r border-border-light">
      <ToolbarButton
        title="Таблица"
        @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 10h18M3 14h18M10 3v18M14 3v18M3 6a3 3 0 013-3h12a3 3 0 013 3v12a3 3 0 01-3 3H6a3 3 0 01-3-3V6z" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Изображение по URL"
        @click="insertImage"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </ToolbarButton>
    </div>

    <!-- Undo / Redo -->
    <div class="flex items-center gap-0.5">
      <ToolbarButton
        title="Отменить (Ctrl+Z)"
        :disabled="!editor.can().undo()"
        @click="editor.chain().focus().undo().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4" />
        </svg>
      </ToolbarButton>
      <ToolbarButton
        title="Повторить (Ctrl+Shift+Z)"
        :disabled="!editor.can().redo()"
        @click="editor.chain().focus().redo().run()"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 10H11a5 5 0 00-5 5v2M21 10l-4-4M21 10l-4 4" />
        </svg>
      </ToolbarButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import ToolbarButton from './ToolbarButton.vue'

const props = defineProps<{
  editor: Editor
}>()

function insertImage() {
  const url = window.prompt('URL изображения:')
  if (url) {
    props.editor.chain().focus().setImage({ src: url }).run()
  }
}
</script>
