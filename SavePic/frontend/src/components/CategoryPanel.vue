<script setup>
defineProps({
  categories: { type: Array, required: true },
  selectedId: { type: Number, default: null },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['select', 'new-category'])

function onSelect(id) {
  emit('select', id)
}
</script>

<template>
  <div class="flex flex-1 flex-col overflow-hidden">
    <div class="border-b border-[var(--color-border)] px-4 py-4">
      <h2 class="text-base font-medium text-zinc-100">选择分类</h2>
      <p class="mt-0.5 text-xs text-zinc-600">点击分类后返回首页查看表情包</p>
    </div>

    <nav class="flex-1 space-y-0.5 overflow-y-auto px-3 py-3">
      <p v-if="loading" class="py-8 text-center text-xs text-zinc-600">加载中…</p>
      <p v-else-if="!categories.length" class="py-8 text-center text-xs text-zinc-600">
        暂无分类
      </p>
      <button
        v-for="cat in categories"
        :key="cat.id"
        type="button"
        class="group flex w-full items-center justify-between rounded-xl px-4 py-3.5 text-left transition-all active:scale-[0.98]"
        :class="
          selectedId === cat.id
            ? 'bg-[var(--color-accent-muted)] text-zinc-100'
            : 'text-zinc-500 active:bg-white/[0.04]'
        "
        @click="onSelect(cat.id)"
      >
        <span class="text-sm font-medium">{{ cat.name }}</span>
        <span
          class="rounded-md px-2 py-0.5 text-[10px] tabular-nums"
          :class="
            selectedId === cat.id
              ? 'bg-[var(--color-accent)]/20 text-[var(--color-accent)]'
              : 'bg-white/5 text-zinc-600'
          "
        >
          {{ cat.count }}
        </span>
      </button>
    </nav>

    <div class="border-t border-[var(--color-border)] p-3">
      <button
        type="button"
        class="flex w-full items-center justify-center gap-2 rounded-xl border border-dashed border-white/10 py-3 text-sm text-zinc-500 active:bg-[var(--color-accent-muted)] active:text-[var(--color-accent)]"
        @click="$emit('new-category')"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        新建分类
      </button>
    </div>
  </div>
</template>
