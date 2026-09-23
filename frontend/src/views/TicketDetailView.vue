<script setup lang="ts">
import { computed, inject, ref } from 'vue'

import { sessionKey } from '../auth/session'
import ClosingReport from '../components/ClosingReport.vue'
import CloseTicketModal from '../components/CloseTicketModal.vue'
import PageLayout from '../components/PageLayout.vue'
import FormField from '../components/FormField.vue'
import TicketActions from '../components/TicketActions.vue'
import TicketEditForm from '../components/TicketEditForm.vue'
import { paths } from '../router/paths'
import { formatDate } from '../tickets/format'
import { allowedActions, type TicketAction } from '../tickets/permissions'
import { priorityLabel } from '../tickets/priority'
import { nameResolver } from '../tickets/responsibles'
import { statusLabel } from '../tickets/status'
import { useTicketDetail } from '../tickets/useTicketDetail'

const props = defineProps<{ id: string }>()

const session = inject(sessionKey)!
const { api, ticket, responsibles, notFound, actionError, perform } = useTicketDetail(props.id)

const reply = ref('')
const closing = ref(false)

async function closeTicket(resolution: string) {
  await perform(() => api.close(props.id, resolution))
  closing.value = false
}
const editing = ref(false)

const actions = computed(() =>
  ticket.value && session.user.value ? allowedActions(ticket.value, session.user.value) : [],
)
const can = (action: TicketAction) => actions.value.includes(action)
const responsibleOptions = computed(() => responsibles.value.map((r) => ({ value: r.id, label: r.name })))

const nameOf = computed(() => nameResolver(responsibles.value))

async function saveEdit(changes: { title: string; description: string }) {
  await perform(() => api.edit(props.id, changes))
  editing.value = false
}

async function sendReply() {
  await perform(() => api.respond(props.id, reply.value))
  reply.value = ''
}
</script>

<template>
  <PageLayout back :back-to="paths.tickets">

    <p v-if="notFound" class="rounded-lg bg-white p-6 text-slate-600 shadow">Chamado não encontrado</p>

    <article v-if="ticket" class="space-y-6 rounded-lg bg-white p-6 shadow">
      <p v-if="actionError" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ actionError }}</p>

      <TicketEditForm
        v-if="editing"
        :title="ticket.title"
        :description="ticket.description"
        @save="saveEdit"
        @cancel="editing = false"
      />
      <header>
        <h2 class="text-xl font-semibold text-slate-800">{{ ticket.title }}</h2>
        <p v-if="!editing" class="mt-2 whitespace-pre-line text-slate-700">{{ ticket.description }}</p>
      </header>

      <ul aria-label="Detalhes" class="grid gap-2 text-sm sm:grid-cols-5">
        <li><span class="block text-slate-500">Status</span>{{ statusLabel(ticket.status) }}</li>
        <li><span class="block text-slate-500">Prioridade</span>{{ priorityLabel(ticket.priority) }}</li>
        <li><span class="block text-slate-500">Responsável</span>{{ nameOf(ticket.assignee_id) }}</li>
        <li><span class="block text-slate-500">Aberto em</span>{{ formatDate(ticket.created_at) }}</li>
        <li v-if="ticket.closed_at"><span class="block text-slate-500">Fechado em</span>{{ formatDate(ticket.closed_at) }}</li>
      </ul>

      <TicketActions
        v-if="actions.length > 0"
        :actions="actions"
        :responsibles="responsibleOptions"
        @assign="(assigneeId) => perform(() => api.assign(id, assigneeId))"
        @auto-assign="perform(() => api.autoAssign(id))"
        @change-priority="(priority) => perform(() => api.changePriority(id, priority))"
        @start="perform(() => api.start(id))"
        @close="closing = true"
        @edit="editing = true"
      />

      <ClosingReport
        v-if="ticket.closed_at"
        :resolution="ticket.resolution"
        :opened-at="ticket.created_at"
        :closed-at="ticket.closed_at"
      />

      <section>
        <h3 class="mb-2 font-semibold text-slate-800">Respostas</h3>
        <ul aria-label="Respostas" class="space-y-3">
          <li v-for="response in ticket.responses ?? []" :key="response.response_id" class="rounded bg-slate-50 p-3">
            <p class="text-xs text-slate-500">{{ nameOf(response.author_id) }} · {{ formatDate(response.created_at) }}</p>
            <p class="text-slate-700">{{ response.content }}</p>
          </li>
        </ul>

        <form v-if="can('respond')" class="mt-4 space-y-2" @submit.prevent="sendReply">
          <FormField id="reply" v-model="reply" label="Sua resposta" multiline />
          <button type="submit" class="action-button">Responder</button>
        </form>
      </section>
    </article>

    <CloseTicketModal
      v-if="ticket"
      :open="closing"
      :opened-at="ticket.created_at"
      @close="closing = false"
      @confirm="closeTicket"
    />
  </PageLayout>
</template>

<style scoped>
@reference 'tailwindcss';

.action-button {
  @apply rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50;
}
</style>
