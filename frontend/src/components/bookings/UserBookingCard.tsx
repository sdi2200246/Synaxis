import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { FiMapPin, FiCalendar, FiMessageSquare } from 'react-icons/fi'
import { EventMap } from '../events/Map'
import { createConversation } from '../../api/messages'
import type { UserBooking } from '../../api/bookings'
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
 
        <div className="ub-card-actions">
          {hasCoords && (
            <button className="ub-map-toggle" onClick={() => setShowMap(!showMap)}>
              <FiMapPin size={14} />
              {showMap ? 'Hide map' : 'View on map'}
            </button>
          )}
          <button className="ub-msg-btn" onClick={handleMessageClick}>
            <FiMessageSquare size={14} />
            Message organizer
          </button>
        </div>
 
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
 
      {/* Confirmation dialog */}
      {showDialog && (
        <div className="ub-dialog-overlay" onClick={() => !creating && setShowDialog(false)}>
          <div className="ub-dialog" onClick={e => e.stopPropagation()}>
            <h3 className="ub-dialog-title">Start a conversation?</h3>
            <p className="ub-dialog-body">
              This will open a direct message thread with the organizer of{' '}
              <strong>{booking.event_title}</strong>. You can use it to ask questions
              about your booking.
            </p>
            {error && <p className="ub-dialog-error">{error}</p>}
            <div className="ub-dialog-actions">
              <button
                className="ub-dialog-cancel"
                onClick={() => setShowDialog(false)}
                disabled={creating}
              >
                Cancel
              </button>
              <button
                className="ub-dialog-confirm"
                onClick={handleCreateConversation}
                disabled={creating}
              >
                {creating ? 'Starting…' : 'Start conversation'}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
 