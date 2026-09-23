import { createRouter, type RouterHistory } from 'vue-router'

import type { Role } from '../auth/authService'
import type { Session } from '../auth/session'
import ComingSoonView from '../views/ComingSoonView.vue'
import HomeView from '../views/HomeView.vue'
import LoginPage from '../views/LoginPage.vue'
import NewTicketView from '../views/NewTicketView.vue'
import TicketDetailView from '../views/TicketDetailView.vue'
import TicketsView from '../views/TicketsView.vue'
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
      { path: paths.register, component: ComingSoonView },
      { path: paths.recoverPassword, component: ComingSoonView },
      { path: paths.userHome, component: HomeView, meta: { role: 'user' } },
      { path: paths.supportHome, component: HomeView, meta: { role: 'support' } },
      { path: paths.adminHome, component: HomeView, meta: { role: 'admin' } },
      { path: paths.tickets, component: TicketsView, meta: { requiresAuth: true } },
      { path: paths.newTicket, component: NewTicketView, meta: { requiresAuth: true } },
      { path: paths.ticketDetail, component: TicketDetailView, props: true, meta: { requiresAuth: true } },
      { path: paths.users, component: ComingSoonView, meta: { role: 'admin' } },
      { path: paths.supports, component: ComingSoonView, meta: { role: 'admin' } },
    ],
  })

  router.beforeEach((to) => resolveNavigation(session.user.value, to))

  return router
}
