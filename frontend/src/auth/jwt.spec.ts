import { fakeJwt } from '../test/fakeJwt'
import { decodeJwtClaims } from './jwt'

describe('decodeJwtClaims', () => {
  it('reads names with accents exactly as the server sent them', () => {
    const token = fakeJwt({ sub: 'user-1', name: 'Usuário Padrão', role: 'user' })

    expect(decodeJwtClaims(token)).toEqual({ sub: 'user-1', name: 'Usuário Padrão', role: 'user' })
  })
})
