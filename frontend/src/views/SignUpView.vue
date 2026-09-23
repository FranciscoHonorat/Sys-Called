<script setup lang="ts">
import { inject, ref } from 'vue'

import { AccountRequestError, authServiceKey } from '../auth/authService'
import { hasErrors } from '../auth/validateCredentials'
import { validateNewPassword } from '../auth/validateNewPassword'
import AuthCard from '../components/AuthCard.vue'
import AuthConfirmation from '../components/AuthConfirmation.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import { paths } from '../router/paths'

const authService = inject(authServiceKey)!

const name = ref('')
const username = ref('')
const password = ref('')
const confirmation = ref('')
const errors = ref<Record<string, string | undefined>>({})
const failure = ref('')
const created = ref(false)

function validate() {
  const found: Record<string, string | undefined> = { ...validateNewPassword(password.value, confirmation.value) }
  if (!name.value.trim()) {
    found.name = 'Informe o nome'
  }
  if (!username.value.trim()) {
    found.username = 'Informe o username'
  }
  return found
}

async function submit() {
  errors.value = validate()
  if (hasErrors(errors.value)) {
    return
  }
  failure.value = ''
  try {
    await authService.signUp({ name: name.value.trim(), username: username.value.trim(), password: password.value })
    created.value = true
  } catch (error) {
    failure.value =
      error instanceof AccountRequestError && error.status === 409
        ? 'Este username já está em uso'
        : 'Não foi possível criar a conta'
  }
}
</script>

<template>
  <AuthCard>
    <BackLink :to="paths.login" />
    <AuthConfirmation v-if="created">Conta criada! Aguarde a aprovação do administrador para entrar.</AuthConfirmation>
    <form v-else class="space-y-4" novalidate @submit.prevent="submit">
      <h2 class="text-lg font-semibold text-slate-800">Criar conta</h2>
      <p v-if="failure" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ failure }}</p>
      <FormField id="name" v-model="name" label="Nome" autocomplete="name" :error="errors.name" />
      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="errors.username" />
      <FormField id="password" v-model="password" label="Senha" type="password" autocomplete="new-password" :error="errors.password" />
      <FormField
        id="confirmation"
        v-model="confirmation"
        label="Confirmar senha"
        type="password"
        autocomplete="new-password"
        :error="errors.confirmation"
      />
      <button type="submit" class="w-full rounded bg-blue-600 py-2 font-medium text-white hover:bg-blue-700">Criar conta</button>
    </form>
  </AuthCard>
</template>
