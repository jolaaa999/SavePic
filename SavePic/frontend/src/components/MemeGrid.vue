<script setup>
import { ref } from 'vue'
import { resolveFileUrl, updateMeme } from '../api'
import MemeWaterfall from './MemeWaterfall.vue'
import MemePreview from './MemePreview.vue'

defineProps({
  memes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['delete', 'updated'])

const toast = ref('')
const previewMeme = ref(null)
let toastTimer = null

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
</script>

<template>
  <div class="relative flex-1 overflow-y-auto px-4 py-4 md:px-6 md:py-5">
    <Transition name="fade">
      <div
        v-if="toast"
        class="fixed left-1/2 z-50 -translate-x-1/2 rounded-lg bg-zinc-800 px-4 py-2 text-sm text-zinc-200 shadow-lg ring-1 ring-white/10 bottom-20 md:bottom-6"
      >
        {{ toast }}
      </div>
    </Transition>

    <MemeWaterfall
      :memes="memes"
      :loading="loading"
      @copy="onCopy"
      @delete="onDelete"
      @edit-tags="onEditTags"
      @preview="previewMeme = $event"
    />

    <MemePreview :meme="previewMeme" @close="previewMeme = null" />
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
</style>
