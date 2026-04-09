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
  await api.delete(`/events/${id}`)
}
