import { computed, readonly, ref, type InjectionKey } from 'vue'

import type { User } from './authService'

export function createSession() {
  const user = ref<User | null>(null)

  return {
    user: readonly(user),
    isAuthenticated: computed(() => user.value !== null),
    start(loggedIn: User) {
      user.value = loggedIn
    },
    end() {
      user.value = null
    },
  }
}

export type Session = ReturnType<typeof createSession>

export const sessionKey: InjectionKey<Session> = Symbol('session')
