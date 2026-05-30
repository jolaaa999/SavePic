<script setup>
import { computed, ref, watch } from 'vue'
import { useIsMobile } from '../composables/useBreakpoint'
import SearchBar from './SearchBar.vue'
import TagFilter from './TagFilter.vue'

defineProps({
  categoryName: { type: String, default: '' },
  total: { type: Number, default: 0 },
  tags: { type: Array, default: () => [] },
  loadingTags: { type: Boolean, default: false },
})

const query = defineModel('query', { type: String, default: '' })
const sort = defineModel('sort', { type: String, default: 'desc' })
const selectedIds = defineModel('selectedIds', { type: Array, default: () => [] })

defineEmits(['rename-tag', 'delete-tag'])

const isMobile = useIsMobile()
const filtersExpanded = ref(true)
const uploadExpanded = ref(true)

watch(
  isMobile,
  (mobile) => {
    filtersExpanded.value = !mobile
  },
  { immediate: true },
)

const hasActiveFilters = computed(
  () => query.value.trim().length > 0 || selectedIds.value.length > 0,
)

const filterSummary = computed(() => {
  const parts = []
  if (query.value.trim()) parts.push('含搜索')
  if (selectedIds.value.length) parts.push(`${selectedIds.value.length} 个标签`)
  return parts.join(' · ')
})
</script>

<template>
  <div class="shrink-0">
    <!-- 顶栏：分类信息 + 折叠按钮 -->
    <div
      class="sticky top-0 z-10 border-b border-[var(--color-border)] bg-[var(--color-surface)]/90 backdrop-blur-md"
    >
      <div class="flex items-center gap-3 px-4 py-3 md:px-6">
        <div class="min-w-0 flex-1">
          <h1 class="truncate text-base font-medium text-zinc-100">
            {{ categoryName || '选择分类' }}
          </h1>
          <p class="mt-0.5 text-xs text-zinc-600">
            {{ total }} 张表情包
            <span v-if="hasActiveFilters"> · {{ filterSummary }}</span>
          </p>
        </div>

        <div class="flex shrink-0 items-center gap-2">
          <span
            v-if="hasActiveFilters && !filtersExpanded"
            class="h-2 w-2 rounded-full bg-[var(--color-accent)]"
            aria-hidden="true"
          />

          <!-- 移动端：合并为一个筛选按钮 -->
          <button
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-[11px] text-zinc-500 transition active:bg-white/[0.06] md:hidden"
            :aria-expanded="filtersExpanded"
            aria-controls="home-filters-panel"
            @click="filtersExpanded = !filtersExpanded"
          >
            <span>{{ filtersExpanded ? '收起' : '筛选' }}</span>
            <svg
              class="h-4 w-4 transition-transform duration-200"
              :class="{ 'rotate-180': filtersExpanded }"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
            </svg>
          </button>

          <!-- Web 端：搜索 / 上传分开折叠 -->
          <button
            type="button"
            class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs text-zinc-500 transition hover:bg-white/[0.06] hover:text-zinc-300 md:inline-flex"
            :aria-expanded="filtersExpanded"
            aria-controls="home-filters-panel"
            @click="filtersExpanded = !filtersExpanded"
          >
            <span>{{ filtersExpanded ? '收起搜索' : '搜索筛选' }}</span>
            <svg
              class="h-3.5 w-3.5 transition-transform duration-200"
              :class="{ 'rotate-180': filtersExpanded }"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
            </svg>
          </button>

          <button
            type="button"
            class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs text-zinc-500 transition hover:bg-white/[0.06] hover:text-zinc-300 md:inline-flex"
            :aria-expanded="uploadExpanded"
            aria-controls="home-upload-panel"
            @click="uploadExpanded = !uploadExpanded"
          >
            <span>{{ uploadExpanded ? '收起上传' : '上传图片' }}</span>
            <svg
              class="h-3.5 w-3.5 transition-transform duration-200"
              :class="{ 'rotate-180': uploadExpanded }"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 搜索 + 标签筛选 -->
    <div
      id="home-filters-panel"
      class="grid transition-[grid-template-rows] duration-200 ease-out"
      :class="filtersExpanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'"
    >
      <div class="overflow-hidden">
        <SearchBar
          v-model:query="query"
          v-model:sort="sort"
          :category-name="categoryName"
          :total="total"
          hide-header
        />
        <TagFilter
          v-model:selected-ids="selectedIds"
          :tags="tags"
          :loading="loadingTags"
          @rename-tag="$emit('rename-tag', $event)"
          @delete-tag="$emit('delete-tag', $event)"
        />
      </div>
    </div>

    <!-- Web 端上传区 -->
    <div
      id="home-upload-panel"
      class="hidden grid transition-[grid-template-rows] duration-200 ease-out md:grid"
      :class="uploadExpanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'"
    >
      <div class="overflow-hidden">
        <div class="px-6 pt-4">
          <slot name="upload" />
        </div>
      </div>
    </div>
  </div>
</template>
