import { createRouter, type RouterHistory } from 'vue-router'

import type { Role } from '../auth/authService'
import type { Session } from '../auth/session'
import ChangePasswordView from '../views/ChangePasswordView.vue'
import HomeView from '../views/HomeView.vue'
import LoginPage from '../views/LoginPage.vue'
import RecoverPasswordView from '../views/RecoverPasswordView.vue'
import SignUpView from '../views/SignUpView.vue'
import TicketDetailView from '../views/TicketDetailView.vue'
import SupportsView from '../views/SupportsView.vue'
import TicketsView from '../views/TicketsView.vue'
import UsersView from '../views/UsersView.vue'
import { resolveNavigation } from './navigationGuard'
import { paths } from './paths'

declare module 'vue-router' {
  interface RouteMeta {
    role?: Role
    requiresAuth?: boolean
  }
}

export function createAppRouter(session: Session, history: RouterHistory) {
  const router = createRouter({
    history,
    routes: [
      { path: paths.login, component: LoginPage },
      { path: paths.register, component: SignUpView },
      { path: paths.recoverPassword, component: RecoverPasswordView },
      { path: paths.changePassword, component: ChangePasswordView, meta: { requiresAuth: true } },
      { path: paths.userHome, component: HomeView, meta: { role: 'user' } },
      { path: paths.supportHome, component: HomeView, meta: { role: 'support' } },
      { path: paths.adminHome, component: HomeView, meta: { role: 'admin' } },
      { path: paths.tickets, component: TicketsView, meta: { requiresAuth: true } },
      { path: paths.ticketDetail, component: TicketDetailView, props: true, meta: { requiresAuth: true } },
      { path: paths.users, component: UsersView, meta: { role: 'admin' } },
      { path: paths.supports, component: SupportsView, meta: { role: 'admin' } },
    ],
  })

  router.beforeEach((to) => resolveNavigation(session.user.value, to))

  return router
}
