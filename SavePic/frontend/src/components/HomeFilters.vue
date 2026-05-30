<script setup>
import { computed, ref } from 'vue'
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
const expanded = ref(false)

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
  <div class="shrink-0 md:contents">
    <!-- 移动端：折叠条 -->
    <div
      class="sticky top-0 z-10 border-b border-[var(--color-border)] bg-[var(--color-surface)]/90 backdrop-blur-md md:hidden"
    >
      <button
        type="button"
        class="flex w-full items-center gap-3 px-4 py-3 text-left active:bg-white/[0.03]"
        :aria-expanded="expanded"
        aria-controls="mobile-home-filters"
        @click="expanded = !expanded"
      >
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
            v-if="hasActiveFilters && !expanded"
            class="h-2 w-2 rounded-full bg-[var(--color-accent)]"
            aria-hidden="true"
          />
          <span class="text-[11px] text-zinc-500">{{ expanded ? '收起' : '筛选' }}</span>
          <svg
            class="h-4 w-4 text-zinc-500 transition-transform duration-200"
            :class="{ 'rotate-180': expanded }"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
            aria-hidden="true"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
          </svg>
        </div>
      </button>
    </div>

    <!-- 移动端可折叠 / 桌面端始终展示 -->
    <div
      id="mobile-home-filters"
      class="grid transition-[grid-template-rows] duration-200 ease-out md:contents"
      :class="expanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr] md:grid-rows-[1fr]'"
    >
      <div class="overflow-hidden md:contents">
        <SearchBar
          v-model:query="query"
          v-model:sort="sort"
          :category-name="categoryName"
          :total="total"
          :hide-header="isMobile"
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
  </div>
</template>
