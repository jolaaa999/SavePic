<script setup>
import { computed, ref } from 'vue'
import { useIsMobile } from '../composables/useBreakpoint'
import { uploadMeme, unwrapUploadResult } from '../api'

const props = defineProps({
  categoryId: { type: Number, default: null },
  disabled: { type: Boolean, default: false },
  categoryName: { type: String, default: '' },
})

const tagsText = defineModel('tagsText', { type: String, default: '' })

const emit = defineEmits(['uploaded', 'error'])

const isMobile = useIsMobile()
const isDragging = ref(false)
const uploading = ref(false)
const fileInput = ref(null)

const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']

const hintPrimary = computed(() => {
  if (uploading.value) return null
  if (props.disabled) return '请先在「分类」中选择分类'
  if (isMobile.value) return null
  return '拖拽图片到此处，或'
})

const hintSecondary = computed(() => {
  if (uploading.value) return ''
  return '支持 JPG · PNG · GIF · WebP，可多选'
})

function isImage(file) {
  return allowedTypes.includes(file.type) || /\.(jpe?g|png|gif|webp)$/i.test(file.name)
}

function parseTags() {
  return tagsText.value
    .split(/[,，]/)
    .map((s) => s.trim())
    .filter(Boolean)
}

async function uploadFiles(files) {
  if (!props.categoryId || props.disabled) return

  const list = [...files].filter(isImage)
  if (!list.length) {
    emit('error', isMobile.value ? '请选择图片文件' : '请拖入 jpg、png、gif、webp 格式的图片')
    return
  }

  const tags = parseTags()
  uploading.value = true
  let duplicateCount = 0
  try {
    for (const file of list) {
      const res = await uploadMeme(props.categoryId, file, tags)
      const data = unwrapUploadResult(res)
      if (res?.duplicate) duplicateCount++
      void data
    }
    emit('uploaded', { count: list.length, duplicateCount })
  } catch (e) {
    emit('error', e.message)
  } finally {
    uploading.value = false
  }
}

function openPicker() {
  if (!props.disabled && !uploading.value) {
    fileInput.value?.click()
  }
}

function onDragOver(e) {
  if (isMobile.value) return
  e.preventDefault()
  if (!props.disabled) isDragging.value = true
}

function onDragLeave() {
  if (isMobile.value) return
  isDragging.value = false
}

function onDrop(e) {
  if (isMobile.value) return
  e.preventDefault()
  isDragging.value = false
  uploadFiles(e.dataTransfer.files)
}

function onFileChange(e) {
  uploadFiles(e.target.files)
  e.target.value = ''
}
</script>

<template>
  <div
    class="relative rounded-xl border border-dashed transition-all duration-200"
    :class="[
      disabled
        ? 'cursor-not-allowed border-white/[0.04] bg-white/[0.01] opacity-50'
        : isMobile
          ? 'cursor-pointer border-white/[0.08] bg-white/[0.02] active:bg-[var(--color-accent-muted)]'
          : isDragging
            ? 'border-[var(--color-accent)]/50 bg-[var(--color-accent-muted)]'
            : 'border-white/[0.08] bg-white/[0.02] hover:border-white/15',
    ]"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
    @click="isMobile ? openPicker() : undefined"
  >
    <div
      class="flex flex-col items-center justify-center gap-2 px-6 py-8"
      :class="{ 'pointer-events-none': !isMobile && (disabled || uploading) }"
    >
      <input
        ref="fileInput"
        type="file"
        class="hidden"
        accept="image/jpeg,image/png,image/gif,image/webp"
        multiple
        :disabled="disabled || uploading"
        @change="onFileChange"
      />

      <!-- 桌面端额外 file input（label 触发） -->
      <label
        v-if="!isMobile"
        class="flex w-full cursor-pointer flex-col items-center gap-2"
        :class="{ 'pointer-events-none': disabled || uploading }"
      >
        <input
          type="file"
          class="hidden"
          accept="image/jpeg,image/png,image/gif,image/webp"
          multiple
          :disabled="disabled || uploading"
          @change="onFileChange"
        />

        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white/[0.04] text-zinc-500">
          <svg
            v-if="!uploading"
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
          <svg v-else class="h-5 w-5 animate-spin text-[var(--color-accent)]" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>

        <p class="text-center text-sm text-zinc-400">
          <template v-if="uploading">上传中…</template>
          <template v-else-if="disabled">{{ hintPrimary }}</template>
          <template v-else>
            {{ hintPrimary }} <span class="text-[var(--color-accent)]">点击上传</span>
          </template>
        </p>
        <p v-if="hintSecondary" class="text-[11px] text-zinc-600">{{ hintSecondary }}</p>
      </label>

      <!-- 移动端：点击整块区域打开相册 -->
      <template v-else>
        <div class="flex h-12 w-12 items-center justify-center rounded-full bg-white/[0.04] text-zinc-500">
          <svg
            v-if="!uploading"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909M3.75 21h16.5A2.25 2.25 0 0 0 22.5 18.75V5.25A2.25 2.25 0 0 0 20.25 3H3.75A2.25 2.25 0 0 0 1.5 5.25v13.5A2.25 2.25 0 0 0 3.75 21Z"
            />
          </svg>
          <svg v-else class="h-6 w-6 animate-spin text-[var(--color-accent)]" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>

        <p v-if="categoryName && !disabled" class="text-xs text-zinc-600">
          当前分类：{{ categoryName }}
        </p>
        <p class="text-center text-sm text-zinc-400">
          <template v-if="uploading">上传中…</template>
          <template v-else-if="disabled">{{ hintPrimary }}</template>
          <template v-else>
            点击 <span class="text-[var(--color-accent)]">从相册选择</span>（可多选）
          </template>
        </p>
        <p v-if="hintSecondary && !uploading" class="text-[11px] text-zinc-600">{{ hintSecondary }}</p>
      </template>

      <!-- 标签输入 -->
      <div
        v-if="!disabled && !uploading"
        class="mt-3 w-full max-w-sm"
        @click.stop
      >
        <input
          v-model="tagsText"
          type="text"
          placeholder="标签：破防, 打工人, 周一（逗号分隔）"
          class="w-full rounded-lg border border-white/[0.06] bg-white/[0.03] px-3 py-2 text-xs text-zinc-400 placeholder-zinc-600 outline-none focus:border-[var(--color-accent)]/30"
        />
      </div>
    </div>
  </div>
</template>
