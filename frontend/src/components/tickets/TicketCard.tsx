import type { TicketType } from '../../api/tickets'
import './Tickets.css'

interface Props {
  ticket: TicketType
  onEdit: (ticket: TicketType) => void
}

export function TicketCard({ ticket , onEdit }: Props) {
  const soldOut = ticket.available === 0
  const sold = ticket.quantity - ticket.available

  return (
    <div className={`tc-card ${soldOut ? 'tc-card--sold-out' : ''}`}>
      <div className="tc-header">
        <span className="tc-name">{ticket.name}</span>
        <span className="tc-price">€{ticket.price.toFixed(2)}</span>
        <button className="tc-edit-btn" onClick={() => onEdit(ticket)}>Edit</button>
      </div>
      <div className="tc-stats">
        <div className="tc-stat">
          <span className="tc-stat-label">Total</span>
          <span className="tc-stat-value">{ticket.quantity}</span>
        </div>
        <div className="tc-stat">
          <span className="tc-stat-label">Sold</span>
          <span className="tc-stat-value">{sold}</span>
        </div>
        <div className="tc-stat">
          <span className="tc-stat-label">Available</span>
          <span className={`tc-stat-value ${soldOut ? 'tc-stat-value--zero' : ''}`}>
            {ticket.available}
          </span>
        </div>
      </div>
    </div>
  )
}