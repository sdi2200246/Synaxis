// in pages/Myevents.tsx
import { useState, useEffect } from 'react'
import type { Event } from '../types'
import { getOrganizerEvents, deleteEvent } from '../api/events'
import { OrganizerEventCard } from '../components/events/OganizerCard'
import { CreateEventForm } from '../components/forms/NewEventForm'
import { EditEventForm } from '../components/forms/EditEventForm'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function MyEventsPage() {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editTarget, setEditTarget] = useState<Event | null>(null)
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [successMessage, setSuccessMessage] = useState('')
  const navigate = useNavigate()
  const {userId} = useAuth()

  const [deleteTarget, setDeleteTarget] = useState<Event | null>(null)
  const [deleteSubmitting, setDeleteSubmitting] = useState(false)

  async function fetchEvents() {
    try {
      if (!userId) return
      const data = await getOrganizerEvents(userId)
      setEvents(data)
    } catch {
      setError('Failed to load events')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchEvents() }, [])

  async function handleConfirmDelete() {
    if (!deleteTarget) return
    setDeleteSubmitting(true)
    try {
      await deleteEvent(deleteTarget.id)
      setDeleteTarget(null)
      setSuccessMessage('Event deleted successfully')
      fetchEvents()
      setTimeout(() => setSuccessMessage(''), 3000)
    } catch (err: any) {
      const msg = err.response?.data?.error || 'Failed to delete event'
      setError(msg)
      setDeleteTarget(null)
      setTimeout(() => setError(''), 3000)
    } finally {
      setDeleteSubmitting(false)
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <h1>My Events</h1>
        <button className="btn-submit" onClick={() => setShowCreateForm(true)}>
          New Event
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="toast">{successMessage}</div>}
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
            onDelete={e => setDeleteTarget(e)}
            onBookings={e => navigate(`/my-events/${e.id}/bookings`, {
              state: { title: e.title, capacity: e.capacity, venue: e.venue?.name }
            })}
          />
        ))}
      </div>

      {deleteTarget && (
        <div className="browse-detail-overlay" onClick={() => !deleteSubmitting && setDeleteTarget(null)}>
          <div className="browse-detail" onClick={e => e.stopPropagation()} style={{ maxWidth: '400px' }}>
            <div className="browse-detail__content">
              <h2 className="browse-detail__title">Delete Event</h2>
              <div className="browse-detail__confirm">
                <p className="browse-detail__confirm-text">
                  Delete <strong>{deleteTarget.title}</strong>?
                </p>
                <p className="browse-detail__confirm-warning">
                  This action cannot be undone.
                </p>
                <div className="browse-detail__confirm-actions">
                  <button
                    className="browse-detail__btn"
                    onClick={() => setDeleteTarget(null)}
                    disabled={deleteSubmitting}
                  >
                    Cancel
                  </button>
                  <button
                    className="browse-detail__confirm-btn"
                    onClick={handleConfirmDelete}
                    disabled={deleteSubmitting}
                    style={{ background: '#ef4444' }}
                  >
                    {deleteSubmitting ? 'Deleting…' : 'Delete'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
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