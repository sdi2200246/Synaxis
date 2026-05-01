import { FiMail, FiPhone} from 'react-icons/fi'
import type { EventBooking } from '../../api/bookings'
import './Bookings.css'

interface EventBookingCardProps {
  booking: EventBooking
}

export function EventBookingCard({ booking }: EventBookingCardProps) {
  return (
    <div className="card event-booking-card">
      <div className="event-booking-card__header">
        <div className="event-booking-card__attendee">
          <span className="event-booking-card__name">{booking.attendee_name}</span>
          <div className="event-booking-card__contact">
            <span><FiMail size={12} />{booking.attendee_email}</span>
            {booking.attendee_phone && (
              <span><FiPhone size={12} />{booking.attendee_phone}</span>
            )}
          </div>
        </div>
        <span className="event-booking-card__date">
          {new Date(booking.booked_at).toLocaleDateString('en-US', { dateStyle: 'medium' })}
        </span>
      </div>

      <div className="event-booking-card__strip">
        <div className="event-booking-card__strip-item">
          <span className="label">Ticket</span>
          <span className="event-booking-card__strip-value">{booking.ticket_name}</span>
        </div>
        <div className="event-booking-card__strip-item">
          <span className="label">Qty</span>
          <span className="event-booking-card__strip-value">{booking.number_of_tickets}</span>
        </div>
        <div className="event-booking-card__strip-item">
          <span className="label">Total</span>
          <span className="event-booking-card__strip-value">€{booking.total_cost.toFixed(2)}</span>
        </div>
      </div>
    </div>
  )
}


