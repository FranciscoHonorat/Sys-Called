export const minPasswordLength = 8

export interface NewPasswordErrors {
  password?: string
  confirmation?: string
}

export function validateNewPassword(password: string, confirmation: string): NewPasswordErrors {
  const errors: NewPasswordErrors = {}
  if (!password) {
    errors.password = 'Informe a senha'
  } else if ([...password].length < minPasswordLength) {
    errors.password = `A senha precisa ter pelo menos ${minPasswordLength} caracteres`
  }
  if (password !== confirmation) {
    errors.confirmation = 'As senhas não conferem'
  }
  return errors
}
