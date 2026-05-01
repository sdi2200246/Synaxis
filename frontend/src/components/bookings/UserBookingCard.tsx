import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { FiMapPin, FiCalendar, FiMessageSquare } from 'react-icons/fi'
import { EventMap } from '../events/Map'
import { createConversation } from '../../api/messages'
import type { UserBooking } from '../../api/bookings'
import { ConfirmDialog } from '../ConfirmDialogue'
import './Bookings.css'

interface UserBookingCardProps {
  booking: UserBooking
  conversationId: string | null
  currentUserId: string
  onConversationCreated: (bookingId: string, conversationId: string) => void
}

export function UserBookingCard({
  booking,
  conversationId,
  currentUserId,
  onConversationCreated,
}: UserBookingCardProps) {
  const navigate = useNavigate()
  const [showMap, setShowMap] = useState(false)
  const [showDialog, setShowDialog] = useState(false)
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')

  const hasCoords = booking.venue_latitude != null && booking.venue_longitude != null

  const dateLabel = new Date(booking.event_start).toLocaleDateString('en-US', { dateStyle: 'medium' })
  const timeLabel = new Date(booking.event_start).toLocaleTimeString('en-US', { timeStyle: 'short' })

  function handleMessageClick() {
    if (conversationId) {
      navigate(`/messages/${conversationId}`)
    } else {
      setError('')
      setShowDialog(true)
    }
  }

  async function handleCreateConversation() {
    setCreating(true)
    setError('')
    try {
      const convId = await createConversation(
        booking.id,
        booking.organizer_id,
        currentUserId
      )
      onConversationCreated(booking.id, convId)
      setShowDialog(false)
      navigate(`/messages/${convId}`)
    } catch {
      setError('Failed to start conversation. Please try again.')
    } finally {
      setCreating(false)
    }
  }

  return (
    <>
      <div className="card user-booking-card">
        <div className="user-booking-card__header">
          <div className="user-booking-card__event">
            <span className="user-booking-card__title">{booking.event_title}</span>
            <div className="user-booking-card__meta">
              <span><FiMapPin size={13} />{booking.venue_name}, {booking.venue_city}</span>
              <span><FiCalendar size={13} />{dateLabel} · {timeLabel}</span>
            </div>
          </div>
          <span className={`badge badge--${booking.status.toLowerCase()}`}>
            {booking.status}
          </span>
        </div>

        <div className="user-booking-card__strip">
          <div className="user-booking-card__strip-item">
            <span className="label">Ticket</span>
            <span className="user-booking-card__strip-value">{booking.ticket_name}</span>
          </div>
          <div className="user-booking-card__strip-item">
            <span className="label">Qty</span>
            <span className="user-booking-card__strip-value">{booking.number_of_tickets}</span>
          </div>
          <div className="user-booking-card__strip-item">
            <span className="label">Total</span>
            <span className="user-booking-card__strip-value">€{booking.total_cost.toFixed(2)}</span>
          </div>
          <div className="user-booking-card__strip-item">
            <span className="label">Booked</span>
            <span className="user-booking-card__strip-value">
              {new Date(booking.booked_at).toLocaleDateString('en-US', { dateStyle: 'medium' })}
            </span>
          </div>
        </div>

        <div className="user-booking-card__actions">
          {hasCoords && (
            <button className="btn btn--soft btn--pill" onClick={() => setShowMap(!showMap)}>
              <FiMapPin size={14} />
              {showMap ? 'Hide map' : 'View on map'}
            </button>
          )}
          <button className="btn btn--soft btn--pill" onClick={handleMessageClick}>
            <FiMessageSquare size={14} />
            Message organizer
          </button>
        </div>

        {showMap && hasCoords && (
          <div className="user-booking-card__map">
            <EventMap
              lat={booking.venue_latitude!}
              lng={booking.venue_longitude!}
              venueName={booking.venue_name}
            />
          </div>
        )}
      </div>

      {showDialog && (
        <ConfirmDialog
          title="Start a conversation?"
          body={`This will open a direct message thread with the organizer of ${booking.event_title}. You can use it to ask questions about your booking.`}
          confirmLabel={creating ? 'Starting…' : 'Start conversation'}
          variant="primary"
          loading={creating}
          error={error || undefined}
          onConfirm={handleCreateConversation}
          onCancel={() => setShowDialog(false)}
        />
      )}
    </>
  )
}