import { fakeJwt } from '../test/fakeJwt'
import { AccountPendingError, AccountRequestError } from './authService'
import { createHttpAuthService } from './httpAuthService'

function jsonResponse(status: number, body: object): Response {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

function serviceAnswering(status: number, body: object) {
  const fetchFn = vi.fn().mockResolvedValue(jsonResponse(status, body))
  return { fetchFn, service: createHttpAuthService(fetchFn) }
}

describe('httpAuthService', () => {
  const token = fakeJwt({ sub: 'admin-1', name: 'Administradora', role: 'admin' })

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

  it('tells when the account still waits for the administrator approval', async () => {
    const { service } = serviceAnswering(403, { error: 'account waiting for the administrator approval' })

    await expect(service.login('maria', 'senha-forte')).rejects.toBeInstanceOf(AccountPendingError)
  })

  it('knows when the user must choose a new password', async () => {
    const temporary = fakeJwt({ sub: 'user-1', name: 'Usuário Padrão', role: 'user', must_change_password: true })
    const { service } = serviceAnswering(200, { access_token: temporary })

    const user = await service.login('usuario', 'Temp-1234')

    expect(user).toEqual({ id: 'user-1', name: 'Usuário Padrão', role: 'user', mustChangePassword: true })
  })

  it('creates an account', async () => {
    const fetchFn = vi.fn().mockResolvedValue(new Response(null, { status: 201 }))
    const service = createHttpAuthService(fetchFn)

    await service.signUp({ name: 'Maria Lima', username: 'maria', password: 'senha-forte' })

    expect(fetchFn).toHaveBeenCalledWith('/api/employees/auth/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: 'Maria Lima', username: 'maria', password: 'senha-forte' }),
    })
  })

  it('explains why an account was refused', async () => {
    const { service } = serviceAnswering(409, { error: 'username already taken' })

    const refusal = service.signUp({ name: 'Outra Ana', username: 'ana', password: 'senha-forte' })

    await expect(refusal).rejects.toBeInstanceOf(AccountRequestError)
    await expect(refusal).rejects.toMatchObject({ status: 409 })
  })

  it('asks the administrator for a new password', async () => {
    const fetchFn = vi.fn().mockResolvedValue(new Response(null, { status: 202 }))
    const service = createHttpAuthService(fetchFn)

    await service.requestPasswordReset('usuario')

    expect(fetchFn).toHaveBeenCalledWith('/api/employees/auth/password-reset-requests', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: 'usuario' }),
    })
  })

  it('changes the password and renews the session without the obligation', async () => {
    const temporary = fakeJwt({ sub: 'user-1', name: 'Usuário Padrão', role: 'user', must_change_password: true })
    const renewed = fakeJwt({ sub: 'user-1', name: 'Usuário Padrão', role: 'user' })
    const fetchFn = vi.fn()
      .mockResolvedValueOnce(jsonResponse(200, { access_token: temporary }))
      .mockResolvedValueOnce(new Response(null, { status: 204 }))
      .mockResolvedValueOnce(jsonResponse(200, { access_token: renewed }))
    const service = createHttpAuthService(fetchFn)
    await service.login('usuario', 'Temp-1234')

    const user = await service.changePassword('Temp-1234', 'nova-senha')

    expect(fetchFn).toHaveBeenNthCalledWith(2, '/api/employees/auth/change-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${temporary}` },
      body: JSON.stringify({ current_password: 'Temp-1234', new_password: 'nova-senha' }),
    })
    expect(user).toEqual({ id: 'user-1', name: 'Usuário Padrão', role: 'user' })
    expect(service.accessToken()).toBe(renewed)
  })

  it('refuses a wrong current password', async () => {
    const fetchFn = vi.fn()
      .mockResolvedValueOnce(jsonResponse(200, { access_token: token }))
      .mockResolvedValueOnce(jsonResponse(401, { error: 'invalid credentials' }))
    const service = createHttpAuthService(fetchFn)
    await service.login('admin', 'senha123')

    await expect(service.changePassword('errada', 'nova-senha')).rejects.toMatchObject({ status: 401 })
  })
})
