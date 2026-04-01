import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getEvent } from '../api'
import type { Event } from '../types'

export function EventDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [event, setEvent] = useState<Event | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      if (!id) return
      try {
        const data = await getEvent(id)
        setEvent(data)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load event')
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [id])

  if (isLoading) return <div>Loading event...</div>
  if (error) return <div className="error">{error}</div>
  if (!event) return <div>Event not found</div>

  return (
    <div className="event-detail">
      <Link to="/events">&larr; Back to Events</Link>

      <h1>{event.title}</h1>
      <span className={`status status-${event.status.toLowerCase()}`}>
        {event.status}
      </span>

      <p className="description">{event.description}</p>

      <div className="meta">
        <div>
          <strong>Date:</strong>{' '}
          {new Date(event.start_datetime).toLocaleString()} –{' '}
          {new Date(event.end_datetime).toLocaleString()}
        </div>
        <div>
          <strong>Capacity:</strong> {event.capacity}
        </div>
        <div>
          <strong>Type:</strong> {event.event_type}
        </div>
      </div>

      {event.venue && (
        <div className="venue">
          <h2>Venue</h2>
          <p>{event.venue.name}</p>
          <p>{event.venue.address}, {event.venue.city}, {event.venue.country}</p>
        </div>
      )}
    </div>
  )
}
