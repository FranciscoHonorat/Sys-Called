import { onMounted, onUnmounted } from 'vue'

export function usePolling(task: () => unknown, everyMs: number) {
  let timer: ReturnType<typeof setInterval> | undefined

  onMounted(() => {
    task()
    timer = setInterval(task, everyMs)
  })

  onUnmounted(() => clearInterval(timer))
}
