import { flushPromises } from '@vue/test-utils'
import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter, type Router } from 'vue-router'
import type { Component } from 'vue'

import { AccountRequestError, authServiceKey, type AuthService, type User } from '../auth/authService'
import { createSession, sessionKey, type Session } from '../auth/session'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import ChangePasswordView from './ChangePasswordView.vue'
import RecoverPasswordView from './RecoverPasswordView.vue'
import SignUpView from './SignUpView.vue'

async function renderView(view: Component, authService: AuthService, session: Session = createSession()) {
  const router: Router = createRouter({
    history: createMemoryHistory(),
    routes: Object.values(paths).map((path) => ({ path, component: { template: '<div />' } })),
  })
  await router.push(paths.login)
  render(view, { global: { plugins: [router], provide: { [authServiceKey]: authService, [sessionKey]: session } } })
  return router
}

async function type(label: string, value: string) {
  await userEvent.type(screen.getByLabelText(label), value)
}

describe('SignUpView', () => {
  async function fillForm(password = 'senha-forte', confirmation = password) {
    await type('Nome', 'Maria Lima')
    await type('Username', 'maria')
    await type('Senha', password)
    await type('Confirmar senha', confirmation)
    await userEvent.click(screen.getByRole('button', { name: 'Criar conta' }))
  }

  it('creates the account and explains it must be approved', async () => {
    const authService = fakeAuthService()
    await renderView(SignUpView, authService)

    await fillForm()

    expect(authService.signUp).toHaveBeenCalledWith({ name: 'Maria Lima', username: 'maria', password: 'senha-forte' })
    expect(await screen.findByRole('status')).toHaveTextContent('Conta criada! Aguarde a aprovação do administrador para entrar.')
    expect(screen.getByRole('link', { name: 'Voltar para o login' })).toHaveAttribute('href', '/')
  })

  it('requires every field', async () => {
    const authService = fakeAuthService()
    await renderView(SignUpView, authService)

    await userEvent.click(screen.getByRole('button', { name: 'Criar conta' }))

    expect(screen.getByText('Informe o nome')).toBeInTheDocument()
    expect(screen.getByText('Informe o username')).toBeInTheDocument()
    expect(screen.getByText('Informe a senha')).toBeInTheDocument()
    expect(authService.signUp).not.toHaveBeenCalled()
  })

  it('asks for a password of at least 8 characters typed twice the same', async () => {
    const authService = fakeAuthService()
    await renderView(SignUpView, authService)

    await fillForm('123', '1234')

    expect(screen.getByText('A senha precisa ter pelo menos 8 caracteres')).toBeInTheDocument()
    expect(screen.getByText('As senhas não conferem')).toBeInTheDocument()
    expect(authService.signUp).not.toHaveBeenCalled()
  })

  it('tells when the username is already in use', async () => {
    const authService = fakeAuthService({ signUp: vi.fn().mockRejectedValue(new AccountRequestError(409, 'username already taken')) })
    await renderView(SignUpView, authService)

    await fillForm()

    expect(await screen.findByRole('alert')).toHaveTextContent('Este username já está em uso')
  })

  it('links back to the login', async () => {
    await renderView(SignUpView, fakeAuthService())

    expect(screen.getByRole('link', { name: '← Voltar' })).toHaveAttribute('href', '/')
  })
})

describe('RecoverPasswordView', () => {
  it('asks the administrator for a temporary password', async () => {
    const authService = fakeAuthService()
    await renderView(RecoverPasswordView, authService)

    await type('Username', 'usuario')
    await userEvent.click(screen.getByRole('button', { name: 'Pedir nova senha' }))

    expect(authService.requestPasswordReset).toHaveBeenCalledWith('usuario')
    expect(await screen.findByRole('status')).toHaveTextContent(
      'Pedido enviado! O administrador vai gerar uma senha temporária para você.',
    )
  })

  it('requires the username', async () => {
    const authService = fakeAuthService()
    await renderView(RecoverPasswordView, authService)

    await userEvent.click(screen.getByRole('button', { name: 'Pedir nova senha' }))

    expect(screen.getByText('Informe o username')).toBeInTheDocument()
    expect(authService.requestPasswordReset).not.toHaveBeenCalled()
  })

  it('links back to the login', async () => {
    await renderView(RecoverPasswordView, fakeAuthService())

    expect(screen.getByRole('link', { name: '← Voltar' })).toHaveAttribute('href', '/')
  })
})

describe('ChangePasswordView', () => {
  const temporary: User = { id: 'user-1', name: 'Usuário Padrão', role: 'user', mustChangePassword: true }
  const renewed: User = { id: 'user-1', name: 'Usuário Padrão', role: 'user' }

  function loggedIn(user: User): Session {
    const session = createSession()
    session.start(user)
    return session
  }

  async function submit(current: string, next: string, confirmation = next) {
    await type('Senha atual', current)
    await type('Nova senha', next)
    await type('Confirmar nova senha', confirmation)
    await userEvent.click(screen.getByRole('button', { name: 'Salvar nova senha' }))
  }

  it('explains why a temporary password must be replaced', async () => {
    await renderView(ChangePasswordView, fakeAuthService(), loggedIn(temporary))

    expect(screen.getByText('Você entrou com uma senha temporária. Escolha uma nova senha para continuar.')).toBeInTheDocument()
  })

  it('saves the new password and goes to the home of the user', async () => {
    const authService = fakeAuthService({ changePassword: vi.fn().mockResolvedValue(renewed) })
    const session = loggedIn(temporary)
    const router = await renderView(ChangePasswordView, authService, session)

    await submit('Temp-1234', 'nova-senha')
    await flushPromises()

    expect(authService.changePassword).toHaveBeenCalledWith('Temp-1234', 'nova-senha')
    expect(session.user.value).toEqual(renewed)
    expect(router.currentRoute.value.path).toBe(paths.userHome)
  })

  it('tells when the current password is wrong', async () => {
    const authService = fakeAuthService({
      changePassword: vi.fn().mockRejectedValue(new AccountRequestError(401, 'invalid credentials')),
    })
    await renderView(ChangePasswordView, authService, loggedIn(temporary))

    await submit('errada', 'nova-senha')

    expect(await screen.findByRole('alert')).toHaveTextContent('Senha atual incorreta')
  })

  it('checks the new password before sending it', async () => {
    const authService = fakeAuthService()
    await renderView(ChangePasswordView, authService, loggedIn(temporary))

    await submit('Temp-1234', '123', '321')

    expect(screen.getByText('A senha precisa ter pelo menos 8 caracteres')).toBeInTheDocument()
    expect(screen.getByText('As senhas não conferem')).toBeInTheDocument()
    expect(authService.changePassword).not.toHaveBeenCalled()
  })
})
