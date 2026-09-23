import type { RouteLocationRaw } from 'vue-router'

import type { Role } from '../auth/authService'
import { TicketStatus } from '../tickets/status'
import { paths } from './paths'

export interface MenuItem {
  label: string
  to: RouteLocationRaw
}

function ticketsWithStatus(status: TicketStatus): RouteLocationRaw {
  return { path: paths.tickets, query: { status } }
}

const menus: Record<Role, MenuItem[]> = {
  user: [{ label: 'Meus chamados', to: paths.tickets }],
  support: [
    { label: 'Chamados abertos', to: ticketsWithStatus(TicketStatus.Open) },
    { label: 'Em atendimento', to: ticketsWithStatus(TicketStatus.InProgress) },
    { label: 'Fechados por mim', to: ticketsWithStatus(TicketStatus.Closed) },
  ],
  admin: [
    { label: 'Todos os chamados', to: paths.tickets },
    { label: 'Usuários', to: paths.users },
    { label: 'Suportes', to: paths.supports },
  ],
}

export function menuFor(role: Role): MenuItem[] {
  return menus[role]
}
