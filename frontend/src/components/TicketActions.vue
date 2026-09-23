<script setup lang="ts">
import { ref } from 'vue'

import type { TicketAction } from '../tickets/permissions'
import { priorityOptions } from '../tickets/priority'
import SelectField from './SelectField.vue'

const props = defineProps<{
  actions: TicketAction[]
  responsibles: ReadonlyArray<{ value: string; label: string }>
}>()

const emit = defineEmits<{
  edit: []
  assign: [assigneeId: string]
  autoAssign: []
  changePriority: [priority: string]
  start: []
  close: []
}>()

const selectedAssignee = ref('')
const selectedPriority = ref('Medium')

const can = (action: TicketAction) => props.actions.includes(action)
</script>

<template>
  <section aria-label="Ações" class="space-y-3 border-t border-slate-100 pt-4">
    <button v-if="can('edit')" type="button" class="action-button" @click="emit('edit')">Editar</button>

    <div v-if="can('manage')" class="flex flex-wrap items-end gap-2">
      <SelectField id="assignee" v-model="selectedAssignee" label="Atribuir a" :options="responsibles" />
      <button type="button" class="action-button" @click="emit('assign', selectedAssignee)">Atribuir</button>
      <button type="button" class="action-button" @click="emit('autoAssign')">Distribuir automaticamente</button>
      <SelectField id="new-priority" v-model="selectedPriority" label="Nova prioridade" :options="priorityOptions" />
      <button type="button" class="action-button" @click="emit('changePriority', selectedPriority)">Alterar prioridade</button>
    </div>

    <div class="flex gap-2">
      <button v-if="can('start')" type="button" class="action-button" @click="emit('start')">Iniciar atendimento</button>
      <button v-if="can('close')" type="button" class="action-button" @click="emit('close')">Fechar chamado</button>
    </div>
  </section>
</template>

<style scoped>
@reference 'tailwindcss';

.action-button {
  @apply rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50;
}
</style>
