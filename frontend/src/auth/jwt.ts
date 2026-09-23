export function decodeJwtClaims(token: string): Record<string, unknown> {
  const payload = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
  const bytes = Uint8Array.from(atob(payload), (char) => char.charCodeAt(0))
  return JSON.parse(new TextDecoder().decode(bytes))
}
