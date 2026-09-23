import { vi } from 'vitest'

import type { AuthService } from '../auth/authService'

export function fakeAuthService(overrides: Partial<AuthService> = {}): AuthService {
  return {
    login: vi.fn(),
    restore: vi.fn().mockResolvedValue(null),
    logout: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }
}
