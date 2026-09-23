import { inject, onMounted, ref } from 'vue'

import { useSessionExit } from '../auth/useSessionExit'
import { SessionExpiredError, ticketsApiKey, type Responsible, type TicketDetail } from './ticketsApi'

export function useTicketDetail(id: string) {
  const api = inject(ticketsApiKey)!
  const { expire } = useSessionExit()

  const ticket = ref<TicketDetail | null>(null)
  const responsibles = ref<Responsible[]>([])
  const notFound = ref(false)
  const actionError = ref('')

  async function handleFailure(error: unknown, explain: (message: string) => void) {
    if (error instanceof SessionExpiredError) {
      await expire()
      return
    }
    explain(error instanceof Error ? error.message : String(error))
  }

  async function perform(action: () => Promise<unknown>) {
    try {
      await action()
      actionError.value = ''
      ticket.value = await api.get(id)
    } catch (error) {
      await handleFailure(error, (message) => {
        actionError.value = `Não foi possível concluir a ação: ${message}`
      })
    }
  }

  onMounted(async () => {
    try {
      const [loaded, list] = await Promise.all([api.get(id), api.responsibles()])
      ticket.value = loaded
      responsibles.value = list
    } catch (error) {
      await handleFailure(error, () => {
        notFound.value = true
      })
    }
  })

  return { api, ticket, responsibles, notFound, actionError, perform }
}
