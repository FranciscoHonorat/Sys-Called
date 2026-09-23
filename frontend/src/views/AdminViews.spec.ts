import { render, screen, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'
import type { Component } from 'vue'

import { authServiceKey } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { employeesApiKey, type EmployeesApi } from '../employees/employeesApi'
import { fakeAuthService } from '../test/fakeAuthService'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey, type TicketsApi } from '../tickets/ticketsApi'
import SupportsView from './SupportsView.vue'
import UsersView from './UsersView.vue'

async function renderAsAdmin(view: Component, tickets: TicketsApi, employees: EmployeesApi) {
  const session = createSession()
  session.start({ id: 'admin-1', name: 'Administradora', role: 'admin' })
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  render(view, {
    global: {
      plugins: [router],
      provide: {
        [sessionKey]: session,
        [authServiceKey]: fakeAuthService(),
        [ticketsApiKey]: tickets,
        [employeesApiKey]: employees,
      },
    },
  })
}

function rows(): string[][] {
  return screen.getAllByRole('row').slice(1).map((row) =>
    within(row).getAllByRole('cell').map((cell) => cell.textContent?.trim() ?? ''),
  )
}

function fakeEmployeesApi(overrides: Partial<EmployeesApi> = {}): EmployeesApi {
  return {
    list: vi.fn().mockResolvedValue([]),
    approve: vi.fn().mockResolvedValue(undefined),
    issueTemporaryPassword: vi.fn().mockResolvedValue('Temp-1234'),
    ...overrides,
  }
}

const noEmployees = fakeEmployeesApi()

describe('SupportsView', () => {
  it('shows how many tickets each agent has open, in progress and closed', async () => {
    const tickets = fakeTicketsApi({
      supportWorkload: vi.fn().mockResolvedValue([
        { id: 'agent-1', name: 'Ana Souza', open: 1, in_progress: 2, closed: 3 },
        { id: 'agent-2', name: 'Bruno Lima', open: 0, in_progress: 0, closed: 0 },
      ]),
    })
    await renderAsAdmin(SupportsView, tickets, noEmployees)

    await screen.findByText('Ana Souza')
    expect(rows()).toEqual([
      ['Ana Souza', '1', '2', '3', '6'],
      ['Bruno Lima', '0', '0', '0', '0'],
    ])
    expect(screen.getByRole('link', { name: '← Voltar' })).toHaveAttribute('href', '/admin')
  })

  it('explains when the numbers cannot be loaded', async () => {
    await renderAsAdmin(SupportsView, fakeTicketsApi({ supportWorkload: vi.fn().mockRejectedValue(new Error('x')) }), noEmployees)

    expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível carregar')
  })
})

describe('UsersView', () => {
  const admin = { id: 'admin-1', name: 'Administradora', username: 'admin', role: 'admin', status: 'active', password_reset_requested: false }
  const ana = { id: 'agent-1', name: 'Ana Souza', username: 'ana', role: 'support', status: 'active', password_reset_requested: false }
  const forgetful = { id: 'user-1', name: 'Usuário Padrão', username: 'usuario', role: 'user', status: 'active', password_reset_requested: true }
  const newcomer = { id: 'user-9', name: 'Maria Lima', username: 'maria', role: 'user', status: 'pending', password_reset_requested: false }

  it('lists every employee with username, role and situation', async () => {
    const employees = fakeEmployeesApi({ list: vi.fn().mockResolvedValue([admin, ana, forgetful, newcomer]) })
    await renderAsAdmin(UsersView, fakeTicketsApi(), employees)

    await screen.findByText('Ana Souza')
    expect(rows().map((row) => row.slice(0, 4))).toEqual([
      ['Administradora', 'admin', 'Administrador', 'Ativo'],
      ['Ana Souza', 'ana', 'Suporte', 'Ativo'],
      ['Usuário Padrão', 'usuario', 'Usuário', 'Pediu nova senha'],
      ['Maria Lima', 'maria', 'Usuário', 'Aguardando aprovação'],
    ])
  })

  it('approves a pending account and shows it active', async () => {
    const list = vi.fn()
      .mockResolvedValueOnce([newcomer])
      .mockResolvedValue([{ ...newcomer, status: 'active' }])
    const employees = fakeEmployeesApi({ list })
    await renderAsAdmin(UsersView, fakeTicketsApi(), employees)

    await userEvent.click(await screen.findByRole('button', { name: 'Aprovar Maria Lima' }))

    expect(employees.approve).toHaveBeenCalledWith('user-9')
    expect(await screen.findByText('Ativo')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Aprovar Maria Lima' })).not.toBeInTheDocument()
  })

  it('shows a temporary password once so the admin can hand it over', async () => {
    const employees = fakeEmployeesApi({ list: vi.fn().mockResolvedValue([forgetful]) })
    await renderAsAdmin(UsersView, fakeTicketsApi(), employees)

    await userEvent.click(await screen.findByRole('button', { name: 'Gerar senha temporária para Usuário Padrão' }))

    expect(employees.issueTemporaryPassword).toHaveBeenCalledWith('user-1')
    const dialog = await screen.findByRole('dialog', { name: 'Senha temporária' })
    expect(within(dialog).getByText('Temp-1234')).toBeInTheDocument()
    expect(within(dialog).getByText(/Usuário Padrão vai precisar trocá-la no próximo acesso/)).toBeInTheDocument()

    await userEvent.click(within(dialog).getByRole('button', { name: 'Fechar' }))
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })

  it('offers no temporary password for accounts still waiting for approval', async () => {
    await renderAsAdmin(UsersView, fakeTicketsApi(), fakeEmployeesApi({ list: vi.fn().mockResolvedValue([newcomer]) }))

    await screen.findByText('Maria Lima')
    expect(screen.queryByRole('button', { name: /Gerar senha temporária/ })).not.toBeInTheDocument()
  })

  it('explains when the users cannot be loaded', async () => {
    await renderAsAdmin(UsersView, fakeTicketsApi(), fakeEmployeesApi({ list: vi.fn().mockRejectedValue(new Error('x')) }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível carregar')
  })
})
