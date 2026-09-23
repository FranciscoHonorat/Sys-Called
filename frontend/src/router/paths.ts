export const paths = {
  login: '/',
  register: '/criar-conta',
  recoverPassword: '/recuperar-senha',
  changePassword: '/trocar-senha',
  userHome: '/usuario',
  supportHome: '/suporte',
  adminHome: '/admin',
  tickets: '/chamados',
  ticketDetail: '/chamados/:id',
  users: '/admin/usuarios',
  supports: '/admin/suportes',
} as const

export function ticketDetailPath(id: string): string {
  return paths.ticketDetail.replace(':id', id)
}
