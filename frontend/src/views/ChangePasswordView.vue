<script setup lang="ts">
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'

import { AccountRequestError, authServiceKey } from '../auth/authService'
import { sessionKey } from '../auth/session'
import { hasErrors } from '../auth/validateCredentials'
import { validateNewPassword } from '../auth/validateNewPassword'
import AuthCard from '../components/AuthCard.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import { homePathFor } from '../router/homePath'

const authService = inject(authServiceKey)!
const session = inject(sessionKey)!
const router = useRouter()

const forced = session.user.value?.mustChangePassword === true
const current = ref('')
const password = ref('')
const confirmation = ref('')
const errors = ref<Record<string, string | undefined>>({})
const failure = ref('')

async function submit() {
  errors.value = {
    ...validateNewPassword(password.value, confirmation.value),
    ...(current.value ? {} : { current: 'Informe a senha atual' }),
  }
  if (hasErrors(errors.value)) {
    return
  }
  failure.value = ''
  try {
    const user = await authService.changePassword(current.value, password.value)
    session.start(user)
    await router.push(homePathFor(user.role))
  } catch (error) {
    failure.value =
      error instanceof AccountRequestError && error.status === 401
        ? 'Senha atual incorreta'
        : 'Não foi possível trocar a senha'
  }
}
</script>

<template>
  <AuthCard>
    <BackLink v-if="!forced" />
    <form class="space-y-4" novalidate @submit.prevent="submit">
      <h2 class="text-lg font-semibold text-slate-800">Trocar senha</h2>
      <p v-if="forced" class="rounded bg-amber-50 p-3 text-sm text-amber-800">
        Você entrou com uma senha temporária. Escolha uma nova senha para continuar.
      </p>
      <p v-if="failure" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ failure }}</p>
      <FormField id="current" v-model="current" label="Senha atual" type="password" autocomplete="current-password" :error="errors.current" />
      <FormField id="password" v-model="password" label="Nova senha" type="password" autocomplete="new-password" :error="errors.password" />
      <FormField
        id="confirmation"
        v-model="confirmation"
        label="Confirmar nova senha"
        type="password"
        autocomplete="new-password"
        :error="errors.confirmation"
      />
      <button type="submit" class="w-full rounded bg-blue-600 py-2 font-medium text-white hover:bg-blue-700">
        Salvar nova senha
      </button>
    </form>
  </AuthCard>
</template>
