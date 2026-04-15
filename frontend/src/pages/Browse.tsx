import { useState, useEffect, useRef, useCallback } from 'react'
import { useStaticData } from '../context/StaticData'
import { searchEvents } from '../api/events'
import { BrowseEventCard } from '../components/events/Browse'
import type { Event } from '../types'
import type { TicketType } from '../api/tickets'

interface CategoryRow {
  id: string
  name: string
  events: Event[]
  loading: boolean
  hasMore: boolean
  offset: number
}

const ROW_LIMIT = 10

export function BrowsePage() {
  const { categories } = useStaticData()
  const [rows, setRows] = useState<CategoryRow[]>([])
  const [initialLoading, setInitialLoading] = useState(true)

  // Initialize rows and fetch first batch per category
  useEffect(() => {
    if (!categories.length) return

    const initial: CategoryRow[] = categories.map(c => ({
      id: c.id,
      name: c.name,
      events: [],
      loading: true,
      hasMore: true,
      offset: 0,
    }))
    setRows(initial)

    Promise.all(
      categories.map(c =>
        searchEvents({ category_id: [c.id], limit: ROW_LIMIT, offset: 0 })
      )
    ).then(results => {
      setRows(prev =>
        prev.map((row, i) => ({
          ...row,
          events: results[i].events,
          hasMore: results[i].has_more,
          loading: false,
          offset: ROW_LIMIT,
        }))
      )
      setInitialLoading(false)
    })
  }, [categories])

  const loadMore = useCallback(async (categoryId: string) => {
    setRows(prev =>
      prev.map(r => r.id === categoryId ? { ...r, loading: true } : r)
    )

    const row = rows.find(r => r.id === categoryId)
    if (!row || !row.hasMore || row.loading) return

    const res = await searchEvents({
      category_id: [categoryId],
      limit: ROW_LIMIT,
      offset: row.offset,
    })

    setRows(prev =>
      prev.map(r =>
        r.id === categoryId
          ? {
              ...r,
              events: [...r.events, ...res.events],
              hasMore: res.has_more,
              loading: false,
              offset: r.offset + ROW_LIMIT,
            }
          : r
      )
    )
  }, [rows])

  function handleBook(event: Event, ticket: TicketType) {
    console.log('book', event.id, ticket.id)
  }

  if (initialLoading) return <div className="page"><p>Loading events…</p></div>

  const nonEmptyRows = rows.filter(r => r.events.length > 0)

  return (
    <div className="page">
      <h1>Browse Events</h1>

      {nonEmptyRows.length === 0 ? (
        <p className="browse-empty">No events to show yet.</p>
      ) : (
        nonEmptyRows.map(row => (
          <CategoryRowView
            key={row.id}
            row={row}
            onLoadMore={() => loadMore(row.id)}
            onBook={handleBook}
          />
        ))
      )}
    </div>
  )
}

interface CategoryRowProps {
  row: CategoryRow
  onLoadMore: () => void
  onBook: (event: Event, ticket: TicketType) => void
}

function CategoryRowView({ row, onLoadMore, onBook }: CategoryRowProps) {
  const sentinelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!row.hasMore || row.loading) return
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting) onLoadMore()
      },
      { root: null, rootMargin: '0px 200px 0px 0px' }
    )
    if (sentinelRef.current) observer.observe(sentinelRef.current)
    return () => observer.disconnect()
  }, [row.hasMore, row.loading, onLoadMore])

  return (
    <div className="browse-section">
      <h2 className="browse-section__title">{row.name}</h2>
      <div className="browse-section__row">
        {row.events.map(ev => (
          <BrowseEventCard key={ev.id} event={ev} onBook={onBook} />
        ))}
        {row.hasMore && <div ref={sentinelRef} className="browse-section__sentinel" />}
      </div>
    </div>
  )
}