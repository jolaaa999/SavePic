import { onMounted, onUnmounted, ref } from 'vue'

/** 屏幕宽度小于 Tailwind md 断点 (768px) */
export function useIsMobile() {
  const isMobile = ref(false)
  let mql

  function update() {
    isMobile.value = mql?.matches ?? false
  }

  onMounted(() => {
    mql = window.matchMedia('(max-width: 767px)')
    update()
    mql.addEventListener('change', update)
  })

  onUnmounted(() => {
    mql?.removeEventListener('change', update)
  })

  return isMobile
}
