import type { RouteMeta } from 'vue-router'

import type { User } from '../auth/authService'
import { homePathFor } from './homePath'
import { paths } from './paths'

export interface Destination {
  path: string
  meta: RouteMeta
}

export function resolveNavigation(user: User | null, to: Destination): true | string {
  if (!user) {
    return isRestricted(to) ? paths.login : true
  }
  if (user.mustChangePassword) {
    return to.path === paths.changePassword ? true : paths.changePassword
  }
  return isOutOfPlace(user, to) ? homePathFor(user.role) : true
}

function isRestricted(to: Destination): boolean {
  return to.meta.role !== undefined || to.meta.requiresAuth === true
}

function isOutOfPlace(user: User, to: Destination): boolean {
  const isLogin = to.path === paths.login
  const isAnotherRoleArea = to.meta.role !== undefined && to.meta.role !== user.role
  return isLogin || isAnotherRoleArea
}
