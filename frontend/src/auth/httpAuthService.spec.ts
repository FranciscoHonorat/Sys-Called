import { createHttpAuthService } from './httpAuthService'

function base64Url(value: object): string {
  return btoa(JSON.stringify(value)).replace(/=+$/, '').replace(/\+/g, '-').replace(/\//g, '_')
}

function fakeToken(claims: object): string {
  return `${base64Url({ alg: 'EdDSA', typ: 'JWT' })}.${base64Url(claims)}.signature`
}

function jsonResponse(status: number, body: object): Response {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

function serviceAnswering(status: number, body: object) {
  const fetchFn = vi.fn().mockResolvedValue(jsonResponse(status, body))
  return { fetchFn, service: createHttpAuthService(fetchFn) }
}

describe('httpAuthService', () => {
  const token = fakeToken({ sub: 'admin-1', name: 'Administradora', role: 'admin' })

  it('logs in against the employees service and keeps the session cookie', async () => {
    const { fetchFn, service } = serviceAnswering(200, { access_token: token, token_type: 'Bearer', expires_in: 900 })

    await service.login('admin', 'senha123')

    expect(fetchFn).toHaveBeenCalledWith('/api/employees/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: 'admin', password: 'senha123' }),
    })
  })

  it('returns who logged in, read from the access token', async () => {
    const { service } = serviceAnswering(200, { access_token: token })

    const user = await service.login('admin', 'senha123')

    expect(user).toEqual({ id: 'admin-1', name: 'Administradora', role: 'admin' })
  })

  it('keeps the access token in memory for the next requests', async () => {
    const { service } = serviceAnswering(200, { access_token: token })

    expect(service.accessToken()).toBeNull()
    await service.login('admin', 'senha123')

    expect(service.accessToken()).toBe(token)
  })

  it('restores the session from the refresh cookie', async () => {
    const { fetchFn, service } = serviceAnswering(200, { access_token: token })

    const user = await service.restore()

    expect(fetchFn).toHaveBeenCalledWith('/api/employees/auth/refresh', { method: 'POST', credentials: 'include' })
    expect(user).toEqual({ id: 'admin-1', name: 'Administradora', role: 'admin' })
    expect(service.accessToken()).toBe(token)
  })

  it('restores nobody when there is no valid session', async () => {
    const { service } = serviceAnswering(401, { error: 'invalid refresh token' })

    expect(await service.restore()).toBeNull()
    expect(service.accessToken()).toBeNull()
  })

  it('ends the session on the server and forgets the token', async () => {
    const fetchFn = vi.fn()
      .mockResolvedValueOnce(jsonResponse(200, { access_token: token }))
      .mockResolvedValueOnce(new Response(null, { status: 204 }))
    const service = createHttpAuthService(fetchFn)
    await service.login('admin', 'senha123')

    await service.logout()

    expect(fetchFn).toHaveBeenLastCalledWith('/api/employees/auth/logout', { method: 'POST', credentials: 'include' })
    expect(service.accessToken()).toBeNull()
  })

  it('forgets the token even when the server cannot be reached', async () => {
    const fetchFn = vi.fn()
      .mockResolvedValueOnce(jsonResponse(200, { access_token: token }))
      .mockRejectedValueOnce(new TypeError('network down'))
    const service = createHttpAuthService(fetchFn)
    await service.login('admin', 'senha123')

    await service.logout()

    expect(service.accessToken()).toBeNull()
  })

  it('rejects refused credentials without keeping any token', async () => {
    const { service } = serviceAnswering(401, { error: 'invalid credentials' })

    await expect(service.login('admin', 'wrong')).rejects.toThrow('invalid credentials')
    expect(service.accessToken()).toBeNull()
  })
})
