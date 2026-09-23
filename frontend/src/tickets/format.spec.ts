import { formatDuration } from './format'

describe('formatDuration', () => {
  const start = '2026-09-22T10:00:00Z'

  it.each([
    ['2026-09-22T10:00:30Z', 'menos de 1 min'],
    ['2026-09-22T10:45:00Z', '45 min'],
    ['2026-09-22T12:15:00Z', '2 h 15 min'],
    ['2026-09-22T13:00:00Z', '3 h'],
    ['2026-09-25T14:30:00Z', '3 dias e 4 h'],
    ['2026-09-23T10:00:00Z', '1 dia'],
  ])('from 10:00 to %s is %s', (end, expected) => {
    expect(formatDuration(start, end)).toBe(expected)
  })
})
