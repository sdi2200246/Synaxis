import api from './client'

export interface UserBooking {
  id: string
  ticket_type_id: string
  ticket_name: string
  number_of_tickets: number
  total_cost: number
  status: string
  booked_at: string
  event_id: string
  event_title: string
  event_start: string
  venue_name: string
  venue_city: string
  venue_latitude?: number
  venue_longitude?: number
}

export interface EventBooking {
  id: string
  ticket_name: string
  number_of_tickets: number
  total_cost: number
  booked_at: string
  attendee_name: string
  attendee_email: string
  attendee_phone?: string
}

export async function getUserBookings(): Promise<UserBooking[]> {
  const res = await api.get<UserBooking[]>('/my-bookings')
  return res.data
}


export async function getEventBookings(eventId: string): Promise<EventBooking[]> {
  const res = await api.get<EventBooking[]>(`/my-events/${eventId}/bookings`)
  return res.data
}

export async function createBooking(eventId: string, ticketTypeId: string, quantity: number): Promise<void> {
  await api.post(`/events/${eventId}/book`, {
    ticket_type_id: ticketTypeId,
    quantity,
  })
}