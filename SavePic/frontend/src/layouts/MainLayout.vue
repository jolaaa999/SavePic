<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import {
  createCategory,
  deleteCategory,
  deleteMeme,
  fetchCategories,
  fetchMemes,
  fetchTags,
} from '../api'
import { useIsMobile } from '../composables/useBreakpoint'
import Sidebar from '../components/Sidebar.vue'
import SearchBar from '../components/SearchBar.vue'
import TagFilter from '../components/TagFilter.vue'
import UploadZone from '../components/UploadZone.vue'
import MemeGrid from '../components/MemeGrid.vue'
import BottomNav from '../components/BottomNav.vue'
import CategoryPanel from '../components/CategoryPanel.vue'

const isMobile = useIsMobile()
const mobileTab = ref('home')

const categoryList = ref([])
const tagList = ref([])
const memeList = ref([])
const selectedCategoryId = ref(null)
const selectedTagIds = ref([])
const searchQuery = ref('')
const sortOrder = ref('desc')
const uploadTagsText = ref('')

const loadingCategories = ref(true)
const loadingTags = ref(true)
const loadingMemes = ref(false)
const globalError = ref('')
const globalMessage = ref('')

const selectedCategory = computed(() =>
  categoryList.value.find((c) => c.id === selectedCategoryId.value),
)

const filteredMemes = computed(() => {
  let list = [...memeList.value]
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return list
  return list.filter((m) => {
    if (String(m.id).includes(q)) return true
    return (m.tags || []).some((t) => t.name.toLowerCase().includes(q))
  })
})

async function loadCategories() {
  loadingCategories.value = true
  try {
    const data = await fetchCategories()
    categoryList.value = data
    if (!selectedCategoryId.value && data.length) {
      selectedCategoryId.value = data[0].id
    }
  } catch (e) {
    globalError.value = e.message
  } finally {
    loadingCategories.value = false
  }
}

async function loadTags() {
  loadingTags.value = true
  try {
    tagList.value = await fetchTags()
  } catch (e) {
    globalError.value = e.message
  } finally {
    loadingTags.value = false
  }
}

async function loadMemes() {
  if (!selectedCategoryId.value) {
    memeList.value = []
    return
  }
  loadingMemes.value = true
  try {
    memeList.value = await fetchMemes({
      categoryId: selectedCategoryId.value,
      tagIds: selectedTagIds.value,
      sort: sortOrder.value,
    })
  } catch (e) {
    globalError.value = e.message
    memeList.value = []
  } finally {
    loadingMemes.value = false
  }
}

function selectCategory(id) {
  selectedCategoryId.value = id
  if (isMobile.value) mobileTab.value = 'home'
}

async function onNewCategory() {
  const name = window.prompt('请输入分类名称')
  if (!name?.trim()) return
  try {
    const cat = await createCategory({
      name: name.trim(),
      sort_order: categoryList.value.length,
    })
    await loadCategories()
    selectedCategoryId.value = cat.id
    globalMessage.value = '分类创建成功'
    setTimeout(() => (globalMessage.value = ''), 2000)
  } catch (e) {
    globalError.value = e.message
  }
}

/**
 * @param {{ id: number, name: string, count?: number }} cat
 */
async function onDeleteCategory(cat) {
  const count = cat.count ?? 0
  const hint =
    count > 0
      ? `确定删除分类「${cat.name}」吗？其下 ${count} 张表情包将一并删除，且无法恢复。`
      : `确定删除分类「${cat.name}」吗？`
  if (!confirm(hint)) return

  const wasSelected = selectedCategoryId.value === cat.id
  try {
    await deleteCategory(cat.id)
    await Promise.all([loadCategories(), loadTags()])
    if (wasSelected) {
      selectedCategoryId.value = categoryList.value[0]?.id ?? null
      if (selectedCategoryId.value) await loadMemes()
      else memeList.value = []
    }
    globalMessage.value = '分类已删除'
    setTimeout(() => (globalMessage.value = ''), 2000)
  } catch (e) {
    globalError.value = e.message
    setTimeout(() => (globalError.value = ''), 3000)
  }
}

async function onUploaded({ count, duplicateCount }) {
  await Promise.all([loadMemes(), loadCategories(), loadTags()])
  if (duplicateCount > 0) {
    globalMessage.value = `上传完成，${duplicateCount} 张已存在（未重复保存）`
  } else {
    globalMessage.value = `成功上传 ${count} 张`
  }
  setTimeout(() => (globalMessage.value = ''), 2500)
  if (isMobile.value) mobileTab.value = 'home'
}

function onUploadError(msg) {
  globalError.value = msg
  setTimeout(() => (globalError.value = ''), 3000)
}

async function onDeleteMeme(id) {
  try {
    await deleteMeme(id)
    await Promise.all([loadMemes(), loadCategories(), loadTags()])
    globalMessage.value = '已删除'
    setTimeout(() => (globalMessage.value = ''), 2000)
  } catch (e) {
    globalError.value = e.message
    setTimeout(() => (globalError.value = ''), 3000)
  }
}

onMounted(() => {
  loadCategories()
  loadTags()
})

watch([selectedCategoryId, selectedTagIds, sortOrder], () => {
  if (selectedCategoryId.value) loadMemes()
}, { deep: true, immediate: true })
</script>

<template>
  <div class="flex h-full overflow-hidden bg-[var(--color-surface)]">
    <div
      class="pointer-events-none fixed inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(110,231,183,0.08),transparent)]"
      aria-hidden="true"
    />

    <Transition name="fade">
      <div
        v-if="globalMessage"
        class="fixed top-4 left-1/2 z-50 -translate-x-1/2 rounded-lg bg-[var(--color-accent-muted)] px-4 py-2 text-sm text-[var(--color-accent)] ring-1 ring-[var(--color-accent)]/20"
      >
        {{ globalMessage }}
      </div>
    </Transition>
    <Transition name="fade">
      <div
        v-if="globalError"
        class="fixed top-4 left-1/2 z-50 -translate-x-1/2 rounded-lg bg-red-500/10 px-4 py-2 text-sm text-red-400 ring-1 ring-red-500/20"
      >
        {{ globalError }}
      </div>
    </Transition>

    <Sidebar
      :categories="categoryList"
      :selected-id="selectedCategoryId"
      :loading="loadingCategories"
      @select="selectCategory"
      @new-category="onNewCategory"
      @delete-category="onDeleteCategory"
    />

    <main
      class="relative flex min-w-0 flex-1 flex-col"
      :class="{ 'pb-[calc(3.5rem+env(safe-area-inset-bottom))]': isMobile }"
    >
      <template v-if="!isMobile || mobileTab === 'home'">
        <SearchBar
          v-model:query="searchQuery"
          v-model:sort="sortOrder"
          :category-name="selectedCategory?.name"
          :total="filteredMemes.length"
        />

        <TagFilter
          v-model:selected-ids="selectedTagIds"
          :tags="tagList"
          :loading="loadingTags"
        />

        <div class="hidden shrink-0 px-6 pt-4 md:block">
          <UploadZone
            v-model:tags-text="uploadTagsText"
            :category-id="selectedCategoryId"
            :disabled="!selectedCategoryId"
            :category-name="selectedCategory?.name"
            @uploaded="onUploaded"
            @error="onUploadError"
          />
        </div>

        <MemeGrid
          :memes="filteredMemes"
          :loading="loadingMemes"
          @delete="onDeleteMeme"
          @updated="loadMemes(); loadTags()"
        />
      </template>

      <CategoryPanel
        v-else-if="mobileTab === 'categories'"
        :categories="categoryList"
        :selected-id="selectedCategoryId"
        :loading="loadingCategories"
        @select="selectCategory"
        @new-category="onNewCategory"
        @delete-category="onDeleteCategory"
      />

      <div v-else-if="mobileTab === 'upload'" class="flex flex-1 flex-col overflow-y-auto px-4 py-4">
        <div class="mb-4">
          <h2 class="text-base font-medium text-zinc-100">上传表情包</h2>
          <p class="mt-0.5 text-xs text-zinc-600">支持添加标签，重复图片将自动跳过</p>
        </div>
        <UploadZone
          v-model:tags-text="uploadTagsText"
          :category-id="selectedCategoryId"
          :disabled="!selectedCategoryId"
          :category-name="selectedCategory?.name"
          @uploaded="onUploaded"
          @error="onUploadError"
        />
      </div>
    </main>

    <BottomNav v-if="isMobile" v-model="mobileTab" />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
