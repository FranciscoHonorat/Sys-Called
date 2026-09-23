import type { InjectionKey } from 'vue'

export type Role = 'user' | 'support' | 'admin'

export interface User {
  id: string
  name: string
  role: Role
  mustChangePassword?: boolean
}

export interface SignUpData {
  name: string
  username: string
  password: string
}

export class AccountPendingError extends Error {
  constructor() {
    super('account waiting for the administrator approval')
  }
}

export class AccountRequestError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export interface AuthService {
  login(username: string, password: string): Promise<User>
  restore(): Promise<User | null>
  logout(): Promise<void>
  signUp(data: SignUpData): Promise<void>
  requestPasswordReset(username: string): Promise<void>
  changePassword(current: string, next: string): Promise<User>
}

export const authServiceKey: InjectionKey<AuthService> = Symbol('authService')
