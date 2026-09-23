<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import AppHeader from '../components/AppHeader.vue'
import { menuFor } from '../router/menu'

const session = inject(sessionKey)!
const menu = computed(() => (session.user.value ? menuFor(session.user.value.role) : []))
</script>

<template>
  <main class="min-h-screen bg-slate-100 p-6">
    <AppHeader />

    <nav class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <RouterLink
        v-for="item in menu"
        :key="item.label"
        :to="item.to"
        class="rounded-lg bg-white p-6 font-medium text-slate-700 shadow hover:bg-blue-50 hover:text-blue-700"
      >
        {{ item.label }}
      </RouterLink>
    </nav>
  </main>
</template>
