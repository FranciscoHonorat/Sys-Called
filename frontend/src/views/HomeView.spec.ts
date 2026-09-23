import { render, screen, waitFor, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { authServiceKey, type User } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey } from '../tickets/ticketsApi'
import HomeView from './HomeView.vue'

function renderHomeFor(user: User, authService = fakeAuthService(), api = fakeTicketsApi()) {
  const session = createSession()
  session.start(user)
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  render(HomeView, {
    global: {
      plugins: [router],
      provide: { [sessionKey]: session, [authServiceKey]: authService, [ticketsApiKey]: api },
    },
  })
  return { session, router }
}

function menuLinks(): Array<[string, string | null]> {
  return screen.getAllByRole('link').map((link) => [link.textContent?.trim() ?? '', link.getAttribute('href')])
}

describe('HomeView', () => {
  it('shows the notification bell next to the logout', async () => {
    renderHomeFor({ id: 'agent-1', name: 'Ana Souza', role: 'support' })

    expect(await screen.findByRole('button', { name: /Notificações/ })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Sair' })).toBeInTheDocument()
  })

  it('greets the logged in user by name', () => {
    renderHomeFor({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })

    expect(screen.getByText('Olá, Usuário Padrão')).toBeInTheDocument()
  })

  it('offers a regular user to open a ticket and follow their own tickets', () => {
    renderHomeFor({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })

    expect(screen.getByRole('button', { name: 'Abrir novo chamado' })).toBeInTheDocument()
    expect(menuLinks()).toEqual([['Meus chamados', '/chamados']])
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

    expect(screen.getByRole('button', { name: 'Abrir novo chamado' })).toBeInTheDocument()
    expect(menuLinks()).toEqual([
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

  describe('opening a ticket without leaving the page', () => {
    const user: User = { id: 'user-1', name: 'Usuário Padrão', role: 'user' }

    it('opens the form in a dialog and confirms with a link to the new ticket', async () => {
      const api = fakeTicketsApi({ open: vi.fn().mockResolvedValue('t-9') })
      const { router } = renderHomeFor(user, fakeAuthService(), api)

      await userEvent.click(screen.getByRole('button', { name: 'Abrir novo chamado' }))
      const dialog = screen.getByRole('dialog', { name: 'Abrir novo chamado' })
      await userEvent.type(within(dialog).getByLabelText('Título'), 'Impressora')
      await userEvent.type(within(dialog).getByLabelText('Descrição'), 'Não imprime')
      await userEvent.click(within(dialog).getByRole('button', { name: 'Abrir chamado' }))

      expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'Medium' })
      await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument())
      expect(screen.getByRole('status')).toHaveTextContent('Chamado aberto')
      expect(screen.getByRole('link', { name: 'Ver chamado' })).toHaveAttribute('href', '/chamados/t-9')
      expect(router.currentRoute.value.path).toBe('/')
    })

    it('closes the dialog without opening anything when going back', async () => {
      const api = fakeTicketsApi()
      renderHomeFor(user, fakeAuthService(), api)

      await userEvent.click(screen.getByRole('button', { name: 'Abrir novo chamado' }))
      await userEvent.click(screen.getByRole('button', { name: '← Voltar' }))

      await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument())
      expect(api.open).not.toHaveBeenCalled()
    })
  })
})
