<script setup lang="ts">
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'

import AppHeader from '../components/AppHeader.vue'
import FormField from '../components/FormField.vue'
import SelectField from '../components/SelectField.vue'
import { paths } from '../router/paths'
import { priorityOptions } from '../tickets/priority'
import { ticketsApiKey } from '../tickets/ticketsApi'

const api = inject(ticketsApiKey)!
const router = useRouter()

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
    await api.open({ title: title.value, description: description.value, priority: priority.value })
    await router.push(paths.tickets)
  } catch {
    submitError.value = 'Não foi possível abrir o chamado'
  }
}
</script>

<template>
  <main class="min-h-screen bg-slate-100 p-6">
    <AppHeader />

    <section class="max-w-xl rounded-lg bg-white p-6 shadow">
      <h2 class="mb-4 text-xl font-semibold text-slate-800">Abrir novo chamado</h2>

      <form class="space-y-4" novalidate @submit.prevent="submit">
        <p v-if="submitError" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ submitError }}</p>

        <FormField id="title" v-model="title" label="Título" :error="titleError" />

        <FormField id="description" v-model="description" label="Descrição" multiline :error="descriptionError" />
        <SelectField id="priority" v-model="priority" label="Prioridade" :options="priorityOptions" />

        <button type="submit" class="rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700">
          Abrir chamado
        </button>
      </form>
    </section>
  </main>
</template>
