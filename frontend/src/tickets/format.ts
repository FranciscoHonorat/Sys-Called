const dateFormat = new Intl.DateTimeFormat('pt-BR')

export function formatDate(iso: string | undefined): string {
  return iso ? dateFormat.format(new Date(iso)) : '—'
}

const MINUTE = 60_000
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

export function formatDuration(fromIso: string, toIso: string): string {
  const elapsed = new Date(toIso).getTime() - new Date(fromIso).getTime()
  const days = Math.floor(elapsed / DAY)
  const hours = Math.floor((elapsed % DAY) / HOUR)
  const minutes = Math.floor((elapsed % HOUR) / MINUTE)

  if (days > 0) {
    const dayLabel = days === 1 ? '1 dia' : `${days} dias`
    return hours > 0 ? `${dayLabel} e ${hours} h` : dayLabel
  }
  if (hours > 0) {
    return minutes > 0 ? `${hours} h ${minutes} min` : `${hours} h`
  }
  return minutes > 0 ? `${minutes} min` : 'menos de 1 min'
}
