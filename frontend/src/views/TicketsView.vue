<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import { useSessionExit } from '../auth/useSessionExit'
import AppHeader from '../components/AppHeader.vue'
import { paths, ticketDetailPath } from '../router/paths'
import { formatDate } from '../tickets/format'
import { priorityLabel } from '../tickets/priority'
import { nameResolver } from '../tickets/responsibles'
import { statusLabel, TicketStatus } from '../tickets/status'
import { SessionExpiredError, ticketsApiKey, type Responsible, type Ticket } from '../tickets/ticketsApi'

const api = inject(ticketsApiKey)!
const { expire } = useSessionExit()
const route = useRoute()

const tickets = ref<Ticket[]>([])
const responsibles = ref<Responsible[]>([])
const loaded = ref(false)

const filters: Array<{ label: string; status?: TicketStatus }> = [
  { label: 'Todos' },
  { label: 'Abertos', status: TicketStatus.Open },
  { label: 'Em atendimento', status: TicketStatus.InProgress },
  { label: 'Fechados', status: TicketStatus.Closed },
]

const nameOf = computed(() => nameResolver(responsibles.value))

const currentStatus = computed(() => route.query.status as TicketStatus | undefined)
const visible = computed(() =>
  currentStatus.value ? tickets.value.filter((t) => t.status === currentStatus.value) : tickets.value,
)

onMounted(async () => {
  try {
    const [list, people] = await Promise.all([api.list(), api.responsibles()])
    tickets.value = list
    responsibles.value = people
    loaded.value = true
  } catch (error) {
    if (error instanceof SessionExpiredError) {
      await expire()
    }
  }
})
</script>

<template>
  <main class="min-h-screen bg-slate-100 p-6">
    <AppHeader />

    <section class="rounded-lg bg-white p-6 shadow">
      <h2 class="mb-4 text-xl font-semibold text-slate-800">Chamados</h2>

      <nav class="mb-4 flex gap-2 text-sm">
        <RouterLink
          v-for="filter in filters"
          :key="filter.label"
          v-slot="{ href, navigate }"
          :to="{ path: paths.tickets, query: filter.status ? { status: filter.status } : {} }"
          custom
        >
          <a
            :href="href"
            :aria-current="currentStatus === filter.status ? 'page' : undefined"
            class="rounded px-3 py-1"
            :class="currentStatus === filter.status ? 'bg-blue-600 text-white' : 'text-slate-600 hover:bg-slate-100'"
            @click="navigate"
          >
            {{ filter.label }}
          </a>
        </RouterLink>
      </nav>

      <p v-if="loaded && visible.length === 0" class="text-slate-500">Nenhum chamado encontrado</p>

      <table v-else-if="visible.length > 0" class="w-full text-left text-sm">
        <thead class="text-slate-500">
          <tr>
            <th class="py-2">Título</th>
            <th>Status</th>
            <th>Prioridade</th>
            <th>Responsável</th>
            <th>Aberto em</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ticket in visible" :key="ticket.ticket_id" class="border-t border-slate-100">
            <td class="py-2 font-medium">
              <RouterLink :to="ticketDetailPath(ticket.ticket_id)" class="text-blue-700 hover:underline">{{ ticket.title }}</RouterLink>
            </td>
            <td>{{ statusLabel(ticket.status) }}</td>
            <td>{{ priorityLabel(ticket.priority) }}</td>
            <td>{{ nameOf(ticket.assignee_id) }}</td>
            <td>{{ formatDate(ticket.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </main>
</template>
