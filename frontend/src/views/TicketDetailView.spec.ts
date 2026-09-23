import { render, screen, waitFor, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { authServiceKey, type User } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { fakeAuthService } from '../test/fakeAuthService'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { SessionExpiredError, ticketsApiKey, type TicketDetail, type TicketsApi } from '../tickets/ticketsApi'
import TicketDetailView from './TicketDetailView.vue'

const requester: User = { id: 'user-1', name: 'Usuário Padrão', role: 'user' }
const assignedAgent: User = { id: 'agent-1', name: 'Ana Souza', role: 'support' }

const openTicket: TicketDetail = {
  ticket_id: 't-1',
  title: 'Impressora',
  description: 'Não imprime nada',
  status: 'Open',
  priority: 'High',
  assignee_id: 'agent-1',
  requester_id: 'user-1',
  created_at: '2026-09-22T13:00:00Z',
  responses: [
    { response_id: 'r-1', author_id: 'agent-1', content: 'Verificando', created_at: '2026-09-22T14:00:00Z' },
    { response_id: 'r-2', author_id: 'user-1', content: 'Obrigado', created_at: '2026-09-22T15:00:00Z' },
  ],
}

async function renderDetail(user: User, ticket: TicketDetail, overrides: Partial<TicketsApi> = {}) {
  const session = createSession()
  session.start(user)
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/chamados/:id', component: TicketDetailView },
      { path: '/:pathMatch(.*)*', component: { template: '<div />' } },
    ],
  })
  await router.push(`/chamados/${ticket.ticket_id}`)
  const api = fakeTicketsApi({
    get: vi.fn().mockResolvedValue(ticket),
    responsibles: vi.fn().mockResolvedValue([
      { id: 'agent-1', name: 'Ana Souza' },
      { id: 'agent-2', name: 'Bruno Lima' },
    ]),
    ...overrides,
  })
  render(TicketDetailView, {
    props: { id: ticket.ticket_id },
    global: {
      plugins: [router],
      provide: { [sessionKey]: session, [authServiceKey]: fakeAuthService(), [ticketsApiKey]: api },
    },
  })
  return { api, router, session }
}

async function renderLoadedDetail(user: User, ticket: TicketDetail, overrides: Partial<TicketsApi> = {}) {
  const rendered = await renderDetail(user, ticket, overrides)
  await screen.findByRole('heading', { name: ticket.title })
  return rendered
}

describe('TicketDetailView', () => {
  it('shows the ticket details with the responsible by name', async () => {
    await renderLoadedDetail(requester, openTicket)

    const details = screen.getByRole('list', { name: 'Detalhes' })
    expect(within(details).getByText('Aberto')).toBeInTheDocument()
    expect(within(details).getByText('Alta')).toBeInTheDocument()
    expect(within(details).getByText('Ana Souza')).toBeInTheDocument()
    expect(within(details).getByText('22/09/2026')).toBeInTheDocument()
    expect(screen.getByText('Não imprime nada')).toBeInTheDocument()
  })

  it('shows the conversation with the author of each response', async () => {
    await renderLoadedDetail(requester, openTicket)

    const responses = within(screen.getByRole('list', { name: 'Respostas' })).getAllByRole('listitem')
    expect(responses.map((item) => item.textContent)).toEqual([
      expect.stringContaining('Ana Souza'),
      expect.stringContaining('user-1'),
    ])
    expect(responses[0]).toHaveTextContent('Verificando')
  })

  describe('offers only the actions the user can perform', () => {
    const actionButtons = () => screen.queryAllByRole('button').map((b) => b.textContent?.trim()).filter((t) => t !== 'Sair')

    it('lets the requester edit and reply', async () => {
      await renderLoadedDetail(requester, openTicket)

      expect(actionButtons()).toEqual(['Editar', 'Responder'])
      expect(screen.queryByLabelText('Atribuir a')).not.toBeInTheDocument()
    })

    it('lets the assigned agent manage, work on and reply to the ticket', async () => {
      await renderLoadedDetail(assignedAgent, openTicket)

      expect(actionButtons()).toEqual([
        'Atribuir', 'Distribuir automaticamente', 'Alterar prioridade', 'Iniciar atendimento', 'Fechar chamado', 'Responder',
      ])
      expect(screen.getByLabelText('Atribuir a')).toBeInTheDocument()
      expect(screen.getByLabelText('Nova prioridade')).toBeInTheDocument()
    })

    it('offers nothing on a closed ticket', async () => {
      await renderLoadedDetail(assignedAgent, { ...openTicket, status: 'Closed' })

      expect(actionButtons()).toEqual([])
      expect(screen.queryByLabelText('Sua resposta')).not.toBeInTheDocument()
    })
  })

  describe('performs the actions and shows the updated ticket', () => {
    it.each([
      ['Iniciar atendimento', undefined, (api: TicketsApi) => expect(api.start).toHaveBeenCalledWith('t-1')],
      ['Fechar chamado', undefined, (api: TicketsApi) => expect(api.close).toHaveBeenCalledWith('t-1')],
      ['Distribuir automaticamente', undefined, (api: TicketsApi) => expect(api.autoAssign).toHaveBeenCalledWith('t-1')],
      ['Atribuir', ['Atribuir a', 'Bruno Lima'], (api: TicketsApi) => expect(api.assign).toHaveBeenCalledWith('t-1', 'agent-2')],
      ['Alterar prioridade', ['Nova prioridade', 'Baixa'], (api: TicketsApi) => expect(api.changePriority).toHaveBeenCalledWith('t-1', 'Low')],
    ] as const)('%s', async (button, selection, expectCall) => {
      const { api } = await renderLoadedDetail(assignedAgent, openTicket)
      if (selection) {
        await userEvent.selectOptions(screen.getByLabelText(selection[0]), selection[1])
      }

      await userEvent.click(screen.getByRole('button', { name: button }))

      expectCall(api)
      await waitFor(() => expect(api.get).toHaveBeenCalledTimes(2))
    })

    it('sends a reply and clears the field', async () => {
      const { api } = await renderLoadedDetail(requester, openTicket)

      await userEvent.type(screen.getByLabelText('Sua resposta'), 'Ainda não funciona')
      await userEvent.click(screen.getByRole('button', { name: 'Responder' }))

      expect(api.respond).toHaveBeenCalledWith('t-1', 'Ainda não funciona')
      await waitFor(() => expect(api.get).toHaveBeenCalledTimes(2))
      expect(screen.getByLabelText('Sua resposta')).toHaveValue('')
    })
  })

  describe('editing', () => {
    it('opens the form with the current title and description and saves the changes', async () => {
      const { api } = await renderLoadedDetail(requester, openTicket)

      await userEvent.click(screen.getByRole('button', { name: 'Editar' }))
      expect(screen.getByLabelText('Título')).toHaveValue('Impressora')
      expect(screen.getByLabelText('Descrição')).toHaveValue('Não imprime nada')
      await userEvent.clear(screen.getByLabelText('Título'))
      await userEvent.type(screen.getByLabelText('Título'), 'Impressora do 2º andar')
      await userEvent.click(screen.getByRole('button', { name: 'Salvar' }))

      expect(api.edit).toHaveBeenCalledWith('t-1', { title: 'Impressora do 2º andar', description: 'Não imprime nada' })
      await waitFor(() => expect(screen.queryByRole('button', { name: 'Salvar' })).not.toBeInTheDocument())
      expect(api.get).toHaveBeenCalledTimes(2)
    })

    it('can be cancelled without saving', async () => {
      const { api } = await renderLoadedDetail(requester, openTicket)

      await userEvent.click(screen.getByRole('button', { name: 'Editar' }))
      await userEvent.click(screen.getByRole('button', { name: 'Cancelar' }))

      expect(screen.queryByLabelText('Título')).not.toBeInTheDocument()
      expect(api.edit).not.toHaveBeenCalled()
    })
  })

  describe('when something goes wrong', () => {
    it('explains why an action was refused and keeps the ticket on screen', async () => {
      await renderLoadedDetail(assignedAgent, openTicket, {
        start: vi.fn().mockRejectedValue(new Error('invalid status transition')),
      })

      await userEvent.click(screen.getByRole('button', { name: 'Iniciar atendimento' }))

      expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível concluir a ação: invalid status transition')
      expect(screen.getByRole('heading', { name: 'Impressora' })).toBeInTheDocument()
    })

    it('tells when the ticket does not exist or is not visible', async () => {
      await renderDetail(requester, openTicket, {
        get: vi.fn().mockRejectedValue(new Error('event stream not found for aggregate')),
      })

      expect(await screen.findByText('Chamado não encontrado')).toBeInTheDocument()
    })

    it('sends the user to the login when the session expired', async () => {
      const { router, session } = await renderDetail(requester, openTicket, {
        get: vi.fn().mockRejectedValue(new SessionExpiredError()),
      })

      await waitFor(() => expect(router.currentRoute.value.path).toBe('/'))
      expect(session.isAuthenticated.value).toBe(false)
    })
  })
})
