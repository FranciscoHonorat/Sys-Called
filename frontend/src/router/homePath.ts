import type { Role } from '../auth/authService'
import { paths } from './paths'

const homePaths: Record<Role, string> = {
  user: paths.userHome,
  support: paths.supportHome,
  admin: paths.adminHome,
}

export function homePathFor(role: Role): string {
  return homePaths[role]
}
