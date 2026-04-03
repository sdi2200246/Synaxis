import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getEvents } from '../api'
import type { Event } from '../types'

export function EventsPage() {
  const [events, setEvents] = useState<Event[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      try {
        const data = await getEvents()
        setEvents(data)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load events')
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [])

  if (isLoading) return <div>Loading events...</div>
  if (error) return <div className="error">{error}</div>

  return (
    <div className="events-page">
      <h1>Events</h1>

      {events.length === 0 ? (
        <p>No events found.</p>
      ) : (
        <div className="events-grid">
          {events.map((event) => (
            <Link to={`/events/${event.id}`} key={event.id} className="event-card">
              <h2>{event.title}</h2>
              <p>{event.description}</p>
              <span className="event-date">
                {new Date(event.start_datetime).toLocaleDateString()}
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
