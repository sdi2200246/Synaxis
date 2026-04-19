import { useState } from 'react'
import type { Event } from '../../types'
import './Events.css'

interface EventCardProps {
  event: Event
  onEdit?: (event: Event) => void
  onPublish?: (event: Event) => void
  onCancel?: (event: Event) => void
  onDelete?: (event: Event) => void
  onTickets?: (event: Event) => void
  onBookings?: (event: Event) => void
}

export function OrganizerEventCard({ event, onEdit, onPublish, onCancel, onDelete , onTickets , onBookings }: EventCardProps) {
  const [expanded, setExpanded] = useState(false)

  const isDraft = event.status === 'DRAFT'
  const isPublished = event.status === 'PUBLISHED'
  const isCancelled = event.status === 'CANCELLED'

  const canPublish = isDraft
  const canCancel = isPublished
  const canDelete = !isCancelled
  const canEdit = !isCancelled
  const canRelease = !isCancelled

  const formatDate = (dt: string) =>
    new Date(dt).toLocaleString('en-US', { dateStyle: 'medium', timeStyle: 'short' })

  const statusClass = isDraft
    ? 'badge-draft'
    : isPublished
    ? 'badge-published'
    : isCancelled
    ? 'badge-cancelled'
    : 'badge-completed'

  return (
    <div className={`ec-card ${expanded ? 'ec-card--expanded' : ''}`}>
      <div className="ec-header" onClick={() => setExpanded(!expanded)}>
        <svg
          className={`ec-chevron ${expanded ? 'ec-chevron--open' : ''}`}
          viewBox="0 0 16 16"
          fill="none"
        >
          <path
            d="M4 6l4 4 4-4"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>

        <div className="ec-main">
          <div className="ec-top">
            <span className="ec-title">{event.title}</span>
            <span className={`ec-badge ${statusClass}`}>{event.status}</span>
          </div>
          <div className="ec-meta">
            <span>{event.venue?.name}, {event.venue?.city}</span>
            <span className="ec-sep">·</span>
            <span>{new Date(event.start_datetime).toLocaleDateString('en-US', { dateStyle: 'medium' })}</span>
            <span className="ec-sep">·</span>
            {!isDraft &&(<button
                className="ec-btn-bookings"
                onClick={(e) => {
                  e.stopPropagation()
                  onBookings?.(event)
                }}
              >
                View bookings
              </button>)}
          </div>
        </div>

        <div className="ec-actions" onClick={e => e.stopPropagation()}>

          {canRelease &&(
            <button className="ec-btn ec-btn--tickets" title="Tickets" onClick={() => onTickets?.(event)}>
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                <rect x="1" y="3" width="13" height="9" rx="1.5" stroke="currentColor" strokeWidth="1.2" />
                <path d="M5 3v9M10 3v9" stroke="currentColor" strokeWidth="1.2" strokeDasharray="1.5 1.5" />
              </svg>
            </button>
          )}

          {canEdit && (
            <button className="ec-btn" title="Edit" onClick={() => onEdit?.(event)}>
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                <path d="M10.5 1.5l3 3-9 9H1.5v-3l9-9z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
              </svg>
            </button>
          )}
          {canPublish && (
            <button className="ec-btn ec-btn--success" title="Publish" onClick={() => onPublish?.(event)}>
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                <path d="M2.5 8l4 4 6-7" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </button>
          )}

          {canCancel && (
            <button className="ec-btn ec-btn--danger" title="Cancel" onClick={() => onCancel?.(event)}>
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                <circle cx="7.5" cy="7.5" r="6" stroke="currentColor" strokeWidth="1.2" />
                <path d="M4.5 7.5h6" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
              </svg>
            </button>
          )}
          {canDelete && (
            <button className="ec-btn ec-btn--danger" title="Delete" onClick={() => onDelete?.(event)}>
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                <path d="M2 4h11M5 4V2.5h5V4M6 7v4M9 7v4M3 4l.8 8.5h7.4L12 4" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </button>
          )}
        </div>
      </div>

      {expanded && (
        <div className="ec-body">
          <div className="ec-grid">
            <div className="ec-detail">
              <span className="ec-label">Type</span>
              <span className="ec-value">{event.event_type}</span>
            </div>
            <div className="ec-detail">
              <span className="ec-label">Capacity</span>
              <span className="ec-value">{event.capacity.toLocaleString()}</span>
            </div>
            <div className="ec-detail">
              <span className="ec-label">Start</span>
              <span className="ec-value">{formatDate(event.start_datetime)}</span>
            </div>
            <div className="ec-detail">
              <span className="ec-label">End</span>
              <span className="ec-value">{formatDate(event.end_datetime)}</span>
            </div>
            <div className="ec-detail ec-detail--full">
              <span className="ec-label">Description</span>
              <span className="ec-value">{event.description}</span>
            </div>
            <div className="ec-detail ec-detail--full">
              <span className="ec-label">Venue</span>
              <div className="ec-venue">
                <span className="ec-venue-name">{event.venue?.name}</span>
                <span className="ec-venue-addr">
                  {event.venue?.address}, {event.venue?.city}, {event.venue?.country}
                  {event.venue?.latitude && ` · ${event.venue.latitude}, ${event.venue.longitude}`}
                </span>
              </div>
            </div>
            {event.categories && event.categories.length > 0 && (
              <div className="ec-detail ec-detail--full">
                <span className="ec-label">Categories</span>
                <div className="ec-cats">
                  {event.categories.map(cat => (
                    <span key={cat.id} className="ec-cat">{cat.name}</span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}