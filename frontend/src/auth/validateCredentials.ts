export interface CredentialErrors {
  username?: string
  password?: string
}

export function validateCredentials(username: string, password: string): CredentialErrors {
  const errors: CredentialErrors = {}
  if (!username) {
    errors.username = 'Informe o username'
  }
  if (!password) {
    errors.password = 'Informe a senha'
  }
  return errors
}

export function hasErrors(errors: object): boolean {
  return Object.keys(errors).length > 0
}
