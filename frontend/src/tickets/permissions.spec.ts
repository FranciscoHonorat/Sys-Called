import type { User } from '../auth/authService'
import { allowedActions } from './permissions'
import type { Ticket } from './ticketsApi'

const requester: User = { id: 'user-1', name: 'Usuário', role: 'user' }
const assignedAgent: User = { id: 'agent-1', name: 'Ana', role: 'support' }
const anotherAgent: User = { id: 'agent-2', name: 'Bruno', role: 'support' }
const admin: User = { id: 'admin-1', name: 'Admin', role: 'admin' }

function ticket(status: Ticket['status'], assignee?: string): Ticket {
  return {
    ticket_id: 't-1',
    title: 'Impressora',
    description: 'x',
    status,
    assignee_id: assignee,
    requester_id: 'user-1',
    created_at: '2026-09-22T10:00:00Z',
  }
}

describe('allowedActions', () => {
  it.each([
    ['requester', requester, ticket('Open', 'agent-1'), ['edit', 'respond']],
    ['requester once the work started', requester, ticket('In Progress', 'agent-1'), ['respond']],
    ['assigned agent on an open ticket', assignedAgent, ticket('Open', 'agent-1'), ['manage', 'start', 'close', 'respond']],
    ['assigned agent on a ticket in progress', assignedAgent, ticket('In Progress', 'agent-1'), ['manage', 'close', 'respond']],
    ['another agent', anotherAgent, ticket('Open', 'agent-1'), ['manage', 'respond']],
    ['admin', admin, ticket('Open'), ['edit', 'manage', 'respond']],
    ['admin on a ticket in progress', admin, ticket('In Progress', 'agent-1'), ['edit', 'manage', 'respond']],
    ['admin on a closed ticket', admin, ticket('Closed', 'agent-1'), []],
    ['requester on a closed ticket', requester, ticket('Closed', 'agent-1'), []],
  ] as const)('%s', (_, user, subject, expected) => {
    expect(allowedActions(subject, user)).toEqual(expected)
  })
})
