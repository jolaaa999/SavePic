<script setup>
/**
 * 高性能瀑布流组件 — 基于后端 width/height 预计算列高，ResizeObserver 响应式分列
 *
 * @typedef {Object} MemeTag
 * @property {number} id
 * @property {string} name
 *
 * @typedef {Object} Meme
 * @property {number} id
 * @property {number} category_id
 * @property {string} file_url
 * @property {string} file_hash
 * @property {number} width
 * @property {number} height
 * @property {number} size
 * @property {string} created_at
 * @property {MemeTag[]} [tags]
 */

import { computed, ref, toRef } from 'vue'
import { resolveFileUrl } from '../api'
import { getColumnCount, useWaterfallLayout } from '../composables/useWaterfallLayout'

const props = defineProps({
  /** @type {import('vue').PropType<Meme[]>} */
  memes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  gap: { type: Number, default: 12 },
  emptyText: { type: String, default: '暂无表情包' },
  emptyHint: { type: String, default: '试试调整标签筛选或上传新图片' },
})

const emit = defineEmits(['copy', 'delete', 'edit-tags', 'preview'])

const containerRef = ref(null)
const memesRef = toRef(props, 'memes')

const { columns, columnCount, gap } = useWaterfallLayout(containerRef, memesRef, {
  gap: props.gap,
})

const skeletonColumns = computed(() => getColumnCount(containerRef.value?.clientWidth || 375))

function aspectStyle(meme) {
  if (meme.width > 0 && meme.height > 0) {
    return { aspectRatio: `${meme.width} / ${meme.height}` }
  }
  return { aspectRatio: '1' }
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>

<template>
  <section
    ref="containerRef"
    class="meme-waterfall relative w-full"
    :style="{ '--wf-gap': `${gap}px` }"
  >
    <!-- Loading skeleton -->
    <div
      v-if="loading"
      class="flex"
      :style="{ gap: `${gap}px` }"
    >
      <div
        v-for="c in skeletonColumns"
        :key="`sk-${c}`"
        class="min-w-0 flex-1 flex flex-col"
        :style="{ gap: `${gap}px` }"
      >
        <div
          v-for="i in 3"
          :key="i"
          class="animate-pulse rounded-xl bg-white/[0.04]"
          :style="{ aspectRatio: i % 2 ? '3/4' : '1/1' }"
        />
      </div>
    </div>

    <!-- Empty -->
    <div
      v-else-if="!memes.length"
      class="flex h-40 flex-col items-center justify-center rounded-2xl border border-dashed border-white/[0.06] text-zinc-600 md:h-48"
    >
      <slot name="empty">
        <p class="text-center text-sm">{{ emptyText }}</p>
        <p v-if="emptyHint" class="mt-1 text-center text-xs text-zinc-600">{{ emptyHint }}</p>
      </slot>
    </div>

    <!-- Waterfall columns -->
    <div
      v-else
      class="flex items-start"
      :style="{ gap: `${gap}px` }"
    >
      <div
        v-for="(col, colIndex) in columns"
        :key="colIndex"
        class="min-w-0 flex-1 flex flex-col"
        :style="{ gap: `${gap}px` }"
      >
        <article
          v-for="meme in col"
          :key="meme.id"
          class="group relative overflow-hidden rounded-xl border border-white/[0.04] bg-[var(--color-surface-raised)] transition duration-200 will-change-transform md:hover:border-[var(--color-accent)]/25 md:hover:shadow-[0_0_24px_-4px_rgba(110,231,183,0.12)]"
        >
          <!-- 图片区：后端 width/height 占位；操作层仅覆盖图片，不遮挡下方标签 -->
          <div class="relative w-full bg-white/[0.02]" :style="aspectStyle(meme)">
            <img
              :src="resolveFileUrl(meme.file_url)"
              :alt="`meme-${meme.id}`"
              :width="meme.width || undefined"
              :height="meme.height || undefined"
              class="absolute inset-0 h-full w-full cursor-zoom-in object-cover"
              loading="lazy"
              decoding="async"
              @click="emit('preview', meme)"
            />

            <slot name="actions" :meme="meme">
              <div
                class="pointer-events-none absolute inset-0 flex flex-col items-center justify-end gap-1.5 bg-gradient-to-t from-black/80 via-black/25 to-transparent p-2 opacity-100 transition-opacity md:opacity-0 md:group-hover:opacity-100 md:p-3"
              >
                <div class="pointer-events-auto flex flex-wrap justify-center gap-1.5">
                  <button
                    type="button"
                    class="rounded-lg bg-white/10 px-2 py-1 text-[11px] text-zinc-200 backdrop-blur-sm md:hover:bg-white/20"
                    @click.stop="emit('copy', meme)"
                  >
                    复制
                  </button>
                  <button
                    type="button"
                    class="rounded-lg bg-white/10 px-2 py-1 text-[11px] text-zinc-200 backdrop-blur-sm md:hover:bg-white/20"
                    @click.stop="emit('edit-tags', meme)"
                  >
                    标签
                  </button>
                  <button
                    type="button"
                    class="rounded-lg bg-red-500/20 px-2 py-1 text-[11px] text-red-300 backdrop-blur-sm md:hover:bg-red-500/35"
                    @click.stop="emit('delete', meme)"
                  >
                    删除
                  </button>
                </div>
              </div>
            </slot>

            <p
              class="pointer-events-none absolute left-2 top-2 rounded bg-black/50 px-1.5 py-0.5 text-[10px] text-zinc-400 md:opacity-0 md:group-hover:opacity-100"
            >
              {{ formatDate(meme.created_at) }}
            </p>
          </div>

          <!-- 标签 -->
          <div v-if="meme.tags?.length" class="flex flex-wrap gap-1 px-2 py-1.5">
            <span
              v-for="tag in meme.tags"
              :key="tag.id"
              class="rounded bg-white/[0.06] px-1.5 py-0.5 text-[10px] text-zinc-500"
            >
              #{{ tag.name }}
            </span>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.meme-waterfall {
  contain: layout style;
}
</style>
