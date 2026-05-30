<script setup>
import { onMounted, onUnmounted, watch } from 'vue'
import { resolveFileUrl } from '../api'

const props = defineProps({
  /** @type {import('vue').PropType<{ id: number, file_url: string, width?: number, height?: number, tags?: { name: string }[], created_at?: string } | null>} */
  meme: { type: Object, default: null },
})

const emit = defineEmits(['close'])

function onKeydown(e) {
  if (e.key === 'Escape') emit('close')
}

watch(
  () => props.meme,
  (meme) => {
    document.body.style.overflow = meme ? 'hidden' : ''
  },
)

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="preview">
      <div
        v-if="meme"
        class="fixed inset-0 z-[100] flex flex-col bg-black/90 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        :aria-label="`预览表情包 ${meme.id}`"
        @click.self="emit('close')"
      >
        <div class="flex shrink-0 items-center justify-between px-4 py-3 md:px-6">
          <div class="min-w-0 flex-1">
            <p v-if="meme.tags?.length" class="flex flex-wrap gap-1.5">
              <span
                v-for="tag in meme.tags"
                :key="tag.id ?? tag.name"
                class="rounded bg-white/10 px-2 py-0.5 text-xs text-zinc-400"
              >
                #{{ tag.name }}
              </span>
            </p>
            <p v-if="meme.created_at" class="mt-1 text-[11px] text-zinc-600">
              {{ formatDate(meme.created_at) }}
            </p>
          </div>
          <button
            type="button"
            class="ml-4 shrink-0 rounded-lg p-2 text-zinc-400 transition hover:bg-white/10 hover:text-zinc-200"
            aria-label="关闭预览"
            @click="emit('close')"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="flex min-h-0 flex-1 items-center justify-center p-4 pt-0">
          <img
            :src="resolveFileUrl(meme.file_url)"
            :alt="`meme-${meme.id}`"
            class="max-h-full max-w-full object-contain select-none"
            draggable="false"
            @click.stop
          />
        </div>

        <p class="shrink-0 pb-[max(1rem,env(safe-area-inset-bottom))] text-center text-[11px] text-zinc-600">
          点击空白处或按 Esc 关闭
        </p>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.preview-enter-active,
.preview-leave-active {
  transition: opacity 0.2s ease;
}
.preview-enter-from,
.preview-leave-to {
  opacity: 0;
}
</style>
