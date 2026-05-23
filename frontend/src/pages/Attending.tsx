import { useState, useEffect, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { getUserBookings } from '../api/bookings'
import { getConversations } from '../api/messages'
import { UserBookingCard } from '../components/bookings/UserBookingCard'
import type { UserBooking } from '../api/bookings'
import { FiCalendar, FiTag, FiMapPin } from 'react-icons/fi'
import { useAuth } from '../context/AuthContext'

const BOOKING_TABS = [
  { key: 'all',      label: 'All' },
  { key: 'upcoming', label: 'Upcoming' },
  { key: 'past',     label: 'Past' },
  { key: 'cancelled', label: 'Cancelled' },
] as const

type TabKey = typeof BOOKING_TABS[number]['key']

export function AttendingPage() {
  const [bookings, setBookings] = useState<UserBooking[]>([])
  const [convByBooking, setConvByBooking] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [activeTab, setActiveTab] = useState<TabKey>('all')
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

  const grouped = useMemo(() => {
    const now = Date.now()
    const active = bookings.filter(b => b.status !== 'CANCELLED')
    return {
      all: bookings,
      upcoming:  active.filter(b => new Date(b.event_start).getTime() >  now),
      past:      active.filter(b => new Date(b.event_start).getTime() <= now),
      cancelled: bookings.filter(b => b.status === 'CANCELLED'),
    }
  }, [bookings])

  const counts: Record<TabKey, number> = {
    all:       grouped.all.length,
    upcoming:  grouped.upcoming.length,
    past:      grouped.past.length,
    cancelled: grouped.cancelled.length,
  }

  const stats = useMemo(() => {
    if (bookings.length === 0) return null
    const totalTickets = bookings.reduce((sum, b) => sum + b.number_of_tickets, 0)
    const totalSpent   = bookings.reduce((sum, b) => sum + b.total_cost, 0)
    const nextEvent = grouped.upcoming.length > 0
      ? [...grouped.upcoming].sort(
          (a, b) => new Date(a.event_start).getTime() - new Date(b.event_start).getTime()
        )[0]
      : null
    return { totalTickets, totalSpent, nextEvent }
  }, [bookings, grouped.upcoming])

  const visibleBookings = grouped[activeTab]

  return (
    <div className="page">
      <div className="page-header">
        <div>
          <h1 className="page-header__title">
            My Bookings
            {bookings.length > 0 && (
              <span className="page-header__count">{bookings.length}</span>
            )}
          </h1>
          <p className="page-header__subtitle">Events you're attending</p>
        </div>
      </div>

      {loading ? (
        <p className="empty-state">Loading bookings…</p>
      ) : error ? (
        <div className="alert alert--error">{error}</div>
      ) : bookings.length === 0 ? (
        <div className="empty-state">
          <p>You haven't booked any events yet.</p>
          <Link
            to="/browse"
            className="btn btn--primary"
            style={{ marginTop: 'var(--syn-space-4)' }}
          >
            Browse events
          </Link>
        </div>
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
              <span className="stat__value">{counts.upcoming}</span>
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

          <div className="tabs" role="tablist">
            {BOOKING_TABS.map(tab => (
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

          {visibleBookings.length === 0 ? (
            <div className="empty-state">
              <p>No {activeTab} bookings.</p>
            </div>
          ) : (
            <div className="ub-list">
              {visibleBookings.map(b => (
                <UserBookingCard
                  key={b.id}
                  booking={b}
                  conversationId={convByBooking[b.id] ?? null}
                  currentUserId={userId!}
                  onConversationCreated={handleConversationCreated}
                />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
} 