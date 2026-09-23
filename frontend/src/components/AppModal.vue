<script setup lang="ts">
import { nextTick, ref, useId, watch } from 'vue'

const props = defineProps<{ open: boolean; title: string }>()
const emit = defineEmits<{ close: [] }>()

const titleId = useId()
const panel = ref<HTMLElement | null>(null)

watch(
  () => props.open,
  async (open) => {
    if (open) {
      await nextTick()
      panel.value?.focus()
    }
  },
  { immediate: true },
)
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        data-testid="modal-backdrop"
        class="fixed inset-0 z-50 flex items-end justify-center bg-slate-900/50 p-4 sm:items-center"
        @click.self="emit('close')"
      >
        <section
          ref="panel"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="titleId"
          tabindex="-1"
          class="modal-panel w-full max-w-lg rounded-lg bg-white p-6 shadow-xl focus:outline-none"
          @keydown.esc="emit('close')"
        >
          <header class="mb-4 flex items-center justify-between gap-4">
            <button type="button" class="text-sm text-blue-700 hover:underline" @click="emit('close')">← Voltar</button>
            <h2 :id="titleId" class="text-lg font-semibold text-slate-800">{{ title }}</h2>
            <span class="w-14" aria-hidden="true" />
          </header>
          <slot />
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
@keyframes slide-up {
  from {
    transform: translateY(2rem);
  }
  to {
    transform: translateY(0);
  }
}

.modal-panel {
  animation: slide-up 200ms ease-out;
}
</style>
