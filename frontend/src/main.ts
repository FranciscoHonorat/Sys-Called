import { createApp } from 'vue'
import { createWebHistory } from 'vue-router'

import App from './App.vue'
import { authServiceKey } from './auth/authService'
import { createHttpAuthService } from './auth/httpAuthService'
import { restoreSession } from './auth/restoreSession'
import { createSession, sessionKey } from './auth/session'
import { createAppRouter } from './router'
import { createTicketsApi, ticketsApiKey } from './tickets/ticketsApi'
import './style.css'

const session = createSession()
const authService = createHttpAuthService()

await restoreSession(authService, session)

createApp(App)
  .use(createAppRouter(session, createWebHistory()))
  .provide(sessionKey, session)
  .provide(authServiceKey, authService)
  .provide(ticketsApiKey, createTicketsApi(authService))
  .mount('#app')
