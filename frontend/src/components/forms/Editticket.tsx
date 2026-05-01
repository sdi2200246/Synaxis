import { useState } from 'react'
import { FiX } from 'react-icons/fi'
import { updateTicketType } from '../../api/tickets'
import type { TicketType } from '../../api/tickets'

interface Props {
  ticket: TicketType
  onClose: () => void
  onSuccess: () => void
}

export function EditTicketForm({ ticket, onClose, onSuccess }: Props) {
  const [form, setForm] = useState({
    name: ticket.name,
    price: String(ticket.price),
    quantity: String(ticket.quantity),
  })
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await updateTicketType(ticket.event_id, ticket.id ,{
        name: form.name !== ticket.name ? form.name : undefined,
        price: Number(form.price) !== ticket.price ? Number(form.price) : undefined,
        quantity: Number(form.quantity) !== ticket.quantity ? Number(form.quantity) : undefined,
      })
      onSuccess()
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update ticket type')
    } finally {
      setSubmitting(false)
    }
  }


   return (
    <div className="overlay" onClick={onClose}>
      <div
        className="dialog dialog--wide"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="dialog__header">
          <h2>Edit Ticket Type</h2>

          <button
            className="btn btn--icon"
            onClick={onClose}
            type="button"
          >
            <FiX size={24} />
          </button>
        </div>

        <div className="dialog__body">
          <form
            onSubmit={handleSubmit}
            className="card ticket-create-form"
          >
            {error && (
              <div className="alert alert--error">
                {error}
              </div>
            )}

            <div className="ticket-create-form__row">
              <div className="field">
                <label className="field__label">Name</label>
                <input
                  className="field__control"
                  value={form.name}
                  onChange={(e) =>
                    setForm({
                      ...form,
                      name: e.target.value,
                    })
                  }
                  required
                  disabled={submitting}
                />
              </div>

              <div className="field">
                <label className="field__label">
                  Price (€)
                </label>

                <input
                  className="field__control"
                  type="number"
                  min="0"
                  step="0.01"
                  value={form.price}
                  onChange={(e) =>
                    setForm({
                      ...form,
                      price: e.target.value,
                    })
                  }
                  required
                  disabled={submitting}
                />
              </div>

              <div className="field">
                <label className="field__label">
                  Quantity
                </label>

                <input
                  className="field__control"
                  type="number"
                  min="1"
                  value={form.quantity}
                  onChange={(e) =>
                    setForm({
                      ...form,
                      quantity: e.target.value,
                    })
                  }
                  required
                  disabled={submitting}
                />
              </div>
            </div>

            <div className="dialog__actions dialog__actions--with-divider">
              <button
                className="btn btn--ghost"
                type="button"
                onClick={onClose}
              >
                Cancel
              </button>

              <button
                className="btn btn--primary"
                type="submit"
                disabled={submitting}
              >
                {submitting
                  ? "Saving..."
                  : "Save"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}