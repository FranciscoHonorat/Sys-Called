<script setup lang="ts">
import { inject, ref } from 'vue'

import { priorityOptions } from '../tickets/priority'
import { ticketsApiKey } from '../tickets/ticketsApi'
import FormField from './FormField.vue'
import SelectField from './SelectField.vue'

const emit = defineEmits<{ opened: [ticketId: string] }>()

const api = inject(ticketsApiKey)!

const title = ref('')
const description = ref('')
const priority = ref('Medium')
const titleError = ref('')
const descriptionError = ref('')
const submitError = ref('')

async function submit() {
  titleError.value = title.value ? '' : 'Informe o título'
  descriptionError.value = description.value ? '' : 'Informe a descrição'
  if (titleError.value || descriptionError.value) {
    return
  }

  try {
    emit('opened', await api.open({ title: title.value, description: description.value, priority: priority.value }))
  } catch {
    submitError.value = 'Não foi possível abrir o chamado'
  }
}
</script>

<template>
  <form class="space-y-4" novalidate @submit.prevent="submit">
    <p v-if="submitError" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ submitError }}</p>

    <FormField id="title" v-model="title" label="Título" :error="titleError" />
    <FormField id="description" v-model="description" label="Descrição" multiline :error="descriptionError" />
    <SelectField id="priority" v-model="priority" label="Prioridade" :options="priorityOptions" />

    <button type="submit" class="w-full rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700">
      Abrir chamado
    </button>
  </form>
</template>
