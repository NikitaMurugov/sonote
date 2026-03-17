import { Node, mergeAttributes } from '@tiptap/core'
import { VueRenderer } from '@tiptap/vue-3'
import Suggestion from '@tiptap/suggestion'
import type { SuggestionOptions } from '@tiptap/suggestion'
import WikiLinkList from './WikiLinkList.vue'

export interface WikiLinkOptions {
  suggestion: Partial<SuggestionOptions>
}

export const WikiLink = Node.create<WikiLinkOptions>({
  name: 'wikiLink',
  group: 'inline',
  inline: true,
  selectable: false,
  atom: true,

  addAttributes() {
    return {
      noteId: { default: null },
      title: { default: '' },
    }
  },

  parseHTML() {
    return [{ tag: 'a[data-wiki-link]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return [
      'a',
      mergeAttributes(HTMLAttributes, {
        'data-wiki-link': '',
        class: 'wiki-link',
        href: '#',
      }),
      `[[${HTMLAttributes.title}]]`,
    ]
  },

  renderText({ node }) {
    return `[[${node.attrs.title}]]`
  },

  addKeyboardShortcuts() {
    return {}
  },

  addProseMirrorPlugins() {
    return [
      Suggestion({
        editor: this.editor,
        ...this.options.suggestion,
      }),
    ]
  },
})

export function createWikiLinkSuggestion(fetchNotes: () => Promise<Array<{ id: number; title: string }>>) {
  return {
    char: '[[',
    allowSpaces: true,

    items: async ({ query }: { query: string }) => {
      const notes = await fetchNotes()
      return notes
        .filter((n) => n.title.toLowerCase().includes(query.toLowerCase()))
        .slice(0, 8)
    },

    render: () => {
      let component: VueRenderer
      let popup: HTMLElement

      return {
        onStart: (props: any) => {
          component = new VueRenderer(WikiLinkList, {
            props,
            editor: props.editor,
          })

          popup = document.createElement('div')
          popup.style.position = 'absolute'
          popup.style.zIndex = '9999'
          popup.appendChild(component.element!)
          document.body.appendChild(popup)

          updatePosition(props)
        },

        onUpdate(props: any) {
          component.updateProps(props)
          updatePosition(props)
        },

        onKeyDown(props: any) {
          if (props.event.key === 'Escape') {
            popup?.remove()
            component?.destroy()
            return true
          }
          return component?.ref?.onKeyDown?.(props.event) ?? false
        },

        onExit() {
          popup?.remove()
          component?.destroy()
        },
      }

      function updatePosition(props: any) {
        if (!props.clientRect || !popup) return
        const rect = props.clientRect()
        if (!rect) return
        popup.style.left = `${rect.left}px`
        popup.style.top = `${rect.bottom + 4}px`
      }
    },

    command: ({ editor, range, props }: any) => {
      editor
        .chain()
        .focus()
        .deleteRange(range)
        .insertContent({
          type: 'wikiLink',
          attrs: {
            noteId: props.id,
            title: props.title,
          },
        })
        .insertContent(' ')
        .run()
    },
  }
}
