import type { AuthService, Role, User } from './authService'
import { decodeJwtClaims } from './jwt'

type Fetch = typeof fetch

export interface HttpAuthService extends AuthService {
  accessToken(): string | null
}

export function createHttpAuthService(fetchFn: Fetch = fetch): HttpAuthService {
  let token: string | null = null

  function adopt(accessToken: string): User {
    token = accessToken
    return userFromToken(accessToken)
  }

  return {
    async login(username, password) {
      const response = await fetchFn('/api/employees/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, password }),
      })
      const body = await response.json()
      if (!response.ok) {
        throw new Error(body.error ?? 'login failed')
      }
      return adopt(body.access_token)
    },
    async restore() {
      const response = await fetchFn('/api/employees/auth/refresh', { method: 'POST', credentials: 'include' })
      if (!response.ok) {
        return null
      }
      const body = await response.json()
      return adopt(body.access_token)
    },
    async logout() {
      token = null
      await fetchFn('/api/employees/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined)
    },
    accessToken() {
      return token
    },
  }
}

function userFromToken(token: string): User {
  const claims = decodeJwtClaims(token)
  return { id: String(claims.sub), name: String(claims.name), role: claims.role as Role }
}
