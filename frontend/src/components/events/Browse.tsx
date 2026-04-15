import { useState, useEffect } from 'react'
import { FiMapPin, FiCalendar, FiTag } from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'
import { EventMap } from './Map'
import { getTicketTypes } from '../../api/tickets'
import type { Event } from '../../types'
import type { TicketType } from '../../api/tickets'
import './Events.css'

interface BrowseEventCardProps {
  event: Event
  onBook?: (event: Event, ticketType: TicketType) => void
}

export function BrowseEventCard({ event, onBook }: BrowseEventCardProps) {
  const [open, setOpen] = useState(false)
  const [tickets, setTickets] = useState<TicketType[]>([])
  const [loadingTickets, setLoadingTickets] = useState(false)
  const [showMap, setShowMap] = useState(false)
  const { isAuthenticated } = useAuth()

  const hasImage = false // TODO: wire to event media
  const hasCoords = event.venue.latitude != null && event.venue.longitude != null

  const dateLabel = new Date(event.start_datetime).toLocaleDateString('en-US', {
    dateStyle: 'medium',
  })

  useEffect(() => {
    if (!open) return
    setLoadingTickets(true)
    getTicketTypes(event.id)
      .then(setTickets)
      .catch(() => setTickets([]))
      .finally(() => setLoadingTickets(false))
  }, [open, event.id])

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
        <div className="browse-detail-overlay" onClick={() => setOpen(false)}>
          <div className="browse-detail" onClick={e => e.stopPropagation()}>
            <button className="browse-detail__close" onClick={() => setOpen(false)}>
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

              <div className="browse-detail__tickets">
                <h3 className="browse-detail__section-title">Tickets</h3>
                {loadingTickets ? (
                  <p className="browse-detail__tickets-loading">Loading tickets…</p>
                ) : tickets.length === 0 ? (
                  <p className="browse-detail__tickets-empty">No tickets available.</p>
                ) : (
                  <div className="browse-detail__ticket-list">
                    {tickets.map(t => (
                      <div key={t.id} className="browse-detail__ticket">
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
                              onClick={() => onBook?.(event, t)}
                            >
                              Book
                            </button>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {!isAuthenticated && (
                <div className="browse-detail__guest-note">
                  Sign up and get approved to book tickets and contact organizers.
                </div>
              )}

              {isAuthenticated && (
                <div className="browse-detail__actions">
                  <button className="browse-detail__btn">Message organizer</button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  )
}