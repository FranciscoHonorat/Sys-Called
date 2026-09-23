import { homePathFor } from './homePath'

describe('homePathFor', () => {
  it.each([
    ['user', '/usuario'],
    ['support', '/suporte'],
    ['admin', '/admin'],
  ] as const)('sends %s to %s', (role, path) => {
    expect(homePathFor(role)).toBe(path)
  })
})
