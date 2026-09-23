<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import { homePathFor } from '../router/homePath'
import { paths } from '../router/paths'

const props = defineProps<{ to?: string }>()

const session = inject(sessionKey)!
const destination = computed(
  () => props.to ?? (session.user.value ? homePathFor(session.user.value.role) : paths.login),
)
</script>

<template>
  <RouterLink :to="destination" class="mb-4 inline-block text-sm text-blue-700 hover:underline">← Voltar</RouterLink>
</template>
