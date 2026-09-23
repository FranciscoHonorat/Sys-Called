import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'

import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey, type TicketsApi } from '../tickets/ticketsApi'
import NewTicketForm from './NewTicketForm.vue'

function renderForm(api: TicketsApi = fakeTicketsApi()) {
  const rendered = render(NewTicketForm, { global: { provide: { [ticketsApiKey]: api } } })
  return { api, ...rendered }
}

async function fillAndSubmit() {
  await userEvent.type(screen.getByLabelText('Título'), 'Impressora')
  await userEvent.type(screen.getByLabelText('Descrição'), 'Não imprime')
  await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))
}

describe('NewTicketForm', () => {
  it('opens a ticket with title, description and priority and reports its id', async () => {
    const { api, emitted } = renderForm(fakeTicketsApi({ open: vi.fn().mockResolvedValue('t-1') }))

    await userEvent.selectOptions(screen.getByLabelText('Prioridade'), 'Alta')
    await fillAndSubmit()

    expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'High' })
    expect(emitted().opened).toEqual([['t-1']])
  })

  it('starts with medium priority', () => {
    renderForm()

    expect(screen.getByLabelText('Prioridade')).toHaveValue('Medium')
  })

  it('requires a title and a description', async () => {
    const { api } = renderForm()

    await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))

    expect(screen.getByText('Informe o título')).toBeInTheDocument()
    expect(screen.getByText('Informe a descrição')).toBeInTheDocument()
    expect(api.open).not.toHaveBeenCalled()
  })

  it('explains when the ticket cannot be opened', async () => {
    const { emitted } = renderForm(fakeTicketsApi({ open: vi.fn().mockRejectedValue(new Error('internal error')) }))

    await fillAndSubmit()

    expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível abrir o chamado')
    expect(emitted().opened).toBeUndefined()
  })
})
