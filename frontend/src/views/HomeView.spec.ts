import { render, screen, waitFor } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { authServiceKey, type User } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import HomeView from './HomeView.vue'

function renderHomeFor(user: User, authService = fakeAuthService()) {
  const session = createSession()
  session.start(user)
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  render(HomeView, {
    global: { plugins: [router], provide: { [sessionKey]: session, [authServiceKey]: authService } },
  })
  return { session, router }
}

function menuLinks(): Array<[string, string | null]> {
  return screen.getAllByRole('link').map((link) => [link.textContent?.trim() ?? '', link.getAttribute('href')])
}

describe('HomeView', () => {
  it('greets the logged in user by name', () => {
    renderHomeFor({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })

    expect(screen.getByText('Olá, Usuário Padrão')).toBeInTheDocument()
  })

  it('offers a regular user to open a ticket and follow their own tickets', () => {
    renderHomeFor({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })

    expect(menuLinks()).toEqual([
      ['Abrir novo chamado', '/chamados/novo'],
      ['Meus chamados', '/chamados'],
    ])
  })

  it('shows support agents the open tickets and their own work', () => {
    renderHomeFor({ id: 'agent-1', name: 'Ana Souza', role: 'support' })

    expect(menuLinks()).toEqual([
      ['Chamados abertos', '/chamados?status=Open'],
      ['Em atendimento', '/chamados?status=In+Progress'],
      ['Fechados por mim', '/chamados?status=Closed'],
    ])
  })

  it('gives admins every ticket, the users and the support team', () => {
    renderHomeFor({ id: 'admin-1', name: 'Administradora', role: 'admin' })

    expect(menuLinks()).toEqual([
      ['Abrir novo chamado', '/chamados/novo'],
      ['Todos os chamados', '/chamados'],
      ['Usuários', '/admin/usuarios'],
      ['Suportes', '/admin/suportes'],
    ])
  })

  it('logs out and goes back to the login', async () => {
    const authService = fakeAuthService()
    const { session, router } = renderHomeFor({ id: 'user-1', name: 'Usuário Padrão', role: 'user' }, authService)

    await userEvent.click(screen.getByRole('button', { name: 'Sair' }))

    expect(authService.logout).toHaveBeenCalled()
    expect(session.isAuthenticated.value).toBe(false)
    await waitFor(() => expect(router.currentRoute.value.path).toBe(paths.login))
  })
})
