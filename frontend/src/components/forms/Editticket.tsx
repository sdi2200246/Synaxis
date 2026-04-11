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
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Edit Ticket Type</h2>
          <button className="close-btn" onClick={onClose} type="button">
            <FiX size={24} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="tc-form">
          {error && <div className="error-message">{error}</div>}
          <div className="tc-form-fields">
            <div className="tc-field">
              <label>Name</label>
              <input
                value={form.name}
                onChange={e => setForm({ ...form, name: e.target.value })}
                required
                disabled={submitting}
              />
            </div>
            <div className="tc-field">
              <label>Price (€)</label>
              <input
                type="number"
                min="0"
                step="0.01"
                value={form.price}
                onChange={e => setForm({ ...form, price: e.target.value })}
                required
                disabled={submitting}
              />
            </div>
            <div className="tc-field">
              <label>Quantity</label>
              <input
                type="number"
                min="1"
                value={form.quantity}
                onChange={e => setForm({ ...form, quantity: e.target.value })}
                required
                disabled={submitting}
              />
            </div>
            <div className="form-actions">
              <button className="btn-cancel" type="button" onClick={onClose}>Cancel</button>
              <button className="btn-submit" type="submit" disabled={submitting}>
                {submitting ? 'Saving...' : 'Save'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}