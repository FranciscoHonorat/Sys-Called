export const TicketStatus = {
  Open: 'Open',
  InProgress: 'In Progress',
  Closed: 'Closed',
} as const

export type TicketStatus = (typeof TicketStatus)[keyof typeof TicketStatus]

const statusLabels: Record<TicketStatus, string> = {
  [TicketStatus.Open]: 'Aberto',
  [TicketStatus.InProgress]: 'Em andamento',
  [TicketStatus.Closed]: 'Fechado',
}

export function statusLabel(status: TicketStatus): string {
  return statusLabels[status]
}
