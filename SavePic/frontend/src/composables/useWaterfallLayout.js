import { ref, watch, onMounted, onUnmounted, unref } from 'vue'

/**
 * @typedef {Object} MemeTag
 * @property {number} id
 * @property {string} name
 */

/**
 * @typedef {Object} MemeItem
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

const TAGS_BLOCK_HEIGHT = 32
const DEFAULT_ASPECT = 1

/**
 * 根据容器宽度计算列数（移动端固定 2 列）
 * @param {number} width
 */
export function getColumnCount(width) {
  if (width < 768) return 2
  if (width < 1024) return 4
  if (width < 1280) return 5
  return 6
}

/**
 * 预估卡片高度（优先使用后端 width/height，避免布局跳动）
 * @param {MemeItem} meme
 * @param {number} columnWidth
 */
export function estimateItemHeight(meme, columnWidth) {
  const ratio =
    meme.width > 0 && meme.height > 0 ? meme.height / meme.width : DEFAULT_ASPECT
  const imageHeight = columnWidth * ratio
  const tagsHeight = meme.tags?.length ? TAGS_BLOCK_HEIGHT : 0
  return imageHeight + tagsHeight
}

/**
 * 最短列贪心分配
 * @param {MemeItem[]} memes
 * @param {number} columnCount
 * @param {number} columnWidth
 * @param {number} gap
 * @returns {MemeItem[][]}
 */
export function distributeToColumns(memes, columnCount, columnWidth, gap) {
  const cols = Array.from({ length: columnCount }, () => /** @type {MemeItem[]} */ ([]))
  const heights = Array(columnCount).fill(0)

  for (const meme of memes) {
    const itemH = estimateItemHeight(meme, columnWidth) + gap
    let minIdx = 0
    for (let i = 1; i < columnCount; i++) {
      if (heights[i] < heights[minIdx]) minIdx = i
    }
    cols[minIdx].push(meme)
    heights[minIdx] += itemH
  }

  return cols
}

/**
 * @param {import('vue').Ref<HTMLElement | null>} containerRef
 * @param {import('vue').Ref<MemeItem[]> | import('vue').ComputedRef<MemeItem[]>} memesSource
 * @param {{ gap?: number }} [options]
 */
export function useWaterfallLayout(containerRef, memesSource, options = {}) {
  const gap = options.gap ?? 12
  const columnCount = ref(2)
  const columns = ref(/** @type {MemeItem[][]} */ ([]))
  const columnWidth = ref(0)

  let resizeObserver = null

  function relayout() {
    const el = containerRef.value
    const memes = unref(memesSource) || []
    if (!el || memes.length === 0) {
      columns.value = []
      return
    }

    const width = el.clientWidth
    const count = getColumnCount(width)
    const totalGap = gap * (count - 1)
    const colW = (width - totalGap) / count

    columnCount.value = count
    columnWidth.value = colW
    columns.value = distributeToColumns(memes, count, colW, gap)
  }

  onMounted(() => {
    const el = containerRef.value
    if (!el) return

    resizeObserver = new ResizeObserver(() => {
      relayout()
    })
    resizeObserver.observe(el)
    relayout()
  })

  onUnmounted(() => {
    resizeObserver?.disconnect()
  })

  watch(
    () => [unref(memesSource), unref(memesSource)?.length],
    () => relayout(),
    { deep: true },
  )

  return { columns, columnCount, columnWidth, gap, relayout }
}
