import api from './client'
import type { ConversationWithParticipants, Message } from '../types'

interface ConversationsResponse {
  conversations: ConversationWithParticipants[]
}

interface MessagesResponse {
  messages: Message[]
}

interface CreateConversationResponse {
  conversation_id: string
}

export async function getConversations(): Promise<ConversationWithParticipants[]> {
  const res = await api.get<ConversationsResponse>('/conversations')
  return res.data.conversations
}

export async function getConversationMessages(conversationId: string): Promise<Message[]> {
  const res = await api.get<MessagesResponse>(`/conversations/${conversationId}/messages`)
  return res.data.messages
}

export async function sendMessage(conversationId: string, content: string): Promise<void> {
  await api.post(`/conversations/${conversationId}/messages`, { content })
}

export async function markConversationAsRead(conversationId: string): Promise<void> {
  await api.patch(`/conversations/${conversationId}/read`)
}

export async function updateMessage(
  messageId: string,
  payload: { content?: string; delete?: number }
): Promise<void> {
  await api.patch(`/messages/${messageId}`, payload)
}

export async function createConversation(
  bookingId: string,
  organizerId: string,
  attendeeId: string
): Promise<string> {
  const res = await api.post<CreateConversationResponse>('/conversations', {
    booking_id: bookingId,
    organizer_id: organizerId,
    attendee_id: attendeeId,
  })
  return res.data.conversation_id
}