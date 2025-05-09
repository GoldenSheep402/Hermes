<!-- eslint-disable vue/v-on-event-hyphenation -->
<!-- eslint-disable no-alert -->
<script lang="ts" setup>
import CodeBlock from '@tiptap/extension-code-block'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Table from '@tiptap/extension-table'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import TableRow from '@tiptap/extension-table-row'
import TaskItem from '@tiptap/extension-task-item'
import TaskList from '@tiptap/extension-task-list'
import TextAlign from '@tiptap/extension-text-align'
import Underline from '@tiptap/extension-underline'
import StarterKit from '@tiptap/starter-kit'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import { onBeforeUnmount, watch } from 'vue'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  minHeight?: string
  theme?: 'light' | 'dark'
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      history: {
        depth: 100,
        newGroupDelay: 500,
      },
    }),
    Image.configure({
      inline: true,
      allowBase64: true,
    }),
    Link.configure({
      openOnClick: false,
    }),
    Placeholder.configure({
      placeholder: props.placeholder || '请输入内容...',
    }),
    Table.configure({
      resizable: true,
    }),
    TableRow,
    TableHeader,
    TableCell,
    TaskList,
    TaskItem.configure({
      nested: true,
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Underline,
    CodeBlock.configure({
      HTMLAttributes: {
        class: 'code-block',
      },
    }),
    Highlight.configure({
      multicolor: true,
    }),
  ],
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
})

watch(() => props.modelValue, (newValue) => {
  const editorContent = editor.value?.getHTML()
  if (newValue !== editorContent) {
    editor.value?.commands.setContent(newValue)
  }
})

watch(() => props.theme, (newValue) => {
  if (editor.value) {
    const element = editor.value.view.dom as HTMLElement
    if (newValue === 'dark') {
      element.classList.add('prose-invert')
    }
    else {
      element.classList.remove('prose-invert')
    }
  }
})

function setLink() {
  const previousUrl = editor.value?.getAttributes('link').href
  const url = window.prompt('URL', previousUrl)

  // cancelled
  if (url === null) {
    return
  }

  // empty
  if (url === '') {
    editor.value?.chain().focus().unsetLink().run()
    return
  }

  // update link
  editor.value?.chain().focus().setLink({ href: url }).run()
}

function addImage() {
  const url = window.prompt('URL')

  if (url) {
    editor.value?.chain().focus().setImage({ src: url }).run()
  }
}

// 监听主题变化
watch(
  () => props.theme,
  (newTheme) => {
    if (editor.value) {
      editor.value.setOptions({
        editorProps: {
          attributes: {
            class: newTheme === 'dark' ? 'dark-theme' : 'light-theme',
          },
        },
      })
    }
  },
  { immediate: true },
)

// 监听内容变化
watch(
  () => props.modelValue,
  (value) => {
    const isSame = editor.value?.getHTML() === value
    if (editor.value && !isSame) {
      editor.value.commands.setContent(value, false)
    }
  },
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<template>
  <div class="editor-container">
    <div v-if="editor" class="editor-toolbar">
      <a-space>
        <a-tooltip content="粗体">
          <a-button
            :type="editor.isActive('bold') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleBold().run()"
          >
            <template #icon>
              <icon-bold />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="斜体">
          <a-button
            :type="editor.isActive('italic') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleItalic().run()"
          >
            <template #icon>
              <icon-italic />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="下划线">
          <a-button
            :type="editor.isActive('underline') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleUnderline().run()"
          >
            <template #icon>
              <icon-underline />
            </template>
          </a-button>
        </a-tooltip>
        <a-divider direction="vertical" />
        <a-tooltip content="标题">
          <a-dropdown trigger="click">
            <a-button type="text">
              <template #icon>
                <icon-h1 />
              </template>
            </a-button>
            <template #content>
              <a-doption @click="editor.chain().focus().toggleHeading({ level: 1 }).run()">
                <template #icon>
                  <icon-h1 />
                </template>
                标题 1
              </a-doption>
              <a-doption @click="editor.chain().focus().toggleHeading({ level: 2 }).run()">
                <template #icon>
                  <icon-h2 />
                </template>
                标题 2
              </a-doption>
              <a-doption @click="editor.chain().focus().toggleHeading({ level: 3 }).run()">
                <template #icon>
                  <icon-h3 />
                </template>
                标题 3
              </a-doption>
            </template>
          </a-dropdown>
        </a-tooltip>
        <a-tooltip content="有序列表">
          <a-button
            :type="editor.isActive('orderedList') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleOrderedList().run()"
          >
            <template #icon>
              <icon-ordered-list />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="无序列表">
          <a-button
            :type="editor.isActive('bulletList') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleBulletList().run()"
          >
            <template #icon>
              <icon-unordered-list />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="任务列表">
          <a-button
            :type="editor.isActive('taskList') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleTaskList().run()"
          >
            <template #icon>
              <icon-select-all />
            </template>
          </a-button>
        </a-tooltip>
        <a-divider direction="vertical" />
        <a-tooltip content="代码块">
          <a-button
            :type="editor.isActive('codeBlock') ? 'primary' : 'text'"
            @click="editor.chain().focus().toggleCodeBlock().run()"
          >
            <template #icon>
              <icon-code />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="表格">
          <a-button
            :type="editor.isActive('table') ? 'primary' : 'text'"
            @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()"
          >
            <template #icon>
              <icon-apps />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="图片">
          <a-button type="text" @click="addImage">
            <template #icon>
              <icon-image />
            </template>
          </a-button>
        </a-tooltip>
        <a-tooltip content="链接">
          <a-button
            :type="editor.isActive('link') ? 'primary' : 'text'"
            @click="setLink"
          >
            <template #icon>
              <icon-link />
            </template>
          </a-button>
        </a-tooltip>
        <a-divider direction="vertical" />
        <a-tooltip content="对齐方式">
          <a-dropdown trigger="click">
            <a-button type="text">
              <template #icon>
                <icon-align-left />
              </template>
            </a-button>
            <template #content>
              <a-doption @click="editor.chain().focus().setTextAlign('left').run()">
                <template #icon>
                  <icon-align-left />
                </template>
                左对齐
              </a-doption>
              <a-doption @click="editor.chain().focus().setTextAlign('center').run()">
                <template #icon>
                  <icon-align-center />
                </template>
                居中
              </a-doption>
              <a-doption @click="editor.chain().focus().setTextAlign('right').run()">
                <template #icon>
                  <icon-align-right />
                </template>
                右对齐
              </a-doption>
            </template>
          </a-dropdown>
        </a-tooltip>
      </a-space>
    </div>
    <EditorContent :editor="editor" class="editor-content" />
  </div>
</template>

<style scoped>
.editor-container {
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-bg-2);
}

.editor-toolbar {
  padding: 8px;
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-bg-2);
}

.editor-content {
  padding: 16px;
  min-height: v-bind(minHeight);
}

:deep(.ProseMirror) {
  outline: none;
  min-height: v-bind(minHeight);
}

:deep(.ProseMirror p.is-editor-empty:first-child::before) {
  color: var(--color-text-3);
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}
</style>
