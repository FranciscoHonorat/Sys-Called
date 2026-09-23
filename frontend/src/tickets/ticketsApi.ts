import type { InjectionKey } from 'vue'

import type { User } from '../auth/authService'
import type { TicketStatus } from './status'

export interface TokenSource {
  accessToken(): string | null
  restore(): Promise<User | null>
}

export interface Ticket {
  ticket_id: string
  title: string
  description: string
  status: TicketStatus
  priority?: string
  assignee_id?: string
  requester_id?: string
  created_at: string
}

export interface TicketResponse {
  response_id: string
  author_id: string
  content: string
  created_at: string
}

export interface TicketDetail extends Ticket {
  responses?: TicketResponse[]
}

export interface Responsible {
  id: string
  name: string
}

export class SessionExpiredError extends Error {
  constructor() {
    super('session expired')
  }
}

export interface NewTicket {
  title: string
  description: string
  priority?: string
}

export interface TicketsApi {
  list(): Promise<Ticket[]>
  open(ticket: NewTicket): Promise<string>
  get(id: string): Promise<TicketDetail>
  responsibles(): Promise<Responsible[]>
  edit(id: string, changes: { title: string; description: string }): Promise<void>
  assign(id: string, assigneeId: string): Promise<void>
  autoAssign(id: string): Promise<string>
  changePriority(id: string, priority: string): Promise<void>
  start(id: string): Promise<void>
  close(id: string): Promise<void>
  respond(id: string, content: string): Promise<void>
}

interface ApiCall {
  method: 'GET' | 'POST' | 'PUT'
  path: string
  payload?: unknown
}

function ticketPath(id: string, action?: string): string {
  return action ? `/tickets/${id}/${action}` : `/tickets/${id}`
}

export function createTicketsApi(tokens: TokenSource, fetchFn: typeof fetch = fetch): TicketsApi {
  function send({ method, path, payload }: ApiCall) {
    return fetchFn(`/api/tickets${path}`, {
      method,
      headers: { Authorization: `Bearer ${tokens.accessToken()}`, 'Content-Type': 'application/json' },
      ...(payload === undefined ? {} : { body: JSON.stringify(payload) }),
    })
  }

  async function sendRenewingSession(call: ApiCall): Promise<Response> {
    const response = await send(call)
    if (response.status !== 401) {
      return response
    }
    if (!(await tokens.restore())) {
      throw new SessionExpiredError()
    }
    return send(call)
  }

  async function request<T>(call: ApiCall): Promise<T> {
    const response = await sendRenewingSession(call)
    if (response.status === 204) {
      return undefined as T
    }
    const body = await response.json()
    if (!response.ok) {
      throw new Error(body.error ?? 'request failed')
    }
    return body as T
  }

  return {
    list: () => request<Ticket[]>({ method: 'GET', path: '/tickets' }),
    open: async (ticket) =>
      (await request<{ ticket_id: string }>({ method: 'POST', path: '/tickets', payload: ticket })).ticket_id,
    get: (id) => request<TicketDetail>({ method: 'GET', path: ticketPath(id) }),
    responsibles: () => request<Responsible[]>({ method: 'GET', path: '/responsibles' }),
    edit: (id, changes) => request<void>({ method: 'PUT', path: ticketPath(id), payload: changes }),
    assign: (id, assigneeId) =>
      request<void>({ method: 'POST', path: ticketPath(id, 'assign'), payload: { assignee_id: assigneeId } }),
    autoAssign: async (id) =>
      (await request<{ assignee_id: string }>({ method: 'POST', path: ticketPath(id, 'assign/auto') })).assignee_id,
    changePriority: (id, priority) =>
      request<void>({ method: 'POST', path: ticketPath(id, 'priority'), payload: { priority } }),
    start: (id) => request<void>({ method: 'POST', path: ticketPath(id, 'start') }),
    close: (id) => request<void>({ method: 'POST', path: ticketPath(id, 'close') }),
    respond: (id, content) => request<void>({ method: 'POST', path: ticketPath(id, 'responses'), payload: { content } }),
  }
}

export const ticketsApiKey: InjectionKey<TicketsApi> = Symbol('ticketsApi')
