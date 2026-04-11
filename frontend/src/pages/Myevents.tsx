import { useState, useEffect} from 'react'
import type { Event } from '../types'
import { getEvents } from '../api/events'
import { OrganizerEventCard } from '../components/events/OganizerCard'
import { CreateEventForm } from '../components/forms/NewEventForm'
import { EditEventForm } from '../components/forms/EditEventForm'
import { useNavigate } from 'react-router-dom'

export function MyEventsPage() {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editTarget, setEditTarget] = useState<Event|null>(null)
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [successMessage , setSuccessMessage] = useState('')
  const navigate = useNavigate()

  async function fetchEvents() {
    try {
      const data = await getEvents()
      setEvents(data)
    } catch {
      setError('Failed to load events')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchEvents() }, [])

  return (
    <div className="page">
      <div className="page-header">
        <h1>My Events</h1>
        <button className="btn-submit" onClick={() => setShowCreateForm(true)}>
          New Event
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}
      {loading && <p>Loading...</p>}

      <div className="events-list">
        {events.map(event => (
          <OrganizerEventCard
            key={event.id}
            event={event}
            onEdit={e => setEditTarget(e)}
            onTickets={e => navigate(`/events/${e.id}/tickets`, { state: { title: e.title, capacity: e.capacity } })}
            onPublish={e => console.log('publish', e.id)}
            onCancel={e => console.log('cancel', e.id)}
            onDelete={e => console.log('delete', e.id)}
          />
        ))}
      </div>

      {successMessage && (
        <div className="toast">{successMessage}</div>
      )}

      {showCreateForm && (
        <CreateEventForm
          onClose={() => setShowCreateForm(false)}
           onSuccess={() => {
              setShowCreateForm(false)
              setSuccessMessage('Event created successfully')
              fetchEvents()
              setTimeout(() => setSuccessMessage(''), 3000)
            }}
        />
      )}

      {editTarget && (
        <EditEventForm
          event={editTarget}
          onClose={() => setEditTarget(null)}
          onSuccess={() => {
              setEditTarget(null)
              setSuccessMessage('Event updated successfully')
              fetchEvents()
              setTimeout(() => setSuccessMessage(''), 3000)
            }}
        />
      )}


    </div>
  )
}