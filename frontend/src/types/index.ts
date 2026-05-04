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
export type UserSummary = {
  id: string
  username: string
  first_name: string
  last_name: string
  email: string
  phone: string
  address? :string
  city?: string
  country?: string
  tax_id?:string
  status?:string
  created_at?: string
  updated_at?:string
}
export interface Category {
  id: string
  name: string
  parent_id?: string | null
}

export interface Venue {
  id: string
  name: string
  address?: string
  city?: string
  country?: string
  latitude?: number
  longitude?: number
  capacity?:number
}

export interface EventMedia {
  id: string
  url: string
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
  venue: Venue
  booking_count?:number
  categories?:Category[]
  category_ids?:string[]
  media?:EventMedia[]
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
  jwt_token: string
}

export interface Category {
  id: string
  name: string
  parent_id?: string | null
}

export interface ConversationParticipant {
  role: string
  user_id: string
}

export interface ConversationData {
  id: string
  booking_id: string
  created_at: string
  unseen_count: number
}

export interface ConversationWithParticipants {
  conversation: ConversationData
  participants: ConversationParticipant[]
  event_title: string
}

export interface Message {
  id: string
  conversation_id: string
  sender_id: string
  content: string
  is_read: boolean
  is_deleted: boolean
  sent_at: string
  updated_at?: string
}