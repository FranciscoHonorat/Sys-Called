import { fakeAuthService } from '../test/fakeAuthService'
import type { User } from './authService'
import { restoreSession } from './restoreSession'
import { createSession } from './session'

describe('restoreSession', () => {
  it('picks up the session that survived a page reload', async () => {
    const ana: User = { id: 'agent-1', name: 'Ana Souza', role: 'support' }
    const session = createSession()

    await restoreSession(fakeAuthService({ restore: vi.fn().mockResolvedValue(ana) }), session)

    expect(session.user.value).toEqual(ana)
  })

  it('leaves nobody logged in when there was no session', async () => {
    const session = createSession()

    await restoreSession(fakeAuthService(), session)

    expect(session.isAuthenticated.value).toBe(false)
  })

  it('opens the app logged out if the session cannot be checked', async () => {
    const session = createSession()

    await restoreSession(fakeAuthService({ restore: vi.fn().mockRejectedValue(new TypeError('network down')) }), session)

    expect(session.isAuthenticated.value).toBe(false)
  })
})
