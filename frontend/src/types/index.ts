// Matches backend entities

export interface User {
  id: string
  username: string
  first_name: string
  last_name: string
  email: string
  phone: string
  address: string
  city: string
  country: string
  tax_id: string
  role: 'admin' | 'user'
  status: 'pending' | 'approved' | 'rejected'
  created_at: string
}

export interface Venue {
  id: string
  name: string
  address: string
  city: string
  country: string
  latitude?: number
  longitude?: number
}

export interface Event {
  id: string
  organizer_id: string
  venue_id: string
  title: string
  event_type: string
  status: 'DRAFT' | 'PUBLISHED' | 'COMPLETED' | 'CANCELLED'
  description: string
  capacity: number
  start_datetime: string
  end_datetime: string
  created_at: string
  venue?: Venue
}

export interface TicketType {
  id: string
  event_id: string
  name: string
  price: number
  quantity: number
  available: number
}

export interface Booking {
  id: string
  user_id: string
  ticket_type_id: string
  number_of_tickets: number
  total_cost: number
  status: 'ACTIVE' | 'COMPLETED' | 'CANCELLED'
  booked_at: string
}

// API payloads
export interface LoginCredentials {
  username: string
  password: string
}

export interface RegisterPayload {
  username: string
  password: string
  first_name: string
  last_name: string
  email: string
  phone: string
  address: string
  city: string
  country: string
  tax_id: string
}

export interface LoginResponse {
  token: string
}
