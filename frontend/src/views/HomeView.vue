<script setup lang="ts">
import { computed, inject, ref } from 'vue'

import { sessionKey } from '../auth/session'
import NewTicketModal from '../components/NewTicketModal.vue'
import PageLayout from '../components/PageLayout.vue'
import { menuFor } from '../router/menu'
import { ticketDetailPath } from '../router/paths'
import { canOpenTickets } from '../tickets/permissions'

const session = inject(sessionKey)!
const role = computed(() => session.user.value?.role)
const menu = computed(() => (role.value ? menuFor(role.value) : []))

const creating = ref(false)
const openedTicketId = ref<string | null>(null)

function onOpened(ticketId: string) {
  creating.value = false
  openedTicketId.value = ticketId
}
</script>

<template>
  <PageLayout>
    <p v-if="openedTicketId" role="status" class="mb-4 rounded bg-green-50 p-3 text-sm text-green-800">
      Chamado aberto com sucesso.
      <RouterLink :to="ticketDetailPath(openedTicketId)" class="font-medium underline">Ver chamado</RouterLink>
    </p>

    <nav class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <button
        v-if="role && canOpenTickets(role)"
        type="button"
        class="rounded-lg bg-blue-600 p-6 text-left font-medium text-white shadow hover:bg-blue-700"
        @click="creating = true"
      >
        Abrir novo chamado
      </button>
      <RouterLink
        v-for="item in menu"
        :key="item.label"
        :to="item.to"
        class="rounded-lg bg-white p-6 font-medium text-slate-700 shadow hover:bg-blue-50 hover:text-blue-700"
      >
        {{ item.label }}
      </RouterLink>
    </nav>

    <NewTicketModal :open="creating" @close="creating = false" @opened="onOpened" />
  </PageLayout>
</template>
