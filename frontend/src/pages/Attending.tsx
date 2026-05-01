import { useState, useEffect, useMemo } from 'react'
import { getUserBookings } from '../api/bookings'
import { getConversations } from '../api/messages'
import { UserBookingCard } from '../components/bookings/UserBookingCard'
import type { UserBooking } from '../api/bookings'
import { FiCalendar, FiTag, FiMapPin } from 'react-icons/fi'
import { useAuth } from '../context/AuthContext'
 
export function AttendingPage() {
  const [bookings, setBookings] = useState<UserBooking[]>([])
  const [convByBooking, setConvByBooking] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const { userId } = useAuth()
 
  useEffect(() => {
    if (!userId) return
 
    Promise.all([
      getUserBookings(userId),
      getConversations(),
    ])
      .then(([b, convs]) => {
        setBookings(b)
        const map: Record<string, string> = {}
        for (const c of convs) {
          map[c.conversation.booking_id] = c.conversation.id
        }
        setConvByBooking(map)
      })
      .catch(() => setError('Failed to load bookings'))
      .finally(() => setLoading(false))
  }, [userId])
 
  function handleConversationCreated(bookingId: string, conversationId: string) {
    setConvByBooking(prev => ({ ...prev, [bookingId]: conversationId }))
  }
 
  const stats = useMemo(() => {
    if (bookings.length === 0) return null
 
    const totalTickets = bookings.reduce((sum, b) => sum + b.number_of_tickets, 0)
    const totalSpent = bookings.reduce((sum, b) => sum + b.total_cost, 0)
 
    const upcoming = bookings
      .filter(b => new Date(b.event_start) > new Date())
      .sort((a, b) => new Date(a.event_start).getTime() - new Date(b.event_start).getTime())
 
    const nextEvent = upcoming.length > 0 ? upcoming[0] : null
 
    return { totalTickets, totalSpent, nextEvent, upcomingCount: upcoming.length }
  }, [bookings])
 
  if (loading) return <div className="page"><p>Loading bookings…</p></div>
  if (error) return <div className="page"><div className="alert alert--error">{error}</div></div>
 
  return (
    <div className="page">
      <h1>My Bookings</h1>
 
      {bookings.length === 0 ? (
        <p className="empty-state">You haven't booked any events yet.</p>
      ) : (
        <>
          <div className="stat-row">
            <div className="stat">
              <span className="stat__value">{bookings.length}</span>
              <span className="label">Bookings</span>
            </div>
            <div className="stat">
              <span className="stat__value">{stats!.totalTickets}</span>
              <span className="label">Tickets</span>
            </div>
            <div className="stat">
              <span className="stat__value">{stats!.upcomingCount}</span>
              <span className="label">Upcoming</span>
            </div>
            <div className="stat">
              <span className="stat__value">€{stats!.totalSpent.toFixed(0)}</span>
              <span className="label">Total spent</span>
            </div>
          </div>
 
          {stats!.nextEvent && (
            <div className="card next-event">
              <span className="next-event__label">Next event</span>
              <div className="ub-next__body">
                <span className="next-event__title">{stats!.nextEvent.event_title}</span>
                <div className="next-event__meta">
                  <span><FiMapPin size={13} />{stats!.nextEvent.venue_name}, {stats!.nextEvent.venue_city}</span>
                  <span>
                    <FiCalendar size={13} />
                    {new Date(stats!.nextEvent.event_start).toLocaleDateString('en-US', { dateStyle: 'medium' })}
                    {' · '}
                    {new Date(stats!.nextEvent.event_start).toLocaleTimeString('en-US', { timeStyle: 'short' })}
                  </span>
                  <span><FiTag size={13} />{stats!.nextEvent.ticket_name} ×{stats!.nextEvent.number_of_tickets}</span>
                </div>
              </div>
            </div>
          )}
 
          <h2 className="">All bookings</h2>
          <div className="ub-list">
            {bookings.map(b => (
              <UserBookingCard
                key={b.id}
                booking={b}
                conversationId={convByBooking[b.id] ?? null}
                currentUserId={userId!}
                onConversationCreated={handleConversationCreated}
              />
            ))}
          </div>
        </>
      )}
    </div>
  )
}
 