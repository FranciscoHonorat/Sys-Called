<script setup lang="ts">
import { computed, ref } from 'vue'

import { formatDuration } from '../tickets/format'
import AppModal from './AppModal.vue'
import FormField from './FormField.vue'

const props = defineProps<{ open: boolean; openedAt: string }>()
const emit = defineEmits<{ close: []; confirm: [resolution: string] }>()

const resolution = ref('')
const error = ref('')
const elapsed = computed(() => formatDuration(props.openedAt, new Date().toISOString()))

function confirm() {
  if (!resolution.value.trim()) {
    error.value = 'Descreva o que foi feito'
    return
  }
  emit('confirm', resolution.value.trim())
}
</script>

<template>
  <AppModal :open="open" title="Fechar chamado" @close="emit('close')">
    <form class="space-y-4" novalidate @submit.prevent="confirm">
      <p class="rounded bg-slate-50 p-3 text-sm text-slate-600">Tempo de atendimento até agora: {{ elapsed }}</p>
      <FormField id="resolution" v-model="resolution" label="O que foi feito" multiline :error="error" />
      <button type="submit" class="w-full rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700">
        Confirmar fechamento
      </button>
    </form>
  </AppModal>
</template>
