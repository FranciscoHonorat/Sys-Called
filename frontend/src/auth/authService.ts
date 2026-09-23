import type { InjectionKey } from 'vue'

export type Role = 'user' | 'support' | 'admin'

export interface User {
  id: string
  name: string
  role: Role
}

export interface AuthService {
  login(username: string, password: string): Promise<User>
  restore(): Promise<User | null>
  logout(): Promise<void>
}

export const authServiceKey: InjectionKey<AuthService> = Symbol('authService')
