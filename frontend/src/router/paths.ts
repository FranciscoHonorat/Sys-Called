export const paths = {
  login: '/',
  register: '/criar-conta',
  recoverPassword: '/recuperar-senha',
  userHome: '/usuario',
  supportHome: '/suporte',
  adminHome: '/admin',
  tickets: '/chamados',
  newTicket: '/chamados/novo',
  ticketDetail: '/chamados/:id',
  users: '/admin/usuarios',
  supports: '/admin/suportes',
} as const

export function ticketDetailPath(id: string): string {
  return paths.ticketDetail.replace(':id', id)
}
