import { onMounted, ref, type Ref } from 'vue'

export function useLoad<T>(
  load: () => Promise<T>,
  initial: T,
): { data: Ref<T>; failed: Ref<boolean>; reload: () => Promise<void> } {
  const data = ref(initial) as Ref<T>
  const failed = ref(false)

  async function reload() {
    try {
      data.value = await load()
      failed.value = false
    } catch {
      failed.value = true
    }
  }

  onMounted(reload)

  return { data, failed, reload }
}
