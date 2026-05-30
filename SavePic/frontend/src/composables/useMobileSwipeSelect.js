import { computed, ref } from 'vue'

const LONG_PRESS_MS = 450
const MOVE_THRESHOLD = 10

function createGestureState() {
  return {
    startX: 0,
    startY: 0,
    moved: false,
    longPressTimer: null,
    longPressFired: false,
    strokeTouched: new Set(),
    strokeAdd: true,
    suppressClick: false,
    activeMemeId: null,
  }
}

/**
 * 移动端长按进入多选，滑动批量选中/取消
 * @param {{ enabled?: import('vue').Ref<boolean> | import('vue').ComputedRef<boolean> }} options
 */
export function useMobileSwipeSelect(options = {}) {
  const selectionMode = ref(false)
  const selectedSet = ref(new Set())

  const selectedIds = computed(() => [...selectedSet.value])
  const selectedCount = computed(() => selectedSet.value.size)

  let gesture = createGestureState()
  let docTouchMove = null
  let docTouchEnd = null

  function isEnabled() {
    return options.enabled?.value !== false
  }

  function isSelected(id) {
    return selectedSet.value.has(id)
  }

  function updateSet(mutator) {
    const next = new Set(selectedSet.value)
    mutator(next)
    selectedSet.value = next
  }

  function select(id) {
    updateSet((set) => set.add(id))
  }

  function deselect(id) {
    updateSet((set) => set.delete(id))
  }

  function toggle(id) {
    if (isSelected(id)) deselect(id)
    else select(id)
  }

  function enterSelection(id) {
    selectionMode.value = true
    select(id)
    if (navigator.vibrate) navigator.vibrate(12)
  }

  function exitSelection() {
    selectionMode.value = false
    selectedSet.value = new Set()
  }

  function getMemeIdFromPoint(x, y) {
    const el = document.elementFromPoint(x, y)
    const card = el?.closest('[data-meme-id]')
    if (!card) return null
    const id = Number(card.dataset.memeId)
    return Number.isFinite(id) ? id : null
  }

  function applyStroke(id) {
    if (gesture.strokeTouched.has(id)) return
    gesture.strokeTouched.add(id)
    if (gesture.strokeAdd) select(id)
    else deselect(id)
    if (selectedSet.value.size === 0) exitSelection()
  }

  function clearDocListeners() {
    if (docTouchMove) {
      document.removeEventListener('touchmove', docTouchMove)
      docTouchMove = null
    }
    if (docTouchEnd) {
      document.removeEventListener('touchend', docTouchEnd)
      document.removeEventListener('touchcancel', docTouchEnd)
      docTouchEnd = null
    }
  }

  function finishGesture() {
    if (gesture.longPressTimer) {
      clearTimeout(gesture.longPressTimer)
      gesture.longPressTimer = null
    }
    clearDocListeners()
    gesture.strokeTouched = new Set()
  }

  function bindDocListeners() {
    clearDocListeners()

    docTouchMove = (e) => {
      if (!isEnabled()) return
      const touch = e.touches[0]
      if (!touch) return

      const dx = touch.clientX - gesture.startX
      const dy = touch.clientY - gesture.startY
      if (Math.hypot(dx, dy) > MOVE_THRESHOLD) {
        gesture.moved = true
        if (gesture.longPressTimer) {
          clearTimeout(gesture.longPressTimer)
          gesture.longPressTimer = null
        }
      }

      if (selectionMode.value || gesture.longPressFired) {
        e.preventDefault()
        const id = getMemeIdFromPoint(touch.clientX, touch.clientY)
        if (id != null) applyStroke(id)
      }
    }

    docTouchEnd = () => {
      if (
        selectionMode.value &&
        !gesture.moved &&
        !gesture.longPressFired &&
        gesture.activeMemeId != null &&
        gesture.strokeTouched.size === 0
      ) {
        toggle(gesture.activeMemeId)
        if (selectedSet.value.size === 0) exitSelection()
      }
      finishGesture()
    }

    document.addEventListener('touchmove', docTouchMove, { passive: false })
    document.addEventListener('touchend', docTouchEnd, { passive: true })
    document.addEventListener('touchcancel', docTouchEnd, { passive: true })
  }

  /** @param {{ id: number }} meme @param {TouchEvent} e */
  function onMemeTouchStart(meme, e) {
    if (!isEnabled()) return

    finishGesture()
    const touch = e.touches[0]
    if (!touch) return

    gesture = createGestureState()
    gesture.startX = touch.clientX
    gesture.startY = touch.clientY
    gesture.activeMemeId = meme.id

    if (selectionMode.value) {
      gesture.strokeAdd = !isSelected(meme.id)
      bindDocListeners()
      return
    }

    gesture.longPressTimer = window.setTimeout(() => {
      gesture.longPressFired = true
      gesture.suppressClick = true
      gesture.strokeAdd = true
      gesture.strokeTouched.add(meme.id)
      enterSelection(meme.id)
    }, LONG_PRESS_MS)

    bindDocListeners()
  }

  function shouldOpenPreview() {
    if (!isEnabled()) return true
    if (gesture.suppressClick) {
      gesture.suppressClick = false
      return false
    }
    if (selectionMode.value) return false
    return true
  }

  function cleanup() {
    finishGesture()
  }

  return {
    selectionMode,
    selectedIds,
    selectedCount,
    isSelected,
    exitSelection,
    onMemeTouchStart,
    shouldOpenPreview,
    cleanup,
  }
}
