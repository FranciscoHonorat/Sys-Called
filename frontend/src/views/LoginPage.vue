<script setup lang="ts">
import { inject } from 'vue'
import { useRouter } from 'vue-router'

import type { User } from '../auth/authService'
import { sessionKey } from '../auth/session'
import { homePathFor } from '../router/homePath'
import LoginView from './LoginView.vue'

const session = inject(sessionKey)!
const router = useRouter()

function enter(user: User) {
  session.start(user)
  router.push(homePathFor(user.role))
}
</script>

<template>
  <LoginView @authenticated="enter" />
</template>
