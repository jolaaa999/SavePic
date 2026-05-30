<script setup>
import { ref } from 'vue'

const props = defineProps({
  categories: { type: Array, required: true },
  selectedId: { type: Number, default: null },
  loading: { type: Boolean, default: false },
  reordering: { type: Boolean, default: false },
})

const emit = defineEmits(['select', 'new-category', 'delete-category', 'rename-category', 'reorder-categories'])

const draggingId = ref(null)
const dragOverId = ref(null)

/**
 * @param {MouseEvent} e
 * @param {{ id: number, name: string, count?: number }} cat
 */
function onDeleteClick(e, cat) {
  e.stopPropagation()
  emit('delete-category', cat)
}

/**
 * @param {MouseEvent} e
 * @param {{ id: number, name: string, count?: number }} cat
 */
function onRenameClick(e, cat) {
  e.stopPropagation()
  emit('rename-category', cat)
}

/** @param {DragEvent} e @param {{ id: number }} cat */
function onDragStart(e, cat) {
  if (props.reordering) {
    e.preventDefault()
    return
  }
  draggingId.value = cat.id
  dragOverId.value = cat.id
  e.dataTransfer?.setData('text/plain', String(cat.id))
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move'
}

function onDragEnd() {
  draggingId.value = null
  dragOverId.value = null
}

/** @param {DragEvent} e @param {{ id: number }} cat */
function onDragOver(e, cat) {
  if (!draggingId.value || draggingId.value === cat.id) return
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  dragOverId.value = cat.id
}

/** @param {{ id: number }} cat */
function onDragLeave(cat) {
  if (dragOverId.value === cat.id) dragOverId.value = null
}

/** @param {DragEvent} e @param {{ id: number }} cat */
function onDrop(e, cat) {
  e.preventDefault()
  const fromId = draggingId.value
  onDragEnd()
  if (!fromId || fromId === cat.id) return

  const ids = props.categories.map((c) => c.id)
  const fromIndex = ids.indexOf(fromId)
  const toIndex = ids.indexOf(cat.id)
  if (fromIndex < 0 || toIndex < 0) return

  ids.splice(fromIndex, 1)
  ids.splice(toIndex, 0, fromId)
  emit('reorder-categories', ids)
}
</script>

<template>
  <aside
    class="relative z-10 hidden w-60 shrink-0 flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)] md:flex"
  >
    <!-- Logo -->
    <div class="flex items-center gap-2.5 px-5 py-6">
      <div
        class="flex h-8 w-8 items-center justify-center rounded-lg bg-[var(--color-accent-muted)] text-sm font-semibold text-[var(--color-accent)]"
      >
        S
      </div>
      <div>
        <p class="text-sm font-medium tracking-tight text-zinc-100">SavePic</p>
        <p class="text-[11px] text-zinc-600">表情包库</p>
      </div>
    </div>

    <div class="flex items-center justify-between px-5 pb-2">
      <p class="text-[10px] font-medium uppercase tracking-widest text-zinc-600">分类</p>
      <p v-if="categories.length > 1" class="text-[10px] text-zinc-600">拖动排序</p>
    </div>

    <!-- 分类列表 -->
    <nav class="flex-1 space-y-0.5 overflow-y-auto px-3">
      <p v-if="loading" class="px-3 py-4 text-center text-xs text-zinc-600">加载中…</p>
      <p v-else-if="!categories.length" class="px-3 py-4 text-center text-xs text-zinc-600">
        暂无分类，点击下方新建
      </p>
      <button
        v-for="cat in categories"
        :key="cat.id"
        type="button"
        class="group flex w-full items-center rounded-lg px-2 py-2.5 text-left transition-all duration-150"
        :class="[
          selectedId === cat.id
            ? 'bg-[var(--color-accent-muted)] text-zinc-100'
            : 'text-zinc-500 hover:bg-white/[0.03] hover:text-zinc-300',
          draggingId === cat.id ? 'opacity-40' : '',
          dragOverId === cat.id && draggingId !== cat.id
            ? 'ring-1 ring-inset ring-[var(--color-accent)]/40'
            : '',
        ]"
        @click="$emit('select', cat.id)"
        @dragover="onDragOver($event, cat)"
        @dragleave="onDragLeave(cat)"
        @drop="onDrop($event, cat)"
      >
        <span
          v-if="categories.length > 1"
          draggable="true"
          class="mr-1.5 flex shrink-0 cursor-grab touch-none rounded p-1 text-zinc-600 opacity-0 transition-all hover:bg-white/10 hover:text-zinc-400 active:cursor-grabbing group-hover:opacity-100"
          :class="[
            selectedId === cat.id ? 'opacity-100' : '',
            reordering ? 'pointer-events-none opacity-30' : '',
          ]"
          title="拖动排序"
          aria-label="拖动排序"
          @click.stop
          @dragstart="onDragStart($event, cat)"
          @dragend="onDragEnd"
        >
          <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="9" cy="7" r="1.4" />
            <circle cx="15" cy="7" r="1.4" />
            <circle cx="9" cy="12" r="1.4" />
            <circle cx="15" cy="12" r="1.4" />
            <circle cx="9" cy="17" r="1.4" />
            <circle cx="15" cy="17" r="1.4" />
          </svg>
        </span>

        <span class="min-w-0 flex-1 truncate text-sm font-medium">{{ cat.name }}</span>

        <span class="ml-2 flex shrink-0 items-center gap-1">
          <span
            class="rounded-md px-1.5 py-0.5 text-[10px] tabular-nums transition-colors"
            :class="
              selectedId === cat.id
                ? 'bg-[var(--color-accent)]/20 text-[var(--color-accent)]'
                : 'bg-white/5 text-zinc-600 group-hover:text-zinc-500'
            "
          >
            {{ cat.count }}
          </span>
          <button
            type="button"
            class="rounded-md p-1 text-zinc-600 opacity-0 transition-all hover:bg-white/10 hover:text-zinc-300 group-hover:opacity-100"
            :class="selectedId === cat.id ? 'opacity-100' : ''"
            title="重命名分类"
            aria-label="重命名分类"
            @click="onRenameClick($event, cat)"
          >
            <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125"
              />
            </svg>
          </button>
          <button
            type="button"
            class="rounded-md p-1 text-zinc-600 opacity-0 transition-all hover:bg-red-500/15 hover:text-red-400 group-hover:opacity-100"
            :class="selectedId === cat.id ? 'opacity-100' : ''"
            title="删除分类"
            aria-label="删除分类"
            @click="onDeleteClick($event, cat)"
          >
            <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
              />
            </svg>
          </button>
        </span>
      </button>
    </nav>

    <!-- 新建分类 -->
    <div class="border-t border-[var(--color-border)] p-3">
      <button
        type="button"
        class="flex w-full items-center justify-center gap-2 rounded-lg border border-dashed border-white/10 px-3 py-2.5 text-sm text-zinc-500 transition-all hover:border-[var(--color-accent)]/40 hover:bg-[var(--color-accent-muted)] hover:text-[var(--color-accent)]"
        @click="$emit('new-category')"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        新建分类
      </button>
    </div>
  </aside>
</template>
