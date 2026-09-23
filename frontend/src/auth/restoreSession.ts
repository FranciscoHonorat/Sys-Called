import type { AuthService } from './authService'
import type { Session } from './session'

export async function restoreSession(authService: AuthService, session: Session): Promise<void> {
  const user = await authService.restore().catch(() => null)
  if (user) {
    session.start(user)
  }
}
