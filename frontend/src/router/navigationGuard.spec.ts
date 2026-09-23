import type { User } from '../auth/authService'
import { resolveNavigation } from './navigationGuard'
import { paths } from './paths'

const admin: User = { id: 'admin-1', name: 'Admin', role: 'admin' }
const agent: User = { id: 'agent-1', name: 'Ana', role: 'support' }

const adminArea = { path: paths.adminHome, meta: { role: 'admin' as const } }
const loginPage = { path: paths.login, meta: {} }
const registerPage = { path: paths.register, meta: {} }

describe('resolveNavigation', () => {
  it('sends a visitor trying a restricted area back to the login', () => {
    expect(resolveNavigation(null, adminArea)).toBe(paths.login)
  })

  it('lets a visitor reach public pages', () => {
    expect(resolveNavigation(null, loginPage)).toBe(true)
    expect(resolveNavigation(null, registerPage)).toBe(true)
  })

  it('lets a user into the area of their own role', () => {
    expect(resolveNavigation(admin, adminArea)).toBe(true)
  })

  it('sends a user trying the area of another role to their own home', () => {
    expect(resolveNavigation(agent, adminArea)).toBe(paths.supportHome)
  })

  it('requires any logged in user on pages shared by every role', () => {
    const ticketsPage = { path: paths.tickets, meta: { requiresAuth: true } }

    expect(resolveNavigation(null, ticketsPage)).toBe(paths.login)
    expect(resolveNavigation(agent, ticketsPage)).toBe(true)
    expect(resolveNavigation(admin, ticketsPage)).toBe(true)
  })

  it('skips the login for someone already logged in', () => {
    expect(resolveNavigation(admin, loginPage)).toBe(paths.adminHome)
  })
})
