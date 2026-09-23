<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import { useSessionExit } from '../auth/useSessionExit'

const session = inject(sessionKey)!
const user = computed(() => session.user.value)
const { logout } = useSessionExit()
</script>

<template>
  <header class="mb-8 flex items-start justify-between">
    <div>
      <h1 class="text-2xl font-bold text-slate-800">Ticket System</h1>
      <p class="text-slate-600">Olá, {{ user?.name }}</p>
    </div>
    <button type="button" class="rounded border border-slate-300 px-3 py-1 text-sm text-slate-700 hover:bg-white" @click="logout">
      Sair
    </button>
  </header>
</template>
