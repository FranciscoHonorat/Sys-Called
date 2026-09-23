import { inject } from 'vue'
import { useRouter } from 'vue-router'

import { paths } from '../router/paths'
import { authServiceKey } from './authService'
import { sessionKey } from './session'

export function useSessionExit() {
  const session = inject(sessionKey)!
  const authService = inject(authServiceKey)!
  const router = useRouter()

  async function expire() {
    session.end()
    await router.push(paths.login)
  }

  async function logout() {
    await authService.logout()
    await expire()
  }

  return { logout, expire }
}
