import { render, screen, waitFor } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { authServiceKey } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey, type TicketsApi } from '../tickets/ticketsApi'
import NewTicketView from './NewTicketView.vue'

async function renderNewTicket(api: TicketsApi) {
  const session = createSession()
  session.start({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  await router.push(paths.newTicket)
  render(NewTicketView, {
    global: {
      plugins: [router],
      provide: { [sessionKey]: session, [authServiceKey]: fakeAuthService(), [ticketsApiKey]: api },
    },
  })
  return router
}

describe('NewTicketView', () => {
  it('opens a ticket with title, description and priority, then shows the ticket list', async () => {
    const api = fakeTicketsApi()
    const router = await renderNewTicket(api)

    await userEvent.type(screen.getByLabelText('Título'), 'Impressora')
    await userEvent.type(screen.getByLabelText('Descrição'), 'Não imprime')
    await userEvent.selectOptions(screen.getByLabelText('Prioridade'), 'Alta')
    await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))

    expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'High' })
    await waitFor(() => expect(router.currentRoute.value.path).toBe(paths.tickets))
  })

  it('starts with medium priority', async () => {
    await renderNewTicket(fakeTicketsApi())

    expect(screen.getByLabelText('Prioridade')).toHaveValue('Medium')
  })

  it('requires a title and a description', async () => {
    const api = fakeTicketsApi()
    await renderNewTicket(api)

    await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))

    expect(screen.getByText('Informe o título')).toBeInTheDocument()
    expect(screen.getByText('Informe a descrição')).toBeInTheDocument()
    expect(api.open).not.toHaveBeenCalled()
  })

  it('stays on the form and explains when the ticket cannot be opened', async () => {
    const api = fakeTicketsApi({ open: vi.fn().mockRejectedValue(new Error('internal error')) })
    const router = await renderNewTicket(api)

    await userEvent.type(screen.getByLabelText('Título'), 'Impressora')
    await userEvent.type(screen.getByLabelText('Descrição'), 'Não imprime')
    await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível abrir o chamado')
    expect(router.currentRoute.value.path).toBe(paths.newTicket)
  })
})
