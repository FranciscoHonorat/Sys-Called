import type { InjectionKey } from 'vue'

import { createApiClient, SessionExpiredError, type TokenSource } from '../api/apiClient'
import type { TicketStatus } from './status'

export interface Ticket {
  ticket_id: string
  title: string
  description: string
  status: TicketStatus
  priority?: string
  assignee_id?: string
  requester_id?: string
  created_at: string
  closed_at?: string
  resolution?: string
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

export interface Notification {
  id: string
  message: string
  ticket_id?: string
  created_at: string
  unread: boolean
}

export interface AgentWorkload {
  id: string
  name: string
  open: number
  in_progress: number
  closed: number
}

export interface NotificationFeed {
  unread: number
  items: Notification[]
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
  close(id: string, resolution: string): Promise<void>
  respond(id: string, content: string): Promise<void>
  notifications(): Promise<NotificationFeed>
  markNotificationsRead(): Promise<void>
  supportWorkload(): Promise<AgentWorkload[]>
}

function ticketPath(id: string, action?: string): string {
  return action ? `/tickets/${id}/${action}` : `/tickets/${id}`
}

export function createTicketsApi(tokens: TokenSource, fetchFn: typeof fetch = fetch): TicketsApi {
  const request = createApiClient('/api/tickets', tokens, fetchFn)

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
    close: (id, resolution) =>
      request<void>({ method: 'POST', path: ticketPath(id, 'close'), payload: { resolution } }),
    respond: (id, content) => request<void>({ method: 'POST', path: ticketPath(id, 'responses'), payload: { content } }),
    notifications: () => request<NotificationFeed>({ method: 'GET', path: '/notifications' }),
    markNotificationsRead: () => request<void>({ method: 'POST', path: '/notifications/read' }),
    supportWorkload: () => request<AgentWorkload[]>({ method: 'GET', path: '/responsibles/workload' }),
  }
}

export const ticketsApiKey: InjectionKey<TicketsApi> = Symbol('ticketsApi')

export { SessionExpiredError, type TokenSource }
