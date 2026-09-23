<script setup lang="ts">
import { inject, ref } from 'vue'

import { authServiceKey } from '../auth/authService'
import AuthCard from '../components/AuthCard.vue'
import AuthConfirmation from '../components/AuthConfirmation.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import { paths } from '../router/paths'

const authService = inject(authServiceKey)!

const username = ref('')
const error = ref('')
const failure = ref('')
const sent = ref(false)

async function submit() {
  error.value = username.value.trim() ? '' : 'Informe o username'
  if (error.value) {
    return
  }
  failure.value = ''
  try {
    await authService.requestPasswordReset(username.value.trim())
    sent.value = true
  } catch {
    failure.value = 'Não foi possível enviar o pedido'
  }
}
</script>

<template>
  <AuthCard>
    <BackLink :to="paths.login" />
    <AuthConfirmation v-if="sent">Pedido enviado! O administrador vai gerar uma senha temporária para você.</AuthConfirmation>
    <form v-else class="space-y-4" novalidate @submit.prevent="submit">
      <h2 class="text-lg font-semibold text-slate-800">Recuperar senha</h2>
      <p class="text-sm text-slate-500">Informe seu username e o administrador vai gerar uma senha temporária.</p>
      <p v-if="failure" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ failure }}</p>
      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="error" />
      <button type="submit" class="w-full rounded bg-blue-600 py-2 font-medium text-white hover:bg-blue-700">
        Pedir nova senha
      </button>
    </form>
  </AuthCard>
</template>
