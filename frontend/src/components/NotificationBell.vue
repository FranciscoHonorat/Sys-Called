<script setup lang="ts">
import { computed, inject, ref } from 'vue'

import { usePolling } from '../composables/usePolling'
import { paths, ticketDetailPath } from '../router/paths'
import { formatDate } from '../tickets/format'
import { ticketsApiKey, type Notification, type NotificationFeed } from '../tickets/ticketsApi'

const props = withDefaults(defineProps<{ pollEveryMs?: number }>(), { pollEveryMs: 30_000 })

const api = inject(ticketsApiKey)!

const feed = ref<NotificationFeed>({ unread: 0, items: [] })
const open = ref(false)

const label = computed(() =>
  feed.value.unread > 0 ? `Notificações (${feed.value.unread} não lidas)` : 'Notificações',
)

async function refresh() {
  feed.value = await api.notifications().catch(() => feed.value)
}

async function toggle() {
  open.value = !open.value
  if (open.value && feed.value.unread > 0) {
    await api.markNotificationsRead()
    feed.value = { unread: 0, items: feed.value.items.map((item) => ({ ...item, unread: false })) }
  }
}

function destination(item: Notification): string {
  return item.ticket_id ? ticketDetailPath(item.ticket_id) : paths.users
}

usePolling(refresh, props.pollEveryMs)
</script>

<template>
  <div class="relative">
    <button
      type="button"
      :aria-label="label"
      :aria-expanded="open"
      class="relative rounded-full p-2 text-slate-600 hover:bg-white hover:text-blue-700"
      @click="toggle"
    >
      <svg aria-hidden="true" viewBox="0 0 24 24" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 17h5l-1.4-1.4A2 2 0 0 1 18 14.2V11a6 6 0 1 0-12 0v3.2a2 2 0 0 1-.6 1.4L4 17h5m6 0v1a3 3 0 1 1-6 0v-1m6 0H9" />
      </svg>
      <span
        v-if="feed.unread > 0"
        class="absolute -right-0.5 -top-0.5 min-w-5 rounded-full bg-red-600 px-1 text-center text-xs font-semibold text-white"
      >{{ feed.unread }}</span>
    </button>

    <section
      v-if="open"
      aria-label="Notificações"
      class="absolute right-0 z-40 mt-2 w-80 rounded-lg bg-white p-2 shadow-xl ring-1 ring-slate-200"
    >
      <p v-if="feed.items.length === 0" class="p-3 text-sm text-slate-500">Nenhuma notificação</p>
      <ul v-else class="max-h-96 divide-y divide-slate-100 overflow-y-auto">
        <li v-for="item in feed.items" :key="item.id">
          <RouterLink
            :to="destination(item)"
            class="block rounded p-3 text-sm hover:bg-slate-50"
            :class="item.unread ? 'font-medium text-slate-800' : 'text-slate-600'"
            @click="open = false"
          >
            {{ item.message }}
            <span class="block text-xs font-normal text-slate-400">{{ formatDate(item.created_at) }}</span>
          </RouterLink>
        </li>
      </ul>
    </section>
  </div>
</template>
