import api from './client'
import type { Event } from '../types'

export async function getEvents(): Promise<Event[]> {
  const response = await api.get<Event[]>('/events')
  return response.data
}

export async function getEvent(id: string): Promise<Event> {
  const response = await api.get<Event>(`/events/${id}`)
  return response.data
}

export async function createEvent(event: Partial<Event>): Promise<Event> {
  const response = await api.post<Event>('/events', event)
  return response.data
}

export async function updateEvent(id: string, event: Partial<Event>): Promise<Event> {
  const response = await api.put<Event>(`/events/${id}`, event)
  return response.data
}

export async function deleteEvent(id: string): Promise<void> {
  await api.delete(`/events/${id}`)
}
