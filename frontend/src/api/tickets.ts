import api from './client'

export interface TicketType {
  id: string
  event_id: string
  name: string
  price: number
  quantity: number
  available: number
}

export interface CreateTicketTypeInput {
  name: string
  price: number
  quantity: number
}

export async function getTicketTypes(eventId: string): Promise<TicketType[]> {
  const res = await api.get<TicketType[]>(`/events/${eventId}/tickets`)
  return res.data
}

export async function createTicketType(eventId: string, input: CreateTicketTypeInput): Promise<void> {
  await api.post(`/events/${eventId}/tickets`, input)
}

export async function updateTicketType(eventId:string ,ticketId: string, data: Partial<Pick<TicketType, 'name' | 'price' | 'quantity'>>): Promise<void> {
  await api.patch(`/events/${eventId}/tickets/${ticketId}`, data)
}