// frontend/src/pages/MyEvents.tsx
import { useState } from 'react'
import { FiPlus } from 'react-icons/fi'
import { CreateEventForm } from '../components/NewEventForm'

export function MyEventsPage() {
  const [showCreateForm, setShowCreateForm] = useState(false)

  function handleSuccess() {
    setShowCreateForm(false)
    // TODO: Refresh events list
    alert('Event created successfully!')
  }

  return (
    <div className="my-events-page">
      <div className="page-header">
        <div>
          <h1>My Events</h1>
          <p>Events you're organizing</p>
        </div>
        <button className="create-event-btn" onClick={() => setShowCreateForm(true)}>
          <FiPlus size={20} />
          <span>Create Event</span>
        </button>
      </div>

      <div className="events-list">
        <p className="empty-state">No events yet. Create your first event!</p>
      </div>

      {showCreateForm && (
        <CreateEventForm
          onClose={() => setShowCreateForm(false)}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  )
}