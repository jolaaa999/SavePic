<script setup>
import { onUnmounted, ref } from 'vue'
import { resolveFileUrl, updateMeme } from '../api'
import { useIsMobile } from '../composables/useBreakpoint'
import { useMobileSwipeSelect } from '../composables/useMobileSwipeSelect'
import MemeWaterfall from './MemeWaterfall.vue'
import MemePreview from './MemePreview.vue'

defineProps({
  memes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['delete', 'updated', 'batch-delete'])

const isMobile = useIsMobile()
const {
  selectionMode,
  selectedIds,
  selectedCount,
  exitSelection,
  onMemeTouchStart,
  shouldOpenPreview,
  cleanup,
} = useMobileSwipeSelect({ enabled: isMobile })

const toast = ref('')
const previewMeme = ref(null)
const batchDeleting = ref(false)
let toastTimer = null

onUnmounted(cleanup)

function showToast(msg) {
  toast.value = msg
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toast.value = ''
  }, 2000)
}

async function onCopy(meme) {
  try {
    await navigator.clipboard.writeText(resolveFileUrl(meme.file_url))
    showToast('已复制图片地址')
  } catch {
    showToast('复制失败')
  }
}

function onDelete(meme) {
  if (confirm('确定删除这张表情包吗？')) {
    emit('delete', meme.id)
  }
}

async function onEditTags(meme) {
  const current = (meme.tags || []).map((t) => t.name).join(', ')
  const input = window.prompt('编辑标签（逗号分隔）', current)
  if (input === null) return
  const tags = input
    .split(/[,，]/)
    .map((s) => s.trim())
    .filter(Boolean)
  try {
    await updateMeme(meme.id, { tags })
    showToast('标签已更新')
    emit('updated')
  } catch (e) {
    showToast(e.message)
  }
}

async function onBatchDelete() {
  if (!selectedCount.value || batchDeleting.value) return
  const count = selectedCount.value
  if (!confirm(`确定删除选中的 ${count} 张表情包吗？`)) return

  batchDeleting.value = true
  const ids = [...selectedIds.value]
  try {
    emit('batch-delete', ids)
    exitSelection()
  } finally {
    batchDeleting.value = false
  }
}
</script>

<template>
  <div class="relative flex-1 overflow-y-auto px-4 py-4 md:px-6 md:py-5">
    <Transition name="fade">
      <div
        v-if="toast"
        class="fixed left-1/2 z-50 -translate-x-1/2 rounded-lg bg-zinc-800 px-4 py-2 text-sm text-zinc-200 shadow-lg ring-1 ring-white/10 bottom-20 md:bottom-6"
        :class="{ 'bottom-36': selectionMode && isMobile }"
      >
        {{ toast }}
      </div>
    </Transition>

    <p
      v-if="isMobile && !selectionMode && memes.length"
      class="mb-3 text-center text-[10px] text-zinc-600"
    >
      长按进入多选，滑动可批量选中
    </p>

    <MemeWaterfall
      :memes="memes"
      :loading="loading"
      :selectable="isMobile"
      :selection-mode="selectionMode"
      :selected-ids="selectedIds"
      :on-meme-touch-start="onMemeTouchStart"
      :should-open-preview="shouldOpenPreview"
      @copy="onCopy"
      @delete="onDelete"
      @edit-tags="onEditTags"
      @preview="previewMeme = $event"
    />

    <MemePreview :meme="previewMeme" @close="previewMeme = null" />

    <Teleport to="body">
      <Transition name="slide-up">
        <div
          v-if="isMobile && selectionMode"
          class="fixed inset-x-0 z-50 border-t border-[var(--color-border)] bg-[var(--color-surface)]/95 px-4 py-3 backdrop-blur-lg md:hidden"
          style="bottom: calc(3.5rem + env(safe-area-inset-bottom, 0px))"
        >
          <div class="flex items-center gap-3">
            <button
              type="button"
              class="shrink-0 rounded-lg px-3 py-2 text-sm text-zinc-400 active:bg-white/[0.06]"
              @click="exitSelection"
            >
              取消
            </button>
            <p class="min-w-0 flex-1 text-center text-sm text-zinc-300">
              已选 <span class="font-medium text-[var(--color-accent)]">{{ selectedCount }}</span> 张
            </p>
            <button
              type="button"
              class="shrink-0 rounded-lg bg-red-500/15 px-3 py-2 text-sm text-red-400 active:bg-red-500/25 disabled:opacity-40"
              :disabled="!selectedCount || batchDeleting"
              @click="onBatchDelete"
            >
              {{ batchDeleting ? '删除中…' : '删除' }}
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}
</style>
