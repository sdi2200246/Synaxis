import { useState, useEffect, useCallback} from 'react'
import { useStaticData } from '../context/StaticData'
import { searchEvents } from '../api/events'
import { BrowseEventCard } from '../components/events/Browse'
import type { Event } from '../types'
import { FiCalendar, FiMapPin, FiChevronLeft, FiChevronRight } from 'react-icons/fi'

interface CategoryRow {
  id: string
  name: string
  events: Event[]
  loading: boolean
  hasMore: boolean
  offset: number
}

const ROW_LIMIT = 10
const HERO_LIMIT = 5
const HERO_INTERVAL_MS = 6000

export function BrowsePage() {
  const { categories } = useStaticData()
  const [rows, setRows] = useState<CategoryRow[]>([])
  const [initialLoading, setInitialLoading] = useState(true)
  const [heroEvents, setHeroEvents] = useState<Event[]>([])
  const [heroIndex, setHeroIndex] = useState(0)
  const [scrollY, setScrollY] = useState(0)

  useEffect(() => {
    const main = document.querySelector('.main-content') as HTMLElement | null
    const target: HTMLElement | Window = main ?? window

    function onScroll() {
      const y = main ? main.scrollTop : window.scrollY
      setScrollY(y)
    }

    target.addEventListener('scroll', onScroll, { passive: true })
    return () => target.removeEventListener('scroll', onScroll)
  }, [])

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
    )
      .then(results => {
        setRows(prev =>
          prev.map((row, i) => ({
            ...row,
            events: results[i]?.events ?? [],
            hasMore: results[i]?.has_more ?? false,
            loading: false,
            offset: ROW_LIMIT,
          }))
        )
      })
      .finally(() => {
        setInitialLoading(false)
      })
  }, [categories])

  useEffect(() => {
    searchEvents({
      start_after: new Date().toISOString(),
      limit: HERO_LIMIT,
      offset: 0,
    })
      .then(res => setHeroEvents(res.events))
      .catch(() => setHeroEvents([]))
  }, [])

  useEffect(() => {
    if (heroEvents.length <= 1) return

    const id = setInterval(() => {
      setHeroIndex(i => (i + 1) % heroEvents.length)
    }, HERO_INTERVAL_MS)

    return () => clearInterval(id)
  }, [heroEvents.length])

  const loadMore = useCallback(async (categoryId: string) => {
    const targetRow = rows.find(r => r.id === categoryId)
    if (!targetRow || !targetRow.hasMore || targetRow.loading) return

    setRows(prev =>
      prev.map(row => (row.id === categoryId ? { ...row, loading: true } : row))
    )

    try {
      const res = await searchEvents({
        category_id: [categoryId],
        limit: ROW_LIMIT,
        offset: targetRow.offset,
      })

      setRows(prev =>
        prev.map(row =>
          row.id === categoryId
            ? {
                ...row,
                events: [...row.events, ...res.events],
                hasMore: res.has_more,
                loading: false,
                offset: row.offset + ROW_LIMIT,
              }
            : row
        )
      )
    } catch {
      setRows(prev =>
        prev.map(row =>
          row.id === categoryId ? { ...row, loading: false } : row
        )
      )
    }
  }, [rows]) 

  if (initialLoading) {
    return <div className="page"><p>Loading events…</p></div>
  }

  const nonEmptyRows = rows.filter(r => r.events.length > 0)
  const heroOpacity = Math.max(0, 1 - scrollY / 280)
  const heroEvent = heroEvents[heroIndex]

  return (
    <div className="browse-page">
      {heroEvent && (
        <div
          className="browse-hero"
          style={{ opacity: heroOpacity }}
          aria-hidden={heroOpacity < 0.1}
        >
          <div
              className="browse-hero__bg"
              style={heroEvent.media?.length ? {
                backgroundImage: `url(${heroEvent.media[0].url})`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
              } : undefined}
            />
          <div className="browse-hero__content">
            <span className="browse-hero__eyebrow">Featured · Upcoming</span>
            <h1 className="browse-hero__title">{heroEvent.title}</h1>
            <div className="browse-hero__meta">
              <span>
                <FiMapPin size={14} />
                {heroEvent.venue.name}, {heroEvent.venue.city}
              </span>
              <span>
                <FiCalendar size={14} />
                {new Date(heroEvent.start_datetime).toLocaleDateString('en-US', {
                  dateStyle: 'long',
                })}
              </span>
            </div>
          </div>

          {heroEvents.length > 1 && (
            <div className="browse-hero__dots">
              {heroEvents.map((_, i) => (
                <button
                  key={i}
                  type="button"
                  className={`browse-hero__dot ${i === heroIndex ? 'is-active' : ''}`}
                  onClick={() => setHeroIndex(i)}
                  aria-label={`Show event ${i + 1}`}
                />
              ))}
            </div>
          )}
        </div>
      )}

      <div className="browse-rows">
        <h1 className="browse-rows__title">Browse Events</h1>

        {nonEmptyRows.length === 0 ? (
          <p className="empty-state">No events to show yet.</p>
        ) : (
          nonEmptyRows.map(row => (
            <CategoryRowView
              key={row.id}
              row={row}
              onLoadMore={() => loadMore(row.id)}
            />
          ))
        )}
      </div>
    </div>
  )
}

const VISIBLE_COUNT = 5; 

interface CategoryRowProps {
  row: CategoryRow
  onLoadMore: () => void
}

function CategoryRowView({ row, onLoadMore }: CategoryRowProps) {
  const [currentIndex, setCurrentIndex] = useState(0)

  useEffect(() => {
    if (
      row.hasMore &&
      !row.loading &&
      currentIndex + VISIBLE_COUNT >= row.events.length
    ) {
      onLoadMore()
    }
  }, [currentIndex, row.events.length, row.hasMore, row.loading, onLoadMore])

  const handlePrev = () => {
    setCurrentIndex(prev => Math.max(0, prev - VISIBLE_COUNT))
  }

  const handleNext = () => {
    if (currentIndex + VISIBLE_COUNT < row.events.length || row.hasMore) {
      setCurrentIndex(prev => prev + VISIBLE_COUNT)
    }
  }

  const visibleEvents = row.events.slice(
    currentIndex, 
    currentIndex + VISIBLE_COUNT
  )

  const isAtStart = currentIndex === 0
  const isAtEnd = currentIndex + VISIBLE_COUNT >= row.events.length && !row.hasMore

  return (
    <section className="browse-row">
      <header className="browse-row__head">
        <h2 className="browse-row__title">{row.name}</h2>

        <div className="browse-row__nav">
          <button
            type="button"
            className="browse-row__arrow"
            onClick={handlePrev}
            disabled={isAtStart}
            aria-label="Previous events"
            style={{ opacity: isAtStart ? 0.5 : 1, cursor: isAtStart ? 'not-allowed' : 'pointer' }}
          >
            <FiChevronLeft size={18} />
          </button>

          <button
            type="button"
            className="browse-row__arrow"
            onClick={handleNext}
            disabled={isAtEnd || (currentIndex + VISIBLE_COUNT >= row.events.length && row.loading)}
            aria-label="Next events"
            style={{ opacity: isAtEnd ? 0.5 : 1, cursor: isAtEnd ? 'not-allowed' : 'pointer' }}
          >
            <FiChevronRight size={18} />
          </button>
        </div>
      </header>

      <div className="browse-row__track">
        {visibleEvents.map(ev => (
          <BrowseEventCard key={ev.id} event={ev} />
        ))}

        {row.loading && currentIndex + VISIBLE_COUNT > row.events.length && (
          <div className="browse-row__loading-state" style={{ padding: '1rem' }}>
            Loading...
          </div>
        )}
      </div>
    </section>
  )
}