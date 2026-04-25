// in pages/Myevents.tsx
import { useState, useEffect } from 'react'
import type { Event } from '../types'
import { getOrganizerEvents, deleteEvent , publishEvent } from '../api/events'
import { OrganizerEventCard } from '../components/events/OganizerCard'
import { CreateEventForm } from '../components/forms/NewEventForm'
import { EditEventForm } from '../components/forms/EditEventForm'
import { ConfirmDialog } from '../components/ConfirmDialogue'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function MyEventsPage() {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editTarget, setEditTarget] = useState<Event | null>(null)
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [successMessage, setSuccessMessage] = useState('')
  const [deleteTarget, setDeleteTarget] = useState<Event | null>(null)
  const [deleteSubmitting, setDeleteSubmitting] = useState(false)
  const [publishTarget, setPublishTarget] = useState<Event | null>(null)
  const [publishSubmitting, setPublishSubmitting] = useState(false)
  const navigate = useNavigate()
  const {userId} = useAuth()

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

  async function handleConfirmPublish() {
      if (!publishTarget) return
      setPublishSubmitting(true)
      try {
        await publishEvent(publishTarget.id)
        setPublishTarget(null)
        setSuccessMessage('Event published successfully')
        fetchEvents()
        setTimeout(() => setSuccessMessage(''), 3000)
      } catch (err: any) {
        const msg = err.response?.data?.error || 'Failed to publish event'
        setError(msg)
        setPublishTarget(null)
        setTimeout(() => setError(''), 3000)
      } finally {
        setPublishSubmitting(false)
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
            onPublish={e => setPublishTarget(e)}
            onCancel={e => console.log('cancel', e.id)}
            onDelete={e => setDeleteTarget(e)}
            onBookings={e => navigate(`/my-events/${e.id}/bookings`, {
              state: { title: e.title, capacity: e.capacity, venue: e.venue?.name }
            })}
          />
        ))}
      </div>

      {deleteTarget && (
          <ConfirmDialog
            title="Delete Event"
            body={`Delete "${deleteTarget.title}"? This action cannot be undone.`}
            confirmLabel={deleteSubmitting ? 'Deleting…' : 'Delete'}
            loading={deleteSubmitting}
            onConfirm={handleConfirmDelete}
            onCancel={() => setDeleteTarget(null)}
            confirmClassName="browse-detail__confirm-btn browse-detail__confirm-btn--danger"
            cancelClassName="browse-detail__btn"
          />
        )}

      {publishTarget && (
          <ConfirmDialog
            title="Publish Event"
            body={`"${publishTarget.title}" will be visible to all users and open for bookings.Events can be cancelled after this action.`}
            confirmLabel={publishSubmitting ? 'Publishing…' : 'Publish'}
            loading={publishSubmitting}
            onConfirm={handleConfirmPublish}
            onCancel={() => setPublishTarget(null)}
            confirmClassName="browse-detail__confirm-btn"
            cancelClassName="browse-detail__btn"
          />
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