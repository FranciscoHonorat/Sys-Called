<script setup lang="ts">
import { inject } from 'vue'

import { authServiceKey, type User } from '../auth/authService'
import { useLoginForm } from '../auth/useLoginForm'
import AuthCard from '../components/AuthCard.vue'
import FormField from '../components/FormField.vue'
import { paths } from '../router/paths'

const emit = defineEmits<{ authenticated: [user: User] }>()

const { username, password, errors, loginError, submit } = useLoginForm(inject(authServiceKey)!, (user) =>
  emit('authenticated', user),
)
</script>

<template>
  <AuthCard>
    <form class="space-y-4" novalidate @submit.prevent="submit">
      <p v-if="loginError" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">{{ loginError }}</p>

      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="errors.username" />
      <FormField
        id="password"
        v-model="password"
        label="Senha"
        type="password"
        autocomplete="current-password"
        :error="errors.password"
      />

      <button
        type="submit"
        class="w-full rounded bg-blue-600 py-2 font-medium text-white hover:bg-blue-700"
      >
        Entrar
      </button>
    </form>

    <footer class="mt-6 flex justify-between text-sm">
      <RouterLink :to="paths.register" class="text-blue-600 hover:underline">Criar conta</RouterLink>
      <RouterLink :to="paths.recoverPassword" class="text-blue-600 hover:underline">Recuperar senha</RouterLink>
    </footer>
  </AuthCard>
</template>
