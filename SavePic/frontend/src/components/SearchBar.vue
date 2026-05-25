<script setup>
defineProps({
  categoryName: { type: String, default: '' },
  total: { type: Number, default: 0 },
})

const query = defineModel('query', { type: String, default: '' })
const sort = defineModel('sort', { type: String, default: 'desc' })
</script>

<template>
  <header
    class="sticky top-0 z-10 flex flex-wrap items-center gap-2 border-b border-[var(--color-border)] bg-[var(--color-surface)]/80 px-4 py-3 backdrop-blur-md md:gap-3 md:px-6 md:py-4"
  >
    <div class="min-w-0 flex-1">
      <h1 class="truncate text-base font-medium text-zinc-100">
        {{ categoryName || '选择分类' }}
      </h1>
      <p class="mt-0.5 text-xs text-zinc-600">{{ total }} 张表情包</p>
    </div>

    <!-- 搜索 -->
    <div class="relative w-full sm:w-56">
      <svg
        class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-600"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="1.5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
        />
      </svg>
      <input
        v-model="query"
        type="search"
        placeholder="搜索…"
        class="w-full rounded-lg border border-white/[0.06] bg-white/[0.03] py-2 pl-9 pr-3 text-sm text-zinc-300 placeholder-zinc-600 outline-none transition focus:border-[var(--color-accent)]/30 focus:ring-1 focus:ring-[var(--color-accent)]/20"
      />
    </div>

    <!-- 排序 -->
    <select
      v-model="sort"
      class="cursor-pointer rounded-lg border border-white/[0.06] bg-white/[0.03] px-3 py-2 text-sm text-zinc-400 outline-none transition hover:border-white/10 focus:border-[var(--color-accent)]/30"
    >
      <option value="desc">最新优先</option>
      <option value="asc">最早优先</option>
    </select>
  </header>
</template>
