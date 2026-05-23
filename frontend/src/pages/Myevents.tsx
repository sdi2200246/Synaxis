import { useState, useEffect, useMemo } from 'react'
import type { Event } from '../types'
import { getOrganizerEvents, deleteEvent, publishEvent, cancelEvent } from '../api/events'
import { OrganizerEventCard } from '../components/events/OganizerCard'
import { CreateEventForm } from '../components/forms/NewEventForm'
import { EditEventForm } from '../components/forms/EditEventForm'
import { ConfirmDialog } from '../components/ConfirmDialogue'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const STATUS_TABS = [
  { key: 'all',       label: 'All' },
  { key: 'DRAFT',     label: 'Draft' },
  { key: 'PUBLISHED', label: 'Published' },
  { key: 'CANCELLED', label: 'Cancelled' },
  { key: 'COMPLETED', label: 'Completed' },
] as const

type TabKey = typeof STATUS_TABS[number]['key']

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
  const [cancelTarget, setCancelTarget] = useState<Event | null>(null)
  const [cancelSubmitting, setCancelSubmitting] = useState(false)
  const [activeTab, setActiveTab] = useState<TabKey>('all')

  const navigate = useNavigate()
  const { userId } = useAuth()

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

  const counts = useMemo(() => {
    const c: Record<TabKey, number> = {
      all: events.length,
      DRAFT: 0, PUBLISHED: 0, CANCELLED: 0, COMPLETED: 0,
    }
    for (const e of events) {
      if (e.status in c) c[e.status as TabKey] += 1
    }
    return c
  }, [events])

  const visibleEvents = useMemo(() => {
    if (activeTab === 'all') return events
    return events.filter(e => e.status === activeTab)
  }, [events, activeTab])

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
      setError(err.response?.data?.error || 'Failed to delete event')
      setDeleteTarget(null)
      setTimeout(() => setError(''), 3000)
    } finally {
      setDeleteSubmitting(false)
    }
  }

  async function handleConfirmCancel() {
    if (!cancelTarget) return
    setCancelSubmitting(true)
    try {
      await cancelEvent(cancelTarget.id)
      setCancelTarget(null)
      setSuccessMessage('Event cancelled successfully')
      fetchEvents()
      setTimeout(() => setSuccessMessage(''), 3000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to cancel event')
      setCancelTarget(null)
      setTimeout(() => setError(''), 3000)
    } finally {
      setCancelSubmitting(false)
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
      setError(err.response?.data?.error || 'Failed to publish event')
      setPublishTarget(null)
      setTimeout(() => setError(''), 3000)
    } finally {
      setPublishSubmitting(false)
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <div>
          <h1 className="page-header__title">
            My Events
            {events.length > 0 && (
              <span className="page-header__count">{events.length}</span>
            )}
          </h1>
          <p className="page-header__subtitle">Events you organize</p>
        </div>
        <button className="btn btn--primary" onClick={() => setShowCreateForm(true)}>
          New Event
        </button>
      </div>

      {error && <div className="alert alert--error">{error}</div>}
      {successMessage && <div className="toast toast--success">{successMessage}</div>}
      {loading && <p>Loading...</p>}

     {loading ? (
          <p className="empty-state">Loading…</p>
        ) : events.length === 0 ? (
          <div className="empty-state">
            <p>You haven't created any events yet.</p>
            <button
              className="btn btn--primary"
              style={{ marginTop: 'var(--syn-space-4)' }}
              onClick={() => setShowCreateForm(true)}
            >
              Create your first event
            </button>
          </div>
        ) : (
          <>
            <div className="tabs" role="tablist">
              {STATUS_TABS.map(tab => (
                <button
                  key={tab.key}
                  type="button"
                  role="tab"
                  aria-selected={activeTab === tab.key}
                  className={`tabs__item ${activeTab === tab.key ? 'is-active' : ''}`}
                  onClick={() => setActiveTab(tab.key)}
                >
                  {tab.label}
                  <span className="tabs__count">{counts[tab.key]}</span>
                </button>
              ))}
            </div>

            {visibleEvents.length === 0 ? (
              <div className="empty-state">
                <p>No {activeTab.toLowerCase()} events.</p>
              </div>
            ) : (
              <div className="list-stack">
                {visibleEvents.map(event => (
                  <OrganizerEventCard
                    key={event.id}
                    event={event}
                    onEdit={e => setEditTarget(e)}
                    onTickets={e => navigate(`/events/${e.id}/tickets`, { state: { title: e.title, capacity: e.capacity } })}
                    onPublish={e => setPublishTarget(e)}
                    onCancel={e => setCancelTarget(e)}
                    onDelete={e => setDeleteTarget(e)}
                    onBookings={e => navigate(`/my-events/${e.id}/bookings`, {
                      state: { title: e.title, capacity: e.capacity, venue: e.venue?.name }
                    })}
                  />
                ))}
              </div>
            )}
          </>
        )}
        
      {deleteTarget && (
        <ConfirmDialog
          title="Delete Event"
          body={`Delete "${deleteTarget.title}"? This action cannot be undone.`}
          confirmLabel={deleteSubmitting ? 'Deleting…' : 'Delete'}
          loading={deleteSubmitting}
          onConfirm={handleConfirmDelete}
          onCancel={() => setDeleteTarget(null)}
        />
      )}

      {publishTarget && (
        <ConfirmDialog
          title="Publish Event"
          body={`"${publishTarget.title}" will be visible to all users and open for bookings. Events can be cancelled after this action.`}
          confirmLabel={publishSubmitting ? 'Publishing…' : 'Publish'}
          loading={publishSubmitting}
          onConfirm={handleConfirmPublish}
          onCancel={() => setPublishTarget(null)}
        />
      )}

      {cancelTarget && (
        <ConfirmDialog
          title="Cancel Event"
          body={`"${cancelTarget.title}" will be cancelled and all attendees will be notified via a direct message. This cannot be undone.`}
          confirmLabel={cancelSubmitting ? 'Cancelling…' : 'Cancel Event'}
          loading={cancelSubmitting}
          onConfirm={handleConfirmCancel}
          onCancel={() => setCancelTarget(null)}
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