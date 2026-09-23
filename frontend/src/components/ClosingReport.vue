<script setup lang="ts">
import { formatDuration } from '../tickets/format'

defineProps<{ resolution?: string; openedAt: string; closedAt: string }>()
</script>

<template>
  <section aria-label="Relatório de fechamento" class="rounded border border-green-200 bg-green-50 p-4 text-sm text-slate-700">
    <h3 class="mb-1 font-semibold text-slate-800">Relatório de fechamento</h3>
    <p class="whitespace-pre-line">{{ resolution }}</p>
    <p class="mt-2 text-slate-500">Tempo até fechar: {{ formatDuration(openedAt, closedAt) }}</p>
  </section>
</template>
