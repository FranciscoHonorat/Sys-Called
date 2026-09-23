import { render, screen, waitFor } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'

import AppModal from './AppModal.vue'

function renderModal(open = true) {
  return render(AppModal, {
    props: { open, title: 'Abrir novo chamado' },
    slots: { default: '<input aria-label="Título" />' },
  })
}

describe('AppModal', () => {
  it('shows its content in a dialog named after its title', () => {
    renderModal()

    const dialog = screen.getByRole('dialog', { name: 'Abrir novo chamado' })
    expect(dialog).toHaveAttribute('aria-modal', 'true')
    expect(screen.getByLabelText('Título')).toBeInTheDocument()
  })

  it('renders nothing while closed', () => {
    renderModal(false)

    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })

  it.each([
    ['the back button', () => userEvent.click(screen.getByRole('button', { name: '← Voltar' }))],
    ['the Escape key', () => userEvent.keyboard('{Escape}')],
    ['a click outside', () => userEvent.click(screen.getByTestId('modal-backdrop'))],
  ])('asks to close through %s', async (_, close) => {
    const { emitted } = renderModal()

    await close()

    expect(emitted().close).toHaveLength(1)
  })

  it('does not close when clicking inside the dialog', async () => {
    const { emitted } = renderModal()

    await userEvent.click(screen.getByLabelText('Título'))

    expect(emitted().close).toBeUndefined()
  })

  it('moves the focus into the dialog', async () => {
    renderModal()

    await waitFor(() => expect(screen.getByRole('dialog')).toContainElement(document.activeElement as HTMLElement))
  })
})
