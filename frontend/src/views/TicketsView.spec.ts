import { render, screen, waitFor, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { authServiceKey } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { paths } from '../router/paths'
import { fakeAuthService } from '../test/fakeAuthService'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { SessionExpiredError, ticketsApiKey, type Ticket, type TicketsApi } from '../tickets/ticketsApi'
import TicketsView from './TicketsView.vue'

const tickets: Ticket[] = [
  { ticket_id: 't-1', title: 'Impressora', description: 'x', status: 'Open', priority: 'High', created_at: '2026-09-22T13:00:00Z' },
  { ticket_id: 't-2', title: 'Cadeira', description: 'x', status: 'In Progress', priority: 'Low', assignee_id: 'agent-1', created_at: '2026-09-21T13:00:00Z' },
  { ticket_id: 't-3', title: 'Monitor', description: 'x', status: 'Closed', assignee_id: 'agent-2', created_at: '2026-09-20T13:00:00Z', closed_at: '2026-09-21T15:00:00Z' },
]

async function renderTicketsAt(url: string, list: TicketsApi['list'], overrides: Partial<TicketsApi> = {}) {
  const session = createSession()
  session.start({ id: 'admin-1', name: 'Administradora', role: 'admin' })
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  await router.push(url)
  const api = fakeTicketsApi({
    list,
    responsibles: vi.fn().mockResolvedValue([{ id: 'agent-1', name: 'Ana Souza' }]),
    ...overrides,
  })
  render(TicketsView, {
    global: {
      plugins: [router],
      provide: { [sessionKey]: session, [authServiceKey]: fakeAuthService(), [ticketsApiKey]: api },
    },
  })
  return { router, session, api }
}

function rowTexts(): string[][] {
  return screen.getAllByRole('row').slice(1).map((row) =>
    within(row).getAllByRole('cell').map((cell) => cell.textContent?.trim() ?? ''),
  )
}

describe('TicketsView', () => {
  it('lists the tickets with their status, priority, responsible and opening date', async () => {
    await renderTicketsAt(paths.tickets, vi.fn().mockResolvedValue(tickets))

    await screen.findByText('Impressora')
    expect(rowTexts()).toEqual([
      ['Impressora', 'Aberto', 'Alta', '—', '22/09/2026', '—'],
      ['Cadeira', 'Em andamento', 'Baixa', 'Ana Souza', '21/09/2026', '—'],
      ['Monitor', 'Fechado', '—', 'agent-2', '20/09/2026', '21/09/2026'],
    ])
  })

  it('goes back to the home of the logged in role', async () => {
    await renderTicketsAt(paths.tickets, vi.fn().mockResolvedValue([]))

    expect(screen.getByRole('link', { name: '← Voltar' })).toHaveAttribute('href', '/admin')
  })

  it('links each ticket to its details', async () => {
    await renderTicketsAt(paths.tickets, vi.fn().mockResolvedValue(tickets))

    expect(await screen.findByRole('link', { name: 'Impressora' })).toHaveAttribute('href', '/chamados/t-1')
  })

  it('shows only the tickets with the status chosen in the address', async () => {
    await renderTicketsAt(`${paths.tickets}?status=In+Progress`, vi.fn().mockResolvedValue(tickets))

    await screen.findByText('Cadeira')
    expect(rowTexts().map((row) => row[0])).toEqual(['Cadeira'])
    expect(screen.getByRole('link', { name: 'Em atendimento' })).toHaveAttribute('aria-current', 'page')
  })

  it('says when there is nothing to show', async () => {
    await renderTicketsAt(paths.tickets, vi.fn().mockResolvedValue([]))

    expect(await screen.findByText('Nenhum chamado encontrado')).toBeInTheDocument()
  })

  it('sends the user to the login when the session expired', async () => {
    const { router, session } = await renderTicketsAt(paths.tickets, vi.fn().mockRejectedValue(new SessionExpiredError()))

    await waitFor(() => expect(router.currentRoute.value.path).toBe(paths.login))
    expect(session.isAuthenticated.value).toBe(false)
  })

  it('opens a ticket in a dialog and refreshes the list right away', async () => {
    const list = vi.fn().mockResolvedValueOnce([]).mockResolvedValueOnce([tickets[0]])
    const { api } = await renderTicketsAt(paths.tickets, list, { open: vi.fn().mockResolvedValue('t-1') })
    await screen.findByText('Nenhum chamado encontrado')

    await userEvent.click(screen.getByRole('button', { name: 'Abrir novo chamado' }))
    const dialog = screen.getByRole('dialog', { name: 'Abrir novo chamado' })
    await userEvent.type(within(dialog).getByLabelText('Título'), 'Impressora')
    await userEvent.type(within(dialog).getByLabelText('Descrição'), 'Não imprime')
    await userEvent.click(within(dialog).getByRole('button', { name: 'Abrir chamado' }))

    expect(api.open).toHaveBeenCalled()
    expect(await screen.findByRole('link', { name: 'Impressora' })).toBeInTheDocument()
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })
})
