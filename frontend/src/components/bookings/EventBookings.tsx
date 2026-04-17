import { FiMail, FiPhone, FiTag } from 'react-icons/fi'
import type { EventBooking } from '../../api/bookings'
import './Bookings.css'

interface EventBookingCardProps {
  booking: EventBooking
}

export function EventBookingCard({ booking }: EventBookingCardProps) {
  return (
    <div className="eb-card">
      <div className="eb-header">
        <div className="eb-attendee">
          <span className="eb-attendee__name">{booking.attendee_name}</span>
          <div className="eb-attendee__contact">
            <span><FiMail size={12} />{booking.attendee_email}</span>
            {booking.attendee_phone && (
              <span><FiPhone size={12} />{booking.attendee_phone}</span>
            )}
          </div>
        </div>
        <span className="eb-date">
          {new Date(booking.booked_at).toLocaleDateString('en-US', { dateStyle: 'medium' })}
        </span>
      </div>

      <div className="eb-ticket">
        <div className="eb-ticket-detail">
          <span className="eb-label">Ticket</span>
          <span className="eb-value">{booking.ticket_name}</span>
        </div>
        <div className="eb-ticket-detail">
          <span className="eb-label">Qty</span>
          <span className="eb-value">{booking.number_of_tickets}</span>
        </div>
        <div className="eb-ticket-detail">
          <span className="eb-label">Total</span>
          <span className="eb-value">€{booking.total_cost.toFixed(2)}</span>
        </div>
      </div>
    </div>
  )
}


