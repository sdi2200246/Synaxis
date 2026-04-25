import api from './client'
import type { Event, Venue, Category } from '../types'


type BareEvent = Omit<Event, 'venue' | 'categories' | 'category_ids'>

export interface SearchEventsParams {
  organizer_id?: string
  status?: 'DRAFT' | 'PUBLISHED' | 'COMPLETED' | 'CANCELLED'
  category_id?: string[]
  title?: string
  description?: string
  city?: string
  country?: string
  start_after?: string
  start_before?: string
  min_price?: number
  max_price?: number
  limit?: number
  offset?: number
}

export interface SearchEventsResponse {
  events: Event[]
  has_more: boolean
}

async function hydrateEvent(bare: BareEvent): Promise<Event> {
  const [venue, categories] = await Promise.all([
    api.get<Venue>(`/venues/${bare.venue_id}`).then(r => r.data),
    api.get<Category[]>(`/events/${bare.id}/categories`).then(r => r.data),
  ])
  return { ...bare, venue, categories }
}

export async function searchEvents(params: SearchEventsParams): Promise<SearchEventsResponse> {
  const query = new URLSearchParams()
  if (params.organizer_id) query.append('organizer_id', params.organizer_id)
  if (params.status) query.append('status', params.status)
  if (params.title) query.append('title', params.title)
  if (params.description) query.append('description', params.description)
  if (params.city) query.append('city', params.city)
  if (params.country) query.append('country', params.country)
  if (params.start_after) query.append('start_after', params.start_after)
  if (params.start_before) query.append('start_before', params.start_before)
  if (params.min_price != null) query.append('min_price', String(params.min_price))
  if (params.max_price != null) query.append('max_price', String(params.max_price))
  if (params.limit) query.append('limit', String(params.limit))
  if (params.offset) query.append('offset', String(params.offset))
  params.category_id?.forEach(id => query.append('category_id', id))

  const res = await api.get<{ events: BareEvent[]; has_more: boolean }>(
    `/events?${query.toString()}`
  )
  const events = await Promise.all(res.data.events.map(hydrateEvent))
  return { events, has_more: res.data.has_more }
}

export async function getOrganizerEvents(organizerID: string): Promise<Event[]> {
  const { events } = await searchEvents({ organizer_id: organizerID })
  return events
}

export async function getEvent(id: string): Promise<Event> {
  const response = await api.get<BareEvent>(`/events/${id}`)
  return hydrateEvent(response.data)
}

export async function createEvent(event: Partial<Event>): Promise<void> {
  if (event.start_datetime) event.start_datetime = new Date(event.start_datetime).toISOString()
  if (event.end_datetime) event.end_datetime = new Date(event.end_datetime).toISOString()
  if (event.capacity) event.capacity = Number(event.capacity)
  await api.post('/events', event)
}

export async function updateEvent(id: string, event: Partial<Event>): Promise<void> {
  await api.patch(`/events/${id}`, event)
}

export async function deleteEvent(id: string): Promise<void> {
  await api.delete(`/events/${id}`)
}

export async function publishEvent(id: string): Promise<void> {
     await api.patch(`/events/${id}`, { status: 'PUBLISHED' })
}

export async function cancelEvent(id: string): Promise<void> {
  await api.patch(`/events/${id}`, { status: 'CANCELLED' })
}