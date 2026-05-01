import { useState, useEffect, useRef, useCallback } from 'react'
import { FiSearch } from 'react-icons/fi'
import { useStaticData } from '../context/StaticData'
import { searchEvents } from '../api/events'
import { BrowseEventCard } from '../components/events/Browse'
import type { Event } from '../types'

const LIMIT = 20

interface Filters {
  q: string
  categoryIDs: string[]
  city: string
  minPrice: string
  maxPrice: string
  startAfter: string
  startBefore: string
}

const emptyFilters: Filters = {
  q: '',
  categoryIDs: [],
  city: '',
  minPrice: '',
  maxPrice: '',
  startAfter: '',
  startBefore: '',
}

export function SearchPage() {
  const { categories } = useStaticData()
  const [filters, setFilters] = useState<Filters>(emptyFilters)
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(false)
  const [hasMore, setHasMore] = useState(false)
  const [offset, setOffset] = useState(0)
  const [searched, setSearched] = useState(false)
  const sentinelRef = useRef<HTMLDivElement>(null)
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(null)

  const fetchEvents = useCallback(async (currentFilters: Filters, currentOffset: number, append: boolean) => {
    setLoading(true)
    try {
      const res = await searchEvents({
        title: currentFilters.q || undefined,
        description: currentFilters.q || undefined,
        category_id: currentFilters.categoryIDs.length ? currentFilters.categoryIDs : undefined,
        city: currentFilters.city || undefined,
        min_price: currentFilters.minPrice ? Number(currentFilters.minPrice) : undefined,
        max_price: currentFilters.maxPrice ? Number(currentFilters.maxPrice) : undefined,
        start_after: currentFilters.startAfter ? new Date(currentFilters.startAfter).toISOString() : undefined,
        start_before: currentFilters.startBefore ? new Date(currentFilters.startBefore).toISOString() : undefined,
        limit: LIMIT,
        offset: currentOffset,
      })
      setEvents(prev => append ? [...prev, ...res.events] : res.events)
      setHasMore(res.has_more)
      setSearched(true)
    } catch {
      setEvents(prev => append ? prev : [])
      setHasMore(false)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)

    debounceRef.current = setTimeout(() => {
      setOffset(0)
      fetchEvents(filters, 0, false)
    }, 300)

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
  }, [filters, fetchEvents])

  useEffect(() => {
    if (!hasMore || loading) return
    const observer = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting) {
        const newOffset = offset + LIMIT
        setOffset(newOffset)
        fetchEvents(filters, newOffset, true)
      }
    })
    if (sentinelRef.current) observer.observe(sentinelRef.current)
    return () => observer.disconnect()
  }, [hasMore, loading, offset, filters, fetchEvents])

  function toggleCategory(id: string) {
    setFilters(prev => ({
      ...prev,
      categoryIDs: prev.categoryIDs.includes(id)
        ? prev.categoryIDs.filter(c => c !== id)
        : [...prev.categoryIDs, id],
    }))
  }

  function removeFilter(key: keyof Filters) {
    setFilters(prev => ({
      ...prev,
      [key]: key === 'categoryIDs' ? [] : '',
    }))
  }

  function removeCategory(id: string) {
    setFilters(prev => ({
      ...prev,
      categoryIDs: prev.categoryIDs.filter(c => c !== id),
    }))
  }

  function clearAll() {
    setFilters(emptyFilters)
  }

  function getCategoryName(id: string) {
    return categories.find(c => c.id === id)?.name ?? ''
  }

  const pills: { label: string; type: string; onRemove: () => void }[] = []
  filters.categoryIDs.forEach(id => {
    pills.push({ label: getCategoryName(id), type: 'cat', onRemove: () => removeCategory(id) })
  })
  if (filters.q) pills.push({ label: `"${filters.q}"`, type: 'text', onRemove: () => removeFilter('q') })
  if (filters.city) pills.push({ label: filters.city, type: 'loc', onRemove: () => removeFilter('city') })
  if (filters.minPrice || filters.maxPrice) {
    const label = filters.minPrice && filters.maxPrice
      ? `€${filters.minPrice} – €${filters.maxPrice}`
      : filters.minPrice ? `Min €${filters.minPrice}` : `Max €${filters.maxPrice}`
    pills.push({ label, type: 'price', onRemove: () => setFilters(prev => ({ ...prev, minPrice: '', maxPrice: '' })) })
  }
  if (filters.startAfter) pills.push({ label: `From ${filters.startAfter}`, type: 'date', onRemove: () => removeFilter('startAfter') })
  if (filters.startBefore) pills.push({ label: `To ${filters.startBefore}`, type: 'date', onRemove: () => removeFilter('startBefore') })

  const hasFilters = pills.length > 0

  return (
    <div className="page">
      <h1>Search Events</h1>

      <div className="search-bar">
        <FiSearch size={18} className="search-bar__icon" />
        <input
          type="text"
          className="search-bar__input"
          placeholder="Search by title or description..."
          value={filters.q}
          onChange={e => setFilters(prev => ({ ...prev, q: e.target.value }))}
        />
      </div>

      {pills.length > 0 && (
        <div className="search-pills">
          {pills.map((p, i) => (
            <span key={i} className={`search-pill search-pill--${p.type}`}>
              {p.label}
              <span className="search-pill__x" onClick={p.onRemove}>&times;</span>
            </span>
          ))}
          <button className="btn btn--ghost" onClick={clearAll}>Clear all</button>
        </div>
      )}

      <div className="search-layout">
        <aside className="search-sidebar">
          <div className="search-sidebar__section">
            <h3 className="search-sidebar__label">Categories</h3>
            <div className="search-sidebar__cats">
              {categories.map(c => (
                <div
                  key={c.id}
                  className={`search-cat ${filters.categoryIDs.includes(c.id) ? 'search-cat--active' : ''}`}
                  onClick={() => toggleCategory(c.id)}
                >
                  <div className="search-cat__check" />
                  {c.name}
                </div>
              ))}
            </div>
          </div>

          <div className="search-sidebar__section">
            <h3 className="search-sidebar__label">Location</h3>
            <input
              type="text"
              placeholder="City..."
              value={filters.city}
              onChange={e => setFilters(prev => ({ ...prev, city: e.target.value }))}
            />
          </div>

          <div className="search-sidebar__section">
            <h3 className="search-sidebar__label">Price range</h3>
            <div className="search-sidebar__price-row">
              <input
                type="number"
                placeholder="Min"
                min="0"
                value={filters.minPrice}
                onChange={e => setFilters(prev => ({ ...prev, minPrice: e.target.value }))}
              />
              <input
                type="number"
                placeholder="Max"
                min="0"
                value={filters.maxPrice}
                onChange={e => setFilters(prev => ({ ...prev, maxPrice: e.target.value }))}
              />
            </div>
          </div>

          <div className="search-sidebar__section">
            <h3 className="search-sidebar__label">Date range</h3>
            <div className="search-sidebar__date-fields">
              <label>From</label>
              <input
                type="date"
                value={filters.startAfter}
                onChange={e => setFilters(prev => ({ ...prev, startAfter: e.target.value }))}
              />
              <label>To</label>
              <input
                type="date"
                value={filters.startBefore}
                onChange={e => setFilters(prev => ({ ...prev, startBefore: e.target.value }))}
              />
            </div>
          </div>
        </aside>

        <main className="search-results">
          {searched && (
            <p className="search-results__count">
              {events.length} event{events.length !== 1 ? 's' : ''} found
              {hasFilters ? '' : ' (showing all)'}
            </p>
          )}

          {events.length > 0 ? (
            <div className="search-grid">
              {events.map(ev => (
                <BrowseEventCard key={ev.id} event={ev} />
              ))}
            </div>
          ) : searched && !loading ? (
            <p className="search-results__empty">No events match your filters.</p>
          ) : null}

          {loading && <p className="search-results__loading">Loading…</p>}
          {hasMore && !loading && <div ref={sentinelRef} className="search-sentinel" />}
        </main>
      </div>
    </div>
  )
}