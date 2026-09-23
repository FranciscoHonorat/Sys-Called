import type { InjectionKey } from 'vue'

import { createApiClient, type TokenSource } from '../api/apiClient'
import type { Role } from '../auth/authService'

export type EmployeeStatus = 'active' | 'pending'

export interface Employee {
  id: string
  name: string
  username: string
  role: Role
  status: EmployeeStatus
  password_reset_requested: boolean
}

export interface EmployeesApi {
  list(): Promise<Employee[]>
  approve(id: string): Promise<void>
  issueTemporaryPassword(id: string): Promise<string>
}

export function createEmployeesApi(tokens: TokenSource, fetchFn: typeof fetch = fetch): EmployeesApi {
  const request = createApiClient('/api/employees', tokens, fetchFn)

  return {
    list: () => request<Employee[]>({ method: 'GET', path: '/employees' }),
    approve: (id) => request<void>({ method: 'POST', path: `/employees/${id}/approve` }),
    issueTemporaryPassword: async (id) => {
      const body = await request<{ temporary_password: string }>({ method: 'POST', path: `/employees/${id}/temporary-password` })
      return body.temporary_password
    },
  }
}

export const employeesApiKey: InjectionKey<EmployeesApi> = Symbol('employeesApi')
