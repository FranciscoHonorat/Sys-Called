import { ref } from 'vue'

import { AccountPendingError, type AuthService, type User } from './authService'
import { hasErrors, validateCredentials, type CredentialErrors } from './validateCredentials'

export function useLoginForm(authService: AuthService, onAuthenticated: (user: User) => void) {
  const username = ref('')
  const password = ref('')
  const errors = ref<CredentialErrors>({})
  const loginError = ref('')

  async function submit() {
    errors.value = validateCredentials(username.value, password.value)
    if (hasErrors(errors.value)) {
      return
    }

    try {
      onAuthenticated(await authService.login(username.value, password.value))
    } catch (error) {
      loginError.value =
        error instanceof AccountPendingError
          ? 'Sua conta ainda aguarda a aprovação do administrador'
          : 'Username ou senha inválidos'
    }
  }

  return { username, password, errors, loginError, submit }
}
