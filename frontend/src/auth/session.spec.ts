import type { User } from './authService'
import { createSession } from './session'

const ana: User = { id: 'agent-1', name: 'Ana Souza', role: 'support' }

describe('session', () => {
  it('starts without anyone logged in', () => {
    const session = createSession()

    expect(session.user.value).toBeNull()
    expect(session.isAuthenticated.value).toBe(false)
  })

  it('remembers who logged in', () => {
    const session = createSession()

    session.start(ana)

    expect(session.user.value).toEqual(ana)
    expect(session.isAuthenticated.value).toBe(true)
  })

  it('forgets the user when the session ends', () => {
    const session = createSession()
    session.start(ana)

    session.end()

    expect(session.user.value).toBeNull()
    expect(session.isAuthenticated.value).toBe(false)
  })
})
