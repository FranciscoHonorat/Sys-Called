export function utf8Base64Url(value: object): string {
  const bytes = new TextEncoder().encode(JSON.stringify(value))
  return btoa(String.fromCharCode(...bytes)).replace(/=+$/, '').replace(/\+/g, '-').replace(/\//g, '_')
}

export function fakeJwt(claims: object): string {
  return `${utf8Base64Url({ alg: 'EdDSA', typ: 'JWT' })}.${utf8Base64Url(claims)}.signature`
}
