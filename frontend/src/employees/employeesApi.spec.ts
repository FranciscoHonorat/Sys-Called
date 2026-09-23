import { createEmployeesApi } from './employeesApi'

describe('employeesApi', () => {
  it('lists the employees with the access token', async () => {
    const employees = [{ id: 'agent-1', name: 'Ana Souza', username: 'ana', role: 'support' }]
    const fetchFn = vi.fn().mockResolvedValue(
      new Response(JSON.stringify(employees), { status: 200, headers: { 'Content-Type': 'application/json' } }),
    )
    const api = createEmployeesApi({ accessToken: () => 'token-1', restore: vi.fn() }, fetchFn)

    expect(await api.list()).toEqual(employees)
    expect(fetchFn).toHaveBeenCalledWith('/api/employees/employees', {
      method: 'GET',
      headers: { Authorization: 'Bearer token-1', 'Content-Type': 'application/json' },
    })
  })

  it('approves a pending account', async () => {
    const fetchFn = vi.fn().mockResolvedValue(new Response(null, { status: 204 }))
    const api = createEmployeesApi({ accessToken: () => 'token-1', restore: vi.fn() }, fetchFn)

    await api.approve('user-9')

    expect(fetchFn).toHaveBeenCalledWith('/api/employees/employees/user-9/approve', {
      method: 'POST',
      headers: { Authorization: 'Bearer token-1', 'Content-Type': 'application/json' },
    })
  })

  it('issues a temporary password', async () => {
    const fetchFn = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ temporary_password: 'Temp-1234' }), { status: 200 }),
    )
    const api = createEmployeesApi({ accessToken: () => 'token-1', restore: vi.fn() }, fetchFn)

    expect(await api.issueTemporaryPassword('user-1')).toBe('Temp-1234')
    expect(fetchFn).toHaveBeenCalledWith('/api/employees/employees/user-1/temporary-password', {
      method: 'POST',
      headers: { Authorization: 'Bearer token-1', 'Content-Type': 'application/json' },
    })
  })
})
