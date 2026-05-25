<script setup>
defineProps({
  categories: { type: Array, required: true },
  selectedId: { type: Number, default: null },
  loading: { type: Boolean, default: false },
})

defineEmits(['select', 'new-category'])
</script>

<template>
  <aside
    class="relative z-10 hidden w-60 shrink-0 flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)] md:flex"
  >
    <!-- Logo -->
    <div class="flex items-center gap-2.5 px-5 py-6">
      <div
        class="flex h-8 w-8 items-center justify-center rounded-lg bg-[var(--color-accent-muted)] text-sm font-semibold text-[var(--color-accent)]"
      >
        S
      </div>
      <div>
        <p class="text-sm font-medium tracking-tight text-zinc-100">SavePic</p>
        <p class="text-[11px] text-zinc-600">表情包库</p>
      </div>
    </div>

    <p class="px-5 pb-2 text-[10px] font-medium uppercase tracking-widest text-zinc-600">
      分类
    </p>

    <!-- 分类列表 -->
    <nav class="flex-1 space-y-0.5 overflow-y-auto px-3">
      <p v-if="loading" class="px-3 py-4 text-center text-xs text-zinc-600">加载中…</p>
      <p v-else-if="!categories.length" class="px-3 py-4 text-center text-xs text-zinc-600">
        暂无分类，点击下方新建
      </p>
      <button
        v-for="cat in categories"
        :key="cat.id"
        type="button"
        class="group flex w-full items-center justify-between rounded-lg px-3 py-2.5 text-left transition-all duration-150"
        :class="
          selectedId === cat.id
            ? 'bg-[var(--color-accent-muted)] text-zinc-100'
            : 'text-zinc-500 hover:bg-white/[0.03] hover:text-zinc-300'
        "
        @click="$emit('select', cat.id)"
      >
        <span class="truncate text-sm font-medium">{{ cat.name }}</span>
        <span
          class="ml-2 shrink-0 rounded-md px-1.5 py-0.5 text-[10px] tabular-nums transition-colors"
          :class="
            selectedId === cat.id
              ? 'bg-[var(--color-accent)]/20 text-[var(--color-accent)]'
              : 'bg-white/5 text-zinc-600 group-hover:text-zinc-500'
          "
        >
          {{ cat.count }}
        </span>
      </button>
    </nav>

    <!-- 新建分类 -->
    <div class="border-t border-[var(--color-border)] p-3">
      <button
        type="button"
        class="flex w-full items-center justify-center gap-2 rounded-lg border border-dashed border-white/10 px-3 py-2.5 text-sm text-zinc-500 transition-all hover:border-[var(--color-accent)]/40 hover:bg-[var(--color-accent-muted)] hover:text-[var(--color-accent)]"
        @click="$emit('new-category')"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        新建分类
      </button>
    </div>
  </aside>
</template>
