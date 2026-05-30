<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { resolveFileUrl } from '../api'

const props = defineProps({
  /** @type {import('vue').PropType<{ id: number, file_url: string, width?: number, height?: number, tags?: { name: string }[], created_at?: string } | null>} */
  meme: { type: Object, default: null },
  memes: { type: Array, default: () => [] },
})

const emit = defineEmits(['close', 'update:meme'])

const SWIPE_THRESHOLD = 48

const touchStartX = ref(0)
const touchStartY = ref(0)

const currentIndex = computed(() => {
  if (!props.meme) return -1
  return props.memes.findIndex((m) => m.id === props.meme.id)
})

const hasPrev = computed(() => currentIndex.value > 0)
const hasNext = computed(
  () => currentIndex.value >= 0 && currentIndex.value < props.memes.length - 1,
)

const positionLabel = computed(() => {
  if (currentIndex.value < 0 || !props.memes.length) return ''
  return `${currentIndex.value + 1} / ${props.memes.length}`
})

function goPrev() {
  if (!hasPrev.value) return
  emit('update:meme', props.memes[currentIndex.value - 1])
}

function goNext() {
  if (!hasNext.value) return
  emit('update:meme', props.memes[currentIndex.value + 1])
}

function onKeydown(e) {
  if (e.key === 'Escape') {
    emit('close')
    return
  }
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    goPrev()
  }
  if (e.key === 'ArrowRight') {
    e.preventDefault()
    goNext()
  }
}

function onTouchStart(e) {
  const touch = e.touches[0]
  if (!touch) return
  touchStartX.value = touch.clientX
  touchStartY.value = touch.clientY
}

function onTouchEnd(e) {
  const touch = e.changedTouches[0]
  if (!touch) return

  const dx = touch.clientX - touchStartX.value
  const dy = touch.clientY - touchStartY.value
  if (Math.abs(dx) < SWIPE_THRESHOLD || Math.abs(dx) < Math.abs(dy)) return

  if (dx > 0) goPrev()
  else goNext()
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
            <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-0.5">
              <p v-if="meme.created_at" class="text-[11px] text-zinc-600">
                {{ formatDate(meme.created_at) }}
              </p>
              <p v-if="positionLabel" class="text-[11px] text-zinc-500">
                {{ positionLabel }}
              </p>
            </div>
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

        <div class="relative flex min-h-0 flex-1 items-center justify-center px-4 pt-0 md:px-16">
          <button
            v-if="memes.length > 1"
            type="button"
            class="absolute left-2 top-1/2 z-10 hidden -translate-y-1/2 rounded-full border border-white/10 bg-black/40 p-2.5 text-zinc-300 backdrop-blur-sm transition hover:bg-black/60 hover:text-zinc-100 disabled:pointer-events-none disabled:opacity-30 md:flex"
            :disabled="!hasPrev"
            aria-label="上一张"
            @click.stop="goPrev"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
            </svg>
          </button>

          <div
            class="flex h-full w-full items-center justify-center"
            @touchstart.passive="onTouchStart"
            @touchend="onTouchEnd"
          >
            <Transition name="slide-image" mode="out-in">
              <img
                :key="meme.id"
                :src="resolveFileUrl(meme.file_url)"
                :alt="`meme-${meme.id}`"
                class="max-h-full max-w-full object-contain select-none"
                draggable="false"
                @click.stop
              />
            </Transition>
          </div>

          <button
            v-if="memes.length > 1"
            type="button"
            class="absolute right-2 top-1/2 z-10 hidden -translate-y-1/2 rounded-full border border-white/10 bg-black/40 p-2.5 text-zinc-300 backdrop-blur-sm transition hover:bg-black/60 hover:text-zinc-100 disabled:pointer-events-none disabled:opacity-30 md:flex"
            :disabled="!hasNext"
            aria-label="下一张"
            @click.stop="goNext"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
            </svg>
          </button>
        </div>

        <p class="shrink-0 pb-[max(1rem,env(safe-area-inset-bottom))] text-center text-[11px] text-zinc-600">
          <span class="md:hidden">左右滑动切换图片 · </span>
          <span class="hidden md:inline">方向键或左右按钮切换 · </span>
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

.slide-image-enter-active,
.slide-image-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.slide-image-enter-from {
  opacity: 0;
  transform: scale(0.98);
}
.slide-image-leave-to {
  opacity: 0;
  transform: scale(0.98);
}
</style>
