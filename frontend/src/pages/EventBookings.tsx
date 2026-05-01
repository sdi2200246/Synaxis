import { useState, useEffect, useMemo } from 'react'
import { useParams, useLocation, Link } from 'react-router-dom'
import { getEventBookings } from '../api/bookings'
import { EventBookingCard } from '../components/bookings/EventBookings'
import type { EventBooking } from '../api/bookings'

export function EventBookingsPage() {
  const { id } = useParams<{ id: string }>()
  const location = useLocation()
  const eventTitle = (location.state as any)?.title || 'Event'
  const eventCapacity = (location.state as any)?.capacity || null
  const venueName = (location.state as any)?.venue || ''

  const [bookings, setBookings] = useState<EventBooking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) return
    getEventBookings(id)
      .then(setBookings)
      .catch(() => setError('Failed to load bookings'))
      .finally(() => setLoading(false))
  }, [id])

  const stats = useMemo(() => {
    const totalTickets = bookings.reduce((sum, b) => sum + b.number_of_tickets, 0)
    const totalRevenue = bookings.reduce((sum, b) => sum + b.total_cost, 0)
    const seatsLeft = eventCapacity != null ? eventCapacity - totalTickets : null

    return { totalTickets, totalRevenue, seatsLeft }
  }, [bookings, eventCapacity])

  if (loading) return <div className="page"><p>Loading bookings…</p></div>
  if (error) return <div className="page"><div className="alert alert--error">{error}</div></div>

  return (
    <div className="page">
      <div className="page-header">
        <Link to="/my-events" className="page-header__back">&larr; Back to My Events</Link>
        <h1>{eventTitle}</h1>
        {venueName && <span className="page-header__subtitle">{venueName}</span>}
      </div>

      <div className="stat-row">
        <div className="stat">
          <span className="stat__value">{bookings.length}</span>
          <span className="label">Bookings</span>
        </div>
        <div className="stat">
          <span className="stat__value">{stats.totalTickets}</span>
          <span className="label">Tickets sold</span>
        </div>
        <div className="stat">
          <span className="stat__value">€{stats.totalRevenue.toFixed(0)}</span>
          <span className="label">Revenue</span>
        </div>
        {stats.seatsLeft != null && (
          <div className="stat">
            <span className="stat__value">{stats.seatsLeft.toLocaleString()}</span>
            <span className="label">Seats left</span>
          </div>
        )}
      </div>

      {bookings.length === 0 ? (
        <p className="empty-state">No bookings yet for this event.</p>
      ) : (
        <>
          <h2 className="eb-section-title">Attendees</h2>
          <div className="list-stack">
            {bookings.map(b => (
              <EventBookingCard key={b.id} booking={b} />
            ))}
          </div>
        </>
      )}
    </div>
  )
}