<script setup>
import { ref } from 'vue'

const props = defineProps({
  tags: { type: Array, default: () => [] },
  selectedIds: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['update:selectedIds', 'rename-tag'])

const longPressTimer = ref(null)
const longPressTriggered = ref(false)

function toggle(id) {
  if (longPressTriggered.value) {
    longPressTriggered.value = false
    return
  }
  const set = new Set(props.selectedIds)
  if (set.has(id)) {
    set.delete(id)
  } else {
    set.add(id)
  }
  emit('update:selectedIds', [...set])
}

function clearAll() {
  emit('update:selectedIds', [])
}

/** @param {{ id: number, name: string }} tag */
function startLongPress(tag) {
  clearLongPress()
  longPressTimer.value = window.setTimeout(() => {
    longPressTriggered.value = true
    emit('rename-tag', tag)
  }, 500)
}

function clearLongPress() {
  if (longPressTimer.value) {
    clearTimeout(longPressTimer.value)
    longPressTimer.value = null
  }
}

/** @param {MouseEvent} e @param {{ id: number, name: string }} tag */
function onContextMenu(e, tag) {
  e.preventDefault()
  emit('rename-tag', tag)
}
</script>

<template>
  <div class="border-b border-[var(--color-border)] px-4 py-2.5 md:px-6">
    <div class="mb-2 flex items-center justify-between">
      <span class="text-[10px] font-medium uppercase tracking-widest text-zinc-600">标签筛选</span>
      <button
        v-if="selectedIds.length"
        type="button"
        class="text-[11px] text-zinc-500 hover:text-[var(--color-accent)]"
        @click="clearAll"
      >
        清除
      </button>
    </div>

    <p v-if="loading" class="text-xs text-zinc-600">加载标签…</p>
    <p v-else-if="!tags.length" class="text-xs text-zinc-600">暂无标签，上传时可添加</p>

    <div v-else class="-mx-1 flex gap-1.5 overflow-x-auto pb-1 scrollbar-none">
      <button
        v-for="tag in tags"
        :key="tag.id"
        type="button"
        class="shrink-0 rounded-full px-3 py-1 text-xs transition"
        :class="
          selectedIds.includes(tag.id)
            ? 'bg-[var(--color-accent-muted)] text-[var(--color-accent)] ring-1 ring-[var(--color-accent)]/30'
            : 'bg-white/[0.04] text-zinc-500 active:bg-white/[0.08]'
        "
        @click="toggle(tag.id)"
        @contextmenu="onContextMenu($event, tag)"
        @touchstart.passive="startLongPress(tag)"
        @touchend="clearLongPress"
        @touchmove="clearLongPress"
        @touchcancel="clearLongPress"
      >
        #{{ tag.name }}
        <span class="ml-1 opacity-60">{{ tag.count }}</span>
      </button>
    </div>

    <p v-if="tags.length && selectedIds.length" class="mt-1.5 text-[10px] text-zinc-600">
      已选 {{ selectedIds.length }} 个标签（同时满足）· 长按或右键可重命名
    </p>
    <p v-else-if="tags.length" class="mt-1.5 text-[10px] text-zinc-600">
      长按或右键标签可重命名
    </p>
  </div>
</template>
