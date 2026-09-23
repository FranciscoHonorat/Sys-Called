import { createTicketsApi, SessionExpiredError, type TicketsApi, type TokenSource } from './ticketsApi'

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

function tokenSource(tokens: string[], restored = true): TokenSource {
  let current = tokens[0]
  return {
    accessToken: () => current,
    restore: vi.fn(async () => {
      if (!restored) return null
      current = tokens[1]
      return { id: 'user-1', name: 'Usuário', role: 'user' as const }
    }),
  }
}

const ticket = {
  ticket_id: 't-1',
  title: 'Impressora',
  description: 'Não imprime',
  status: 'Open',
  priority: 'High',
  created_at: '2026-09-22T10:00:00Z',
  requester_id: 'user-1',
}

function apiAnswering(status: number, body: unknown) {
  const fetchFn = vi.fn().mockResolvedValue(jsonResponse(status, body))
  return { fetchFn, api: createTicketsApi(tokenSource(['token-1']), fetchFn) }
}

function lastCall(fetchFn: ReturnType<typeof vi.fn>) {
  const [url, init] = fetchFn.mock.calls[fetchFn.mock.calls.length - 1]
  return { url, method: init.method, body: init.body === undefined ? undefined : JSON.parse(init.body) }
}

describe('ticketsApi', () => {
  it('lists the tickets with the access token', async () => {
    const fetchFn = vi.fn().mockResolvedValue(jsonResponse(200, [ticket]))
    const api = createTicketsApi(tokenSource(['token-1']), fetchFn)

    const tickets = await api.list()

    expect(fetchFn).toHaveBeenCalledWith('/api/tickets/tickets', {
      method: 'GET',
      headers: { Authorization: 'Bearer token-1', 'Content-Type': 'application/json' },
    })
    expect(tickets).toEqual([ticket])
  })

  it('opens a ticket and returns its id', async () => {
    const fetchFn = vi.fn().mockResolvedValue(jsonResponse(201, { ticket_id: 't-9' }))
    const api = createTicketsApi(tokenSource(['token-1']), fetchFn)

    const id = await api.open({ title: 'Cadeira', description: 'Quebrada', priority: 'Medium' })

    expect(fetchFn).toHaveBeenCalledWith('/api/tickets/tickets', {
      method: 'POST',
      headers: { Authorization: 'Bearer token-1', 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: 'Cadeira', description: 'Quebrada', priority: 'Medium' }),
    })
    expect(id).toBe('t-9')
  })

  it('renews an expired session once and retries with the new token', async () => {
    const fetchFn = vi.fn()
      .mockResolvedValueOnce(jsonResponse(401, { error: 'unauthenticated' }))
      .mockResolvedValueOnce(jsonResponse(200, [ticket]))
    const tokens = tokenSource(['expired', 'fresh'])
    const api = createTicketsApi(tokens, fetchFn)

    const tickets = await api.list()

    expect(tokens.restore).toHaveBeenCalledTimes(1)
    expect(fetchFn.mock.calls[1][1].headers.Authorization).toBe('Bearer fresh')
    expect(tickets).toEqual([ticket])
  })

  it('reports an expired session when it cannot be renewed', async () => {
    const fetchFn = vi.fn().mockResolvedValue(jsonResponse(401, { error: 'unauthenticated' }))
    const api = createTicketsApi(tokenSource(['expired'], false), fetchFn)

    await expect(api.list()).rejects.toBeInstanceOf(SessionExpiredError)
  })

  it('surfaces the server message on other failures', async () => {
    const fetchFn = vi.fn().mockResolvedValue(jsonResponse(500, { error: 'internal error' }))
    const api = createTicketsApi(tokenSource(['token-1']), fetchFn)

    await expect(api.list()).rejects.toThrow('internal error')
  })

  it('gets one ticket with its responses', async () => {
    const detail = { ...ticket, responses: [{ response_id: 'r-1', author_id: 'agent-1', content: 'Verificando', created_at: '2026-09-22T11:00:00Z' }] }
    const { fetchFn, api } = apiAnswering(200, detail)

    expect(await api.get('t-1')).toEqual(detail)
    expect(lastCall(fetchFn)).toEqual({ url: '/api/tickets/tickets/t-1', method: 'GET', body: undefined })
  })

  it('lists the responsibles', async () => {
    const responsibles = [{ id: 'agent-1', name: 'Ana Souza' }]
    const { fetchFn, api } = apiAnswering(200, responsibles)

    expect(await api.responsibles()).toEqual(responsibles)
    expect(lastCall(fetchFn)).toEqual({ url: '/api/tickets/responsibles', method: 'GET', body: undefined })
  })

  it.each([
    ['edit', (api: TicketsApi) => api.edit('t-1', { title: 'Novo', description: 'Texto' }), 'PUT', '/api/tickets/tickets/t-1', { title: 'Novo', description: 'Texto' }],
    ['assign', (api: TicketsApi) => api.assign('t-1', 'agent-2'), 'POST', '/api/tickets/tickets/t-1/assign', { assignee_id: 'agent-2' }],
    ['change priority', (api: TicketsApi) => api.changePriority('t-1', 'High'), 'POST', '/api/tickets/tickets/t-1/priority', { priority: 'High' }],
    ['start', (api: TicketsApi) => api.start('t-1'), 'POST', '/api/tickets/tickets/t-1/start', undefined],
    ['close', (api: TicketsApi) => api.close('t-1'), 'POST', '/api/tickets/tickets/t-1/close', undefined],
    ['respond', (api: TicketsApi) => api.respond('t-1', 'Resolvido'), 'POST', '/api/tickets/tickets/t-1/responses', { content: 'Resolvido' }],
  ] as const)('%s calls the matching endpoint and accepts an empty answer', async (_, run, method, url, body) => {
    const fetchFn = vi.fn().mockResolvedValue(new Response(null, { status: 204 }))
    const api = createTicketsApi(tokenSource(['token-1']), fetchFn)

    await run(api)

    expect(lastCall(fetchFn)).toEqual({ url, method, body })
  })

  it('assigns automatically and tells who got the ticket', async () => {
    const { fetchFn, api } = apiAnswering(200, { assignee_id: 'agent-3' })

    expect(await api.autoAssign('t-1')).toBe('agent-3')
    expect(lastCall(fetchFn)).toEqual({ url: '/api/tickets/tickets/t-1/assign/auto', method: 'POST', body: undefined })
  })
})
