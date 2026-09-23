import type { Role } from '../auth/authService'
import type { Employee } from './employeesApi'

const roleLabels: Record<Role, string> = {
  user: 'Usuário',
  support: 'Suporte',
  admin: 'Administrador',
}

export function roleLabel(role: Role): string {
  return roleLabels[role]
}

export function situationLabel(employee: Pick<Employee, 'status' | 'password_reset_requested'>): string {
  if (employee.status === 'pending') {
    return 'Aguardando aprovação'
  }
  return employee.password_reset_requested ? 'Pediu nova senha' : 'Ativo'
}

export function needsAttention(employee: Pick<Employee, 'status' | 'password_reset_requested'>): boolean {
  return employee.status === 'pending' || employee.password_reset_requested
}
