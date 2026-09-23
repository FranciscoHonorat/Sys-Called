<script setup lang="ts">
import { inject } from 'vue'

import PageLayout from '../components/PageLayout.vue'
import { useLoad } from '../composables/useLoad'
import { ticketsApiKey, type AgentWorkload } from '../tickets/ticketsApi'

const api = inject(ticketsApiKey)!

const { data: workloads, failed } = useLoad<AgentWorkload[]>(() => api.supportWorkload(), [])
</script>

<template>
  <PageLayout back>
    <section class="rounded-lg bg-white p-6 shadow">
      <h2 class="mb-4 text-xl font-semibold text-slate-800">Suportes</h2>
      <p v-if="failed" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">Não foi possível carregar os suportes</p>
      <table v-else class="w-full text-left text-sm">
        <thead class="text-slate-500">
          <tr>
            <th class="py-2">Nome</th>
            <th>Abertos</th>
            <th>Em andamento</th>
            <th>Fechados</th>
            <th>Total</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="agent in workloads" :key="agent.id" class="border-t border-slate-100">
            <td class="py-2 font-medium text-slate-800">{{ agent.name }}</td>
            <td>{{ agent.open }}</td>
            <td>{{ agent.in_progress }}</td>
            <td>{{ agent.closed }}</td>
            <td>{{ agent.open + agent.in_progress + agent.closed }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </PageLayout>
</template>
