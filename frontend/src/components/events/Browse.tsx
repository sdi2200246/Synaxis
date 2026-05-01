import { useState, useEffect } from 'react'
import { FiMapPin, FiCalendar, FiTag, FiX } from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'
import { EventMap } from './Map'
import { getTicketTypes } from '../../api/tickets'
import type { Event } from '../../types'
import type { TicketType } from '../../api/tickets'
import { createBooking } from '../../api/bookings'
import { recordVisit } from '../../api/visits'
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
      .finally(() => setLoadingTickets(false))
    recordVisit(event.id).catch(() => {})
  }, [open, event.id])

  function handleSelectTicket(t: TicketType) {
    setSelectedTicket(t)
    setQuantity(1)
    setBookingError('')
    setBookingSuccess('')
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
      const updated = await getTicketTypes(event.id)
      setTickets(updated)
    } catch (err: any) {
      setBookingError(err.response?.data?.error || 'Booking failed')
      setConfirming(false)
    } finally {
      setSubmitting(false)
    }
  }

  function handleClose() {
    setOpen(false)
    setSelectedTicket(null)
    setConfirming(false)
    setBookingError('')
    setBookingSuccess('')
    setQuantity(1)
    setShowMap(false)
  }

  const totalCost = selectedTicket ? selectedTicket.price * quantity : 0

  return (
    <>
      {/* ── Card ─────────────────────────────────────────────────── */}
      <div className="card card--hoverable browse-card" onClick={() => setOpen(true)}>
        <div className="media media--16x10">
          {hasImage ? (
            <img className="media__img" src="" alt={event.title} />
          ) : (
            <div className="media__placeholder">
              <FiCalendar size={28} />
            </div>
          )}
        </div>
        <div className="browse-card__body">
          <span className="browse-card__title">{event.title}</span>
          <div className="browse-card__meta">
            <span>{dateLabel}</span>
            <span className="browse-card__dot" />
            <span>{event.venue.city}</span>
          </div>
        </div>
      </div>

      {/* ── Detail Dialog ────────────────────────────────────────── */}
      {open && (
        <div className="overlay" onClick={handleClose}>
          <div
            className="dialog dialog--medium browse-detail"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Hero */}
            <div className="browse-detail__hero">
              <div className="media media--16x7">
                {hasImage ? (
                  <img className="media__img" src="" alt={event.title} />
                ) : (
                  <div className="media__placeholder">
                    <FiCalendar size={42} />
                  </div>
                )}
              </div>
              <button
                type="button"
                className="browse-detail__close"
                onClick={handleClose}
                aria-label="Close"
              >
                <FiX size={18} />
              </button>
            </div>

            {/* Header */}
            <header className="section">
              <h2 className="browse-detail__title">{event.title}</h2>
              <div className="browse-detail__info">
                <span className="icon-text"><FiTag size={14} />{event.event_type}</span>
                <span className="icon-text"><FiMapPin size={14} />{event.venue.name}, {event.venue.city}</span>
                <span className="icon-text"><FiCalendar size={14} />{dateLabel}</span>
              </div>
            </header>

            {/* About */}
            <section className="section">
              <h3 className="label">About</h3>
              <p className="browse-detail__desc">{event.description}</p>
            </section>

            {/* Location */}
            {hasCoords && (
              <section className="section">
                <div className="section__head">
                  <h3 className="label">Location</h3>
                  <button
                    type="button"
                    className="btn btn--soft btn--pill"
                    onClick={() => setShowMap(!showMap)}
                  >
                    <FiMapPin size={14} />
                    {showMap ? 'Hide map' : 'Show map'}
                  </button>
                </div>
                {showMap && (
                  <EventMap
                    lat={event.venue.latitude!}
                    lng={event.venue.longitude!}
                    venueName={event.venue.name}
                  />
                )}
              </section>
            )}

            {/* Feedback */}
            {(bookingSuccess || bookingError) && (
              <section className="section">
                {bookingSuccess && <div className="alert alert--success">{bookingSuccess}</div>}
                {bookingError && <div className="alert alert--error">{bookingError}</div>}
              </section>
            )}

            {/* Tickets */}
            <section className="section">
              <h3 className="label">Tickets</h3>
              {loadingTickets ? (
                <p className="empty-state">Loading tickets…</p>
              ) : tickets.length === 0 ? (
                <p className="empty-state">No tickets available.</p>
              ) : (
                <div className="list-stack">
                  {tickets.map((t) => {
                    const selected = selectedTicket?.id === t.id
                    const soldOut = t.available <= 0
                    return (
                      <div
                        key={t.id}
                        className={`card ticket-row ${selected ? 'is-selected' : ''} ${soldOut ? 'is-sold-out' : ''}`}
                      >
                        <div className="ticket-row__main">
                          <span className="ticket-row__name">{t.name}</span>
                          <span className="ticket-row__avail">
                            {soldOut ? 'Sold out' : `${t.available} left`}
                          </span>
                        </div>
                        <div className="ticket-row__right">
                          <span className="ticket-row__price">
                            {t.price === 0 ? 'Free' : `€${t.price.toFixed(2)}`}
                          </span>
                          {isAuthenticated && !soldOut && (
                            <button
                              type="button"
                              className={`btn ${selected ? 'btn--primary' : 'btn--ghost'}`}
                              onClick={() => handleSelectTicket(t)}
                            >
                              {selected ? 'Selected' : 'Select'}
                            </button>
                          )}
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </section>

            {/* Booking form */}
            {selectedTicket && !confirming && (
              <section className="section">
                <div className="booking-form">
                  <div className="booking-form__qty-row">
                    <span className="label" style={{ marginBottom: 0 }}>Quantity</span>
                    <div className="qty-stepper">
                      <button
                        type="button"
                        className="qty-stepper__btn"
                        onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                        disabled={quantity <= 1}
                        aria-label="Decrease"
                      >
                        −
                      </button>
                      <span className="qty-stepper__value">{quantity}</span>
                      <button
                        type="button"
                        className="qty-stepper__btn"
                        onClick={() => setQuantity((q) => Math.min(selectedTicket.available, q + 1))}
                        disabled={quantity >= selectedTicket.available}
                        aria-label="Increase"
                      >
                        +
                      </button>
                    </div>
                  </div>

                  <div className="booking-form__total-row">
                    <span className="label" style={{ marginBottom: 0 }}>Total</span>
                    <span className="booking-form__total-value">
                      {totalCost === 0 ? 'Free' : `€${totalCost.toFixed(2)}`}
                    </span>
                  </div>

                  <button
                    type="button"
                    className="btn btn--primary btn--block"
                    onClick={() => setConfirming(true)}
                  >
                    Book {quantity} ticket{quantity > 1 ? 's' : ''}
                  </button>
                </div>
              </section>
            )}

            {/* Confirm */}
            {confirming && selectedTicket && (
              <section className="section section--compact">
                <h3 className="label">Confirm booking</h3>
                <p className="browse-detail__confirm-text">
                  You're booking <strong>{quantity} × {selectedTicket.name}</strong> for{' '}
                  <strong>{totalCost === 0 ? 'Free' : `€${totalCost.toFixed(2)}`}</strong>.
                </p>
                <div className="alert alert--warning">
                  This action cannot be undone.
                </div>
                <div className="dialog__actions dialog__actions--with-divider ">
                  <button
                    type="button"
                    className="btn btn--ghost"
                    onClick={() => setConfirming(false)}
                    disabled={submitting}
                  >
                    Back
                  </button>
                  <button
                    type="button"
                    className="btn btn--primary"
                    onClick={handleConfirmBook}
                    disabled={submitting}
                  >
                    {submitting ? 'Booking…' : 'Confirm'}
                  </button>
                </div>
              </section>
            )}

            {/* Guest note */}
            {!isAuthenticated && (
              <section className="section">
                <div className="alert alert--warning">
                  Sign up and get approved to book tickets and contact organizers.
                </div>
              </section>
            )}
          </div>
        </div>
      )}
    </>
  )
}