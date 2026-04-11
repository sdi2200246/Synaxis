import { useState, useEffect } from 'react'
import { useParams, Link , useLocation } from 'react-router-dom'
import { getTicketTypes, createTicketType } from '../api/tickets'
import type { TicketType } from '../api/tickets'
import { TicketCard } from '../components/tickets/TicketCard'
import { EditTicketForm } from '../components/forms/Editticket'
import '../components/tickets/Tickets.css'

export function EventTicketsPage() {
  const { id } = useParams<{ id: string }>()
  const [tickets, setTickets] = useState<TicketType[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({ name: '', price: '', quantity: '' })
  const [successMessage , setSuccessMessage] = useState('')
  const { state } = useLocation() as { state: { title: string; capacity: number } | null }
  const [editTarget, setEditTarget] = useState<TicketType | null>(null)

  async function fetchTickets() {
    if (!id) return
    try {
      const data = await getTicketTypes(id)
      setTickets(data)
    } catch {
      setError('Failed to load tickets')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchTickets() }, [id])

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) {
    setForm({ ...form, [e.target.name]: e.target.value })
    if (error) setError('')
  }


  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!id) return
    setSubmitting(true)
    setError('')

    try {
        await createTicketType(id, {
            name: form.name,
            price: Number(form.price),
            quantity: Number(form.quantity),
        })
        setForm({ name: '', price: '', quantity: '' })
        setSuccessMessage('Tickets Released successfully')
        fetchTickets()
        setTimeout(() => setSuccessMessage(''), 3000)

    } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to create ticket type')
    } finally {
        setSubmitting(false)
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <div>
          <Link to="/my-events" style={{ fontSize: '13px', color: '#888' }}>&larr; Back to My Events</Link>
          <h1>Ticket Types</h1>
        </div>
      </div>

    <h1>Tickets — {state?.title ?? 'Event'}</h1>
    <p style={{ color: '#888', fontSize: '13px' }}>Capacity: {state?.capacity}</p>  

      {error && <div className="error-message">{error}</div>}

    {successMessage && (
        <div className="toast">{successMessage}</div>
      )}

      <form className="tc-form" onSubmit={handleCreate}>
        <h3>Release a new ticket type</h3>
        <div className="tc-form-fields">
          <div className="tc-field">
            <label>Name</label>
            <input
              name="name"
              value={form.name}
              onChange={e => handleChange(e)}
              placeholder="e.g. General Admission"
              required
            />
          </div>
          <div className="tc-field">
            <label>Price (€)</label>
            <input
              name="price"
              type="number"
              min="0"
              step="0.01"
              value={form.price}
              onChange={e => handleChange(e)}
              required
            />
          </div>
          <div className="tc-field">
            <label>Quantity</label>
            <input
              name="quantity"
              type="number"
              min="1"
              value={form.quantity}
              onChange={e => handleChange(e)}
              required
            />
          </div>
          <button className="tc-submit" type="submit" disabled={submitting}>
            {submitting ? 'Creating...' : 'Release'}
          </button>
        </div>
      </form>

      {loading ? (
        <p>Loading...</p>
      ) : tickets.length === 0 ? (
        <div className="tc-empty">No ticket types yet — release one above.</div>
      ) : (
        <div className="tc-list">
         {tickets.map(t => (
              <TicketCard key={t.id} ticket={t} onEdit={t => setEditTarget(t)} />
            ))}
        </div>
      )}

        {editTarget && (
          <EditTicketForm
            ticket={editTarget}
            onClose={() => setEditTarget(null)}
            onSuccess={() => {
              setEditTarget(null)
              setSuccessMessage('Ticket updated successfully')
              fetchTickets()
              setTimeout(() => setSuccessMessage(''), 3000)
            }}
          />
        )}
    </div>
    
  )
}