import type { User } from '../auth/authService'

export interface TokenSource {
  accessToken(): string | null
  restore(): Promise<User | null>
}

export class SessionExpiredError extends Error {
  constructor() {
    super('session expired')
  }
}

export interface ApiCall {
  method: 'GET' | 'POST' | 'PUT'
  path: string
  payload?: unknown
}

export type Request = <T>(call: ApiCall) => Promise<T>

export function createApiClient(basePath: string, tokens: TokenSource, fetchFn: typeof fetch = fetch): Request {
  function send({ method, path, payload }: ApiCall) {
    return fetchFn(`${basePath}${path}`, {
      method,
      headers: { Authorization: `Bearer ${tokens.accessToken()}`, 'Content-Type': 'application/json' },
      ...(payload === undefined ? {} : { body: JSON.stringify(payload) }),
    })
  }

  async function sendRenewingSession(call: ApiCall): Promise<Response> {
    const response = await send(call)
    if (response.status !== 401) {
      return response
    }
    if (!(await tokens.restore())) {
      throw new SessionExpiredError()
    }
    return send(call)
  }

  return async function request<T>(call: ApiCall): Promise<T> {
    const response = await sendRenewingSession(call)
    if (response.status === 204) {
      return undefined as T
    }
    const body = await response.json()
    if (!response.ok) {
      throw new Error(body.error ?? 'request failed')
    }
    return body as T
  }
}
