import api from './client'
import type { Event } from '../types'

export interface SearchEventsParams {
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


export async function searchEvents(params: SearchEventsParams): Promise<SearchEventsResponse> {
  const query = new URLSearchParams()
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

  const res = await api.get<SearchEventsResponse>(`/events?${query.toString()}`)
  return res.data
}


export async function getEvents(): Promise<Event[]> {
  const response = await api.get<Event[]>('/my-events')
  return response.data
}

export async function getEvent(id: string): Promise<Event> {
  const response = await api.get<Event>(`/events/${id}`)
  return response.data
}

export async function createEvent(event: Partial<Event>): Promise<Event> {
  if (event.start_datetime) {
      event.start_datetime = new Date(event.start_datetime).toISOString();
    }
  
  if (event.end_datetime) {
    event.end_datetime = new Date(event.end_datetime).toISOString();
  } 

  if (event.capacity) {
    event.capacity = Number(event.capacity);
  }
  console.log("sending event" , event)
  const response = await api.post<Event>('/events', event)
  return response.data
}

export async function updateEvent(id: string, event: Partial<Event>): Promise<Event> {
  const response = await api.patch<Event>(`/events/${id}`, event)
  return response.data
}

export async function deleteEvent(id: string): Promise<void> {
  await api.delete(`/my-events/${id}`)
}
