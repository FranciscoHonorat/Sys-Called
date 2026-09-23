const dateFormat = new Intl.DateTimeFormat('pt-BR')

export function formatDate(iso: string): string {
  return dateFormat.format(new Date(iso))
}
