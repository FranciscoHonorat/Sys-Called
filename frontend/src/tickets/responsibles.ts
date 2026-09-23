import type { Responsible } from './ticketsApi'

export function nameResolver(responsibles: Responsible[]): (id: string | undefined) => string {
  const names = new Map(responsibles.map((r) => [r.id, r.name]))
  return (id) => (id ? (names.get(id) ?? id) : '—')
}
