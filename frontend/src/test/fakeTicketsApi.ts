import { vi } from 'vitest'

import type { TicketsApi } from '../tickets/ticketsApi'

export function fakeTicketsApi(overrides: Partial<TicketsApi> = {}): TicketsApi {
  return {
    list: vi.fn().mockResolvedValue([]),
    open: vi.fn().mockResolvedValue('t-1'),
    get: vi.fn(),
    responsibles: vi.fn().mockResolvedValue([]),
    edit: vi.fn().mockResolvedValue(undefined),
    assign: vi.fn().mockResolvedValue(undefined),
    autoAssign: vi.fn().mockResolvedValue('agent-1'),
    changePriority: vi.fn().mockResolvedValue(undefined),
    start: vi.fn().mockResolvedValue(undefined),
    close: vi.fn().mockResolvedValue(undefined),
    respond: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }
}
