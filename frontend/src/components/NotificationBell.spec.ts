import { flushPromises } from '@vue/test-utils'
import { render, screen, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory, createRouter } from 'vue-router'

import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey, type NotificationFeed, type TicketsApi } from '../tickets/ticketsApi'
import NotificationBell from './NotificationBell.vue'

const feed: NotificationFeed = {
  unread: 2,
  items: [
    { id: 'n-2', message: 'Nova mensagem em: Impressora', ticket_id: 't-1', created_at: '2026-09-23T10:05:00Z', unread: true },
    { id: 'n-1', message: 'Novo chamado: Monitor', ticket_id: 't-2', created_at: '2026-09-23T10:00:00Z', unread: true },
  ],
}

async function renderBell(api: TicketsApi) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  const rendered = render(NotificationBell, { global: { plugins: [router], provide: { [ticketsApiKey]: api } } })
  await flushPromises()
  return rendered
}

describe('NotificationBell', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it('shows how many notifications are unread', async () => {
    await renderBell(fakeTicketsApi({ notifications: vi.fn().mockResolvedValue(feed) }))

    const bell = screen.getByRole('button', { name: 'Notificações (2 não lidas)' })
    expect(within(bell).getByText('2')).toBeInTheDocument()
  })

  it('lists the notifications linking to their tickets and marks them as read when opened', async () => {
    const api = fakeTicketsApi({ notifications: vi.fn().mockResolvedValue(feed) })
    await renderBell(api)

    await userEvent.click(screen.getByRole('button', { name: /Notificações/ }))

    const panel = screen.getByRole('region', { name: 'Notificações' })
    expect(within(panel).getByRole('link', { name: /Nova mensagem em: Impressora/ })).toHaveAttribute('href', '/chamados/t-1')
    expect(within(panel).getByRole('link', { name: /Novo chamado: Monitor/ })).toHaveAttribute('href', '/chamados/t-2')
    expect(api.markNotificationsRead).toHaveBeenCalled()
    expect(screen.getByRole('button', { name: 'Notificações' })).toBeInTheDocument()
  })

  it('takes the admin to the users screen from an account notification', async () => {
    const accountFeed: NotificationFeed = {
      unread: 1,
      items: [{ id: 'n-3', message: 'Nova conta aguardando aprovação: Maria Lima', created_at: '2026-09-23T10:10:00Z', unread: true }],
    }
    await renderBell(fakeTicketsApi({ notifications: vi.fn().mockResolvedValue(accountFeed) }))

    await userEvent.click(screen.getByRole('button', { name: /Notificações/ }))

    expect(screen.getByRole('link', { name: /Nova conta aguardando aprovação: Maria Lima/ })).toHaveAttribute('href', '/admin/usuarios')
  })

  it('says when there is nothing new', async () => {
    await renderBell(fakeTicketsApi())

    await userEvent.click(screen.getByRole('button', { name: 'Notificações' }))

    expect(screen.getByText('Nenhuma notificação')).toBeInTheDocument()
  })

  it('checks again every 30 seconds while on screen', async () => {
    vi.useFakeTimers()
    const notifications = vi.fn().mockResolvedValueOnce({ unread: 0, items: [] }).mockResolvedValue(feed)
    const { unmount } = await renderBell(fakeTicketsApi({ notifications }))

    await vi.advanceTimersByTimeAsync(30_000)

    expect(notifications).toHaveBeenCalledTimes(2)
    expect(screen.getByRole('button', { name: 'Notificações (2 não lidas)' })).toBeInTheDocument()

    unmount()
    await vi.advanceTimersByTimeAsync(60_000)
    expect(notifications).toHaveBeenCalledTimes(2)
  })
})
