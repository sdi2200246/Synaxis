import { useState } from 'react'
import { FiMapPin, FiCalendar, FiTag } from 'react-icons/fi'
import { EventMap } from '../events/Map'
import type { UserBooking } from '../../api/bookings'
import './Bookings.css'

interface UserBookingCardProps {
  booking: UserBooking
}

export function UserBookingCard({ booking }: UserBookingCardProps) {
  const [showMap, setShowMap] = useState(false)

  const hasCoords = booking.venue_latitude != null && booking.venue_longitude != null

  const dateLabel = new Date(booking.event_start).toLocaleDateString('en-US', {
    dateStyle: 'medium',
  })

  const timeLabel = new Date(booking.event_start).toLocaleTimeString('en-US', {
    timeStyle: 'short',
  })

  return (
    <div className="ub-card">
      <div className="ub-header">
        <div className="ub-event">
          <span className="ub-title">{booking.event_title}</span>
          <div className="ub-meta">
            <span><FiMapPin size={13} />{booking.venue_name}, {booking.venue_city}</span>
            <span><FiCalendar size={13} />{dateLabel} · {timeLabel}</span>
          </div>
        </div>
        <span className={`ub-status ub-status--${booking.status.toLowerCase()}`}>
          {booking.status}
        </span>
      </div>

      <div className="ub-ticket">
        <div className="ub-ticket-detail">
          <span className="ub-label">Ticket</span>
          <span className="ub-value">{booking.ticket_name}</span>
        </div>
        <div className="ub-ticket-detail">
          <span className="ub-label">Qty</span>
          <span className="ub-value">{booking.number_of_tickets}</span>
        </div>
        <div className="ub-ticket-detail">
          <span className="ub-label">Total</span>
          <span className="ub-value">€{booking.total_cost.toFixed(2)}</span>
        </div>
        <div className="ub-ticket-detail">
          <span className="ub-label">Booked</span>
          <span className="ub-value">
            {new Date(booking.booked_at).toLocaleDateString('en-US', { dateStyle: 'medium' })}
          </span>
        </div>
      </div>

      {hasCoords && (
        <button
          className="ub-map-toggle"
          onClick={() => setShowMap(!showMap)}
        >
          <FiMapPin size={14} />
          {showMap ? 'Hide map' : 'View on map'}
        </button>
      )}

      {showMap && hasCoords && (
        <div className="ub-map">
          <EventMap
            lat={booking.venue_latitude!}
            lng={booking.venue_longitude!}
            venueName={booking.venue_name}
          />
        </div>
      )}
    </div>
  )
}