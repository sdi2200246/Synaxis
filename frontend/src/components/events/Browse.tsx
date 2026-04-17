import { useState, useEffect } from 'react'
import { FiMapPin, FiCalendar, FiTag } from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'
import { EventMap } from './Map'
import { getTicketTypes } from '../../api/tickets'
import type { Event } from '../../types'
import type { TicketType } from '../../api/tickets'
import { createBooking } from '../../api/bookings'
import './Events.css'

interface BrowseEventCardProps {
  event: Event
}

export function BrowseEventCard({ event }: BrowseEventCardProps) {
  const [open, setOpen] = useState(false)
  const [showMap, setShowMap] = useState(false)
  const { isAuthenticated } = useAuth()

  const [tickets, setTickets] = useState<TicketType[]>([])
  const [loadingTickets, setLoadingTickets] = useState(false)

  const [selectedTicket, setSelectedTicket] = useState<TicketType | null>(null)
  const [quantity, setQuantity] = useState(1)
  const [confirming, setConfirming] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [bookingError, setBookingError] = useState('')
  const [bookingSuccess, setBookingSuccess] = useState('')

  const hasImage = false
  const hasCoords = event.venue.latitude != null && event.venue.longitude != null

  const dateLabel = new Date(event.start_datetime).toLocaleDateString('en-US', {
    dateStyle: 'medium',
  })

  useEffect(() => {
    if (!open) return
    setLoadingTickets(true)
    getTicketTypes(event.id)
      .then(setTickets)
      .catch(() => {})
      .finally(() => setLoadingTickets(false))
  }, [open, event.id])

  function handleSelectTicket(t: TicketType) {
    setSelectedTicket(t)
    setQuantity(1)
    setBookingError('')
    setBookingSuccess('')
  }

  function handleRequestBook() {
    setConfirming(true)
  }

  async function handleConfirmBook() {
    if (!selectedTicket) return
    setSubmitting(true)
    setBookingError('')

    try {
      await createBooking(event.id, selectedTicket.id, quantity)
      setBookingSuccess(`Booked ${quantity} × ${selectedTicket.name}!`)
      setSelectedTicket(null)
      setConfirming(false)
      setQuantity(1)

      // Refresh ticket availability
      const updated = await getTicketTypes(event.id)
      setTickets(updated)
    } catch (err: any) {
      const msg = err.response?.data?.error || 'Booking failed'
      setBookingError(msg)
      setConfirming(false)
    } finally {
      setSubmitting(false)
    }
  }

  function handleCancelBook() {
    setConfirming(false)
  }

  function handleClose() {
    setOpen(false)
    setSelectedTicket(null)
    setConfirming(false)
    setBookingError('')
    setBookingSuccess('')
    setQuantity(1)
  }

  const totalCost = selectedTicket ? selectedTicket.price * quantity : 0

  return (
    <>
      <div className="browse-card" onClick={() => setOpen(true)}>
        {hasImage ? (
          <img className="browse-card__img" src="" alt={event.title} />
        ) : (
          <div className="browse-card__placeholder">
            <FiCalendar size={28} />
          </div>
        )}
        <div className="browse-card__body">
          <span className="browse-card__title">{event.title}</span>
          <div className="browse-card__meta">
            <span>{dateLabel}</span>
            <span className="browse-card__dot" />
            <span>{event.venue.city}</span>
          </div>
        </div>
      </div>

      {open && (
        <div className="browse-detail-overlay" onClick={handleClose}>
          <div className="browse-detail" onClick={e => e.stopPropagation()}>
            <button className="browse-detail__close" onClick={handleClose}>
              &times;
            </button>

            {hasImage && (
              <img className="browse-detail__img" src="" alt={event.title} />
            )}

            <div className="browse-detail__content">
              <h2 className="browse-detail__title">{event.title}</h2>

              <div className="browse-detail__info">
                <span><FiTag size={14} />{event.event_type}</span>
                <span><FiMapPin size={14} />{event.venue.name}, {event.venue.city}</span>
                <span><FiCalendar size={14} />{dateLabel}</span>
              </div>

              {hasCoords && (
                <button
                  className="browse-detail__map-toggle"
                  onClick={() => setShowMap(!showMap)}
                >
                  <FiMapPin size={14} />
                  {showMap ? 'Hide map' : 'View on map'}
                </button>
              )}

              {showMap && hasCoords && (
                <EventMap
                  lat={event.venue.latitude!}
                  lng={event.venue.longitude!}
                  venueName={event.venue.name}
                />
              )}

              <p className="browse-detail__desc">{event.description}</p>

              {bookingSuccess && (
                <div className="browse-detail__success">{bookingSuccess}</div>
              )}

              {bookingError && (
                <div className="browse-detail__error">{bookingError}</div>
              )}

              <div className="browse-detail__tickets">
                <h3 className="browse-detail__section-title">Tickets</h3>
                {loadingTickets ? (
                  <p className="browse-detail__tickets-loading">Loading tickets…</p>
                ) : tickets.length === 0 ? (
                  <p className="browse-detail__tickets-empty">No tickets available.</p>
                ) : (
                  <div className="browse-detail__ticket-list">
                    {tickets.map(t => (
                      <div
                        key={t.id}
                        className={`browse-detail__ticket ${selectedTicket?.id === t.id ? 'browse-detail__ticket--selected' : ''}`}
                      >
                        <div className="browse-detail__ticket-info">
                          <span className="browse-detail__ticket-name">{t.name}</span>
                          <span className="browse-detail__ticket-avail">
                            {t.available > 0 ? `${t.available} left` : 'Sold out'}
                          </span>
                        </div>
                        <div className="browse-detail__ticket-right">
                          <span className="browse-detail__ticket-price">
                            {t.price === 0 ? 'Free' : `€${t.price.toFixed(2)}`}
                          </span>
                          {isAuthenticated && t.available > 0 && (
                            <button
                              className="browse-detail__btn browse-detail__btn--primary"
                              onClick={() => handleSelectTicket(t)}
                            >
                              {selectedTicket?.id === t.id ? 'Selected' : 'Select'}
                            </button>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {selectedTicket && !confirming && (
                <div className="browse-detail__booking-form">
                  <div className="browse-detail__qty-row">
                    <label className="browse-detail__qty-label">Quantity</label>
                    <div className="browse-detail__qty-controls">
                      <button
                        className="browse-detail__qty-btn"
                        onClick={() => setQuantity(q => Math.max(1, q - 1))}
                        disabled={quantity <= 1}
                      >−</button>
                      <span className="browse-detail__qty-value">{quantity}</span>
                      <button
                        className="browse-detail__qty-btn"
                        onClick={() => setQuantity(q => Math.min(selectedTicket.available, q + 1))}
                        disabled={quantity >= selectedTicket.available}
                      >+</button>
                    </div>
                  </div>
                  <div className="browse-detail__total-row">
                    <span className="browse-detail__total-label">Total</span>
                    <span className="browse-detail__total-value">
                      {totalCost === 0 ? 'Free' : `€${totalCost.toFixed(2)}`}
                    </span>
                  </div>
                  <button
                    className="browse-detail__book-btn"
                    onClick={handleRequestBook}
                  >
                    Book {quantity} ticket{quantity > 1 ? 's' : ''}
                  </button>
                </div>
              )}

              {confirming && selectedTicket && (
                <div className="browse-detail__confirm">
                  <p className="browse-detail__confirm-text">
                    Confirm booking of <strong>{quantity} × {selectedTicket.name}</strong> for{' '}
                    <strong>{totalCost === 0 ? 'Free' : `€${totalCost.toFixed(2)}`}</strong>?
                  </p>
                  <p className="browse-detail__confirm-warning">
                    This action cannot be undone.
                  </p>
                  <div className="browse-detail__confirm-actions">
                    <button
                      className="browse-detail__btn"
                      onClick={handleCancelBook}
                      disabled={submitting}
                    >
                      Cancel
                    </button>
                    <button
                      className="browse-detail__confirm-btn"
                      onClick={handleConfirmBook}
                      disabled={submitting}
                    >
                      {submitting ? 'Booking…' : 'Confirm'}
                    </button>
                  </div>
                </div>
              )}

              {!isAuthenticated && (
                <div className="browse-detail__guest-note">
                  Sign up and get approved to book tickets and contact organizers.
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  )
}