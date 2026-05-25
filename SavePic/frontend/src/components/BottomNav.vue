<script setup>
const model = defineModel({ type: String, default: 'home' })

const items = [
  { id: 'home', label: '首页', icon: 'home' },
  { id: 'categories', label: '分类', icon: 'categories' },
  { id: 'upload', label: '上传', icon: 'upload' },
]
</script>

<template>
  <nav
    class="fixed inset-x-0 bottom-0 z-40 border-t border-[var(--color-border)] bg-[var(--color-surface)]/95 backdrop-blur-lg md:hidden"
    style="padding-bottom: env(safe-area-inset-bottom, 0px)"
  >
    <div class="flex h-14 items-stretch">
      <button
        v-for="item in items"
        :key="item.id"
        type="button"
        class="flex flex-1 flex-col items-center justify-center gap-0.5 transition-colors"
        :class="
          model === item.id
            ? 'text-[var(--color-accent)]'
            : 'text-zinc-500 active:text-zinc-300'
        "
        @click="model = item.id"
      >
        <!-- 首页 -->
        <svg
          v-if="item.icon === 'home'"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
          />
        </svg>
        <!-- 分类 -->
        <svg
          v-else-if="item.icon === 'categories'"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"
          />
        </svg>
        <!-- 上传 -->
        <svg
          v-else
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5"
          />
        </svg>
        <span class="text-[10px] font-medium">{{ item.label }}</span>
      </button>
    </div>
  </nav>
</template>
