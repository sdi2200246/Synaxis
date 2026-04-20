import api from './client'
import type { Event, Venue, TicketType , User } from '../types'

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

type BareBooking = {
  id: string
  user_id: string
  ticket_type_id: string
  number_of_tickets: number
  total_cost: number
  status: string
  booked_at: string
}

type BareEvent = Omit<Event, 'venue' | 'categories' | 'category_ids'>

async function hydrateEventLite(bare: BareEvent): Promise<Event> {
  const venue = await api.get<Venue>(`/venues/${bare.venue_id}`).then(r => r.data)
  return { ...bare, venue }
}

export async function getUserBookings(userID: string): Promise<UserBooking[]> {
  const bookings = await api.get<BareBooking[]>(`/users/${userID}/bookings`).then(r => r.data)

  const ticketIds = [...new Set(bookings.map(b => b.ticket_type_id))]
  const tickets = await Promise.all(
    ticketIds.map(id => api.get<TicketType>(`/tickets/${id}`).then(r => r.data))
  )
  const ticketsById = Object.fromEntries(tickets.map(t => [t.id, t]))

  const eventIds = [...new Set(tickets.map(t => t.event_id))]
  const events = await Promise.all(
    eventIds.map(id =>
      api.get<BareEvent>(`/events/${id}`).then(r => r.data).then(hydrateEventLite)
    )
  )
  const eventsById = Object.fromEntries(events.map(e => [e.id, e]))

  return bookings.map(b => {
    const ticket = ticketsById[b.ticket_type_id]
    const event = eventsById[ticket.event_id]
    return {
      id: b.id,
      ticket_type_id: b.ticket_type_id,
      ticket_name: ticket.name,
      number_of_tickets: b.number_of_tickets,
      total_cost: b.total_cost,
      status: b.status,
      booked_at: b.booked_at,
      event_id: event.id,
      event_title: event.title,
      event_start: event.start_datetime,
      venue_name: event.venue.name,
      venue_city: event.venue.city ?? '',
      venue_latitude: event.venue.latitude,
      venue_longitude: event.venue.longitude,
    }
  })
}

export async function getEventBookings(eventId: string): Promise<EventBooking[]> {
  const bookings = await api.get<BareBooking[]>(`/events/${eventId}/bookings`).then(r => r.data)

  const ticketIds = [...new Set(bookings.map(b => b.ticket_type_id))]
  const userIds = [...new Set(bookings.map(b => b.user_id))]

  const [tickets, users] = await Promise.all([
    Promise.all(ticketIds.map(id => api.get<TicketType>(`/tickets/${id}`).then(r => r.data))),
    Promise.all(userIds.map(id => api.get<User>(`/users/${id}`).then(r => r.data))),
  ])

  const ticketsById = Object.fromEntries(tickets.map(t => [t.id, t]))
  const usersById = Object.fromEntries(users.map(u => [u.id, u]))

  return bookings.map(b => {
    const ticket = ticketsById[b.ticket_type_id]
    const user = usersById[b.user_id]
    return {
      id: b.id,
      ticket_name: ticket.name,
      number_of_tickets: b.number_of_tickets,
      total_cost: b.total_cost,
      booked_at: b.booked_at,
      attendee_name: `${user.first_name} ${user.last_name}`,
      attendee_email: user.email,
      attendee_phone: user.phone,
    }
  })
}

export async function createBooking(eventId: string, ticketTypeId: string, quantity: number): Promise<void> {
  await api.post(`/events/${eventId}/bookings`, {
    ticket_type_id: ticketTypeId,
    quantity,
  })
}