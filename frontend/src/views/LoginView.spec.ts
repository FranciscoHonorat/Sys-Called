import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { AccountPendingError, authServiceKey, type AuthService, type User } from '../auth/authService'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import LoginView from './LoginView.vue'

function renderLogin(authService: AuthService = fakeAuthService()) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: Object.values(paths).map((path) => ({ path, component: { template: '<div />' } })),
  })

  return render(LoginView, {
    global: { plugins: [router], provide: { [authServiceKey]: authService } },
  })
}

async function submitCredentials(username: string, password: string) {
  await userEvent.type(screen.getByLabelText('Username'), username)
  await userEvent.type(screen.getByLabelText('Senha'), password)
  await userEvent.click(screen.getByRole('button', { name: 'Entrar' }))
}

describe('LoginView', () => {
  it('shows the system name and subtitle', () => {
    renderLogin()

    expect(screen.getByRole('heading', { name: 'Ticket System' })).toBeInTheDocument()
    expect(screen.getByText('Sistema de chamados internos')).toBeInTheDocument()
  })

  it('asks for username and a masked password', () => {
    renderLogin()

    expect(screen.getByLabelText('Username')).toHaveAttribute('type', 'text')
    expect(screen.getByLabelText('Senha')).toHaveAttribute('type', 'password')
    expect(screen.getByRole('button', { name: 'Entrar' })).toBeInTheDocument()
  })

  it('links to account creation and password recovery', () => {
    renderLogin()

    expect(screen.getByRole('link', { name: 'Criar conta' })).toHaveAttribute('href', '/criar-conta')
    expect(screen.getByRole('link', { name: 'Recuperar senha' })).toHaveAttribute('href', '/recuperar-senha')
  })

  it('requires username and password before trying to log in', async () => {
    const authService = fakeAuthService()
    renderLogin(authService)

    await userEvent.click(screen.getByRole('button', { name: 'Entrar' }))

    expect(screen.getByText('Informe o username')).toBeInTheDocument()
    expect(screen.getByText('Informe a senha')).toBeInTheDocument()
    expect(authService.login).not.toHaveBeenCalled()
  })

  it('logs in with the typed credentials and reports the authenticated user', async () => {
    const user: User = { id: 'admin-1', name: 'Ana', role: 'admin' }
    const authService = fakeAuthService({ login: vi.fn().mockResolvedValue(user) })
    const { emitted } = renderLogin(authService)

    await submitCredentials('ana', 'secret')

    expect(authService.login).toHaveBeenCalledWith('ana', 'secret')
    expect(emitted().authenticated).toEqual([[user]])
  })

  it('shows an error and stays on the page when the credentials are rejected', async () => {
    const authService = fakeAuthService({ login: vi.fn().mockRejectedValue(new Error('invalid credentials')) })
    const { emitted } = renderLogin(authService)

    await submitCredentials('ana', 'wrong')

    expect(await screen.findByRole('alert')).toHaveTextContent('Username ou senha inválidos')
    expect(emitted().authenticated).toBeUndefined()
  })

  it('explains that a new account still waits for the administrator', async () => {
    const authService = fakeAuthService({ login: vi.fn().mockRejectedValue(new AccountPendingError()) })
    renderLogin(authService)

    await submitCredentials('maria', 'senha-forte')

    expect(await screen.findByRole('alert')).toHaveTextContent('Sua conta ainda aguarda a aprovação do administrador')
  })
})
