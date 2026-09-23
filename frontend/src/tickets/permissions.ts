import type { Role, User } from '../auth/authService'
import { TicketStatus } from './status'
import type { Ticket } from './ticketsApi'

export type TicketAction = 'edit' | 'manage' | 'start' | 'close' | 'respond'

export function allowedActions(ticket: Ticket, user: User): TicketAction[] {
  if (ticket.status === TicketStatus.Closed) {
    return []
  }

  const isAdmin = user.role === 'admin'
  const isSupport = user.role === 'support'
  const isRequester = ticket.requester_id === user.id
  const isOpen = ticket.status === TicketStatus.Open
  const worksOnIt = isSupport && ticket.assignee_id === user.id
  const rules: Array<[TicketAction, boolean]> = [
    ['edit', isAdmin || (isRequester && isOpen)],
    ['manage', isAdmin || isSupport],
    ['start', worksOnIt && isOpen],
    ['close', worksOnIt],
    ['respond', true],
  ]
  return rules.filter(([, allowed]) => allowed).map(([action]) => action)
}

export function canOpenTickets(role: Role): boolean {
  return role !== 'support'
}
