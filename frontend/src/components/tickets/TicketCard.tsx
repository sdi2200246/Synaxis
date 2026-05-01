import type { TicketType } from '../../api/tickets'
import './Tickets.css'

interface Props {
  ticket: TicketType
  onEdit: (ticket: TicketType) => void
}

export function TicketCard({ ticket, onEdit }: Props) {
  const soldOut = ticket.available === 0
  const sold = ticket.quantity - ticket.available

  return (
    <div className={`card card--hoverable ticket-type-card ${soldOut ? 'is-sold-out' : ''}`}>
      <div className="ticket-type-card__header">
        <span className="card__title">{ticket.name}</span>
        <span className="ticket-type-card__price">€{ticket.price.toFixed(2)}</span>
        <button className="btn btn--ghost" onClick={() => onEdit(ticket)}>Edit</button>
      </div>
      <div className="ticket-type-card__stats">
        <div className="ticket-type-card__stat">
          <span className="label">Total</span>
          <span className="ticket-type-card__stat-value">{ticket.quantity}</span>
        </div>
        <div className="ticket-type-card__stat">
          <span className="label">Sold</span>
          <span className="ticket-type-card__stat-value">{sold}</span>
        </div>
        <div className="ticket-type-card__stat">
          <span className="label">Available</span>
          <span className={`ticket-type-card__stat-value ${soldOut ? 'is-zero' : ''}`}>
            {ticket.available}
          </span>
        </div>
      </div>
    </div>
  )
}