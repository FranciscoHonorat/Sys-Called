import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory } from 'vue-router'

import App from './App.vue'
import { authServiceKey, type User } from './auth/authService'
import { createSession, sessionKey } from './auth/session'
import { createAppRouter } from './router'
import { paths } from './router/paths'
import { fakeAuthService } from './test/fakeAuthService'

async function renderAppAt(path: string, loggedInAs: User) {
  const session = createSession()
  const router = createAppRouter(session, createMemoryHistory())
  const authService = fakeAuthService({ login: vi.fn().mockResolvedValue(loggedInAs) })
  await router.push(path)
  render(App, {
    global: { plugins: [router], provide: { [sessionKey]: session, [authServiceKey]: authService } },
  })
  return router
}

describe('App', () => {
  const ana: User = { id: 'agent-1', name: 'Ana Souza', role: 'support' }

  it('takes a support agent from the login to the support home', async () => {
    const router = await renderAppAt(paths.login, ana)

    await userEvent.type(screen.getByLabelText('Username'), 'ana')
    await userEvent.type(screen.getByLabelText('Senha'), 'senha123')
    await userEvent.click(screen.getByRole('button', { name: 'Entrar' }))

    expect(await screen.findByText('Olá, Ana Souza')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Chamados abertos' })).toBeInTheDocument()
    expect(router.currentRoute.value.path).toBe(paths.supportHome)
  })

  it('shows the login to a visitor opening a restricted area', async () => {
    const router = await renderAppAt(paths.adminHome, ana)

    expect(await screen.findByRole('heading', { name: 'Ticket System' })).toBeInTheDocument()
    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(router.currentRoute.value.path).toBe(paths.login)
  })

  it('tells that pages still under construction are coming soon', async () => {
    await renderAppAt(paths.register, ana)

    expect(await screen.findByText('Em breve')).toBeInTheDocument()
  })
})
