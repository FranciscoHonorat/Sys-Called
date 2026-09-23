import { AccountPendingError, AccountRequestError, type AuthService, type Role, type User } from './authService'
import { decodeJwtClaims } from './jwt'

type Fetch = typeof fetch

export interface HttpAuthService extends AuthService {
  accessToken(): string | null
}

const jsonHeaders = { 'Content-Type': 'application/json' }

export function createHttpAuthService(fetchFn: Fetch = fetch): HttpAuthService {
  let token: string | null = null

  function adopt(accessToken: string): User {
    token = accessToken
    return userFromToken(accessToken)
  }

  async function post(path: string, payload: object, headers: Record<string, string> = {}): Promise<void> {
    const response = await fetchFn(`/api/employees/auth/${path}`, {
      method: 'POST',
      headers: { ...jsonHeaders, ...headers },
      body: JSON.stringify(payload),
    })
    if (!response.ok) {
      const body = await response.json().catch(() => ({}))
      throw new AccountRequestError(response.status, body.error ?? 'request failed')
    }
  }

  async function restore(): Promise<User | null> {
    const response = await fetchFn('/api/employees/auth/refresh', { method: 'POST', credentials: 'include' })
    if (!response.ok) {
      return null
    }
    const body = await response.json()
    return adopt(body.access_token)
  }

  return {
    async login(username, password) {
      const response = await fetchFn('/api/employees/auth/login', {
        method: 'POST',
        headers: jsonHeaders,
        credentials: 'include',
        body: JSON.stringify({ username, password }),
      })
      const body = await response.json()
      if (response.status === 403) {
        throw new AccountPendingError()
      }
      if (!response.ok) {
        throw new Error(body.error ?? 'login failed')
      }
      return adopt(body.access_token)
    },
    restore,
    async logout() {
      token = null
      await fetchFn('/api/employees/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined)
    },
    signUp(data) {
      return post('signup', data)
    },
    requestPasswordReset(username) {
      return post('password-reset-requests', { username })
    },
    async changePassword(current, next) {
      await post(
        'change-password',
        { current_password: current, new_password: next },
        { Authorization: `Bearer ${token}` },
      )
      const renewed = await restore()
      if (!renewed) {
        throw new AccountRequestError(401, 'session expired')
      }
      return renewed
    },
    accessToken() {
      return token
    },
  }
}

function userFromToken(token: string): User {
  const claims = decodeJwtClaims(token)
  const user: User = { id: String(claims.sub), name: String(claims.name), role: claims.role as Role }
  return claims.must_change_password === true ? { ...user, mustChangePassword: true } : user
}
