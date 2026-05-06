import { createContext, useContext, useState, useEffect, useCallback } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from './AuthContext'
import { getConversations, markConversationAsRead as apiMarkAsRead } from '../api/messages'
import type { ConversationWithParticipants } from '../types'

interface MessagesContextType {
  conversations: ConversationWithParticipants[]
  hasUnread: boolean
  refresh: () => Promise<void>
  markAsRead: (conversationId: string) => Promise<void>
}

const MessagesContext = createContext<MessagesContextType | null>(null)
const POLL_INTERVAL_MS = 15000

export function MessagesProvider({ children }: { children: ReactNode }) {
  const { userRole, isAuthenticated } = useAuth()
  const [conversations, setConversations] = useState<ConversationWithParticipants[]>([])

  const shouldFetch = isAuthenticated && userRole !== 'admin'

  const refresh = useCallback(async () => {
    if (!shouldFetch) return
    try {
      const convs = await getConversations()
      setConversations(convs)
    } catch {
      /* silent — non-critical */
    }
  }, [shouldFetch])

  const markAsRead = useCallback(async (conversationId: string) => {
    await apiMarkAsRead(conversationId)
    setConversations(prev =>
      prev.map(c =>
        c.conversation.id === conversationId
          ? { ...c, conversation: { ...c.conversation, unseen_count: 0 } }
          : c
      )
    )
  }, [])

  useEffect(() => {
    if (!shouldFetch) {
      setConversations([])
      return
    }
    refresh()
    const interval = setInterval(refresh, POLL_INTERVAL_MS)
    return () => clearInterval(interval)
  }, [shouldFetch, refresh])

  const hasUnread = conversations.some(c => c.conversation.unseen_count > 0)

  return (
    <MessagesContext.Provider value={{ conversations, hasUnread, refresh, markAsRead }}>
      {children}
    </MessagesContext.Provider>
  )
}

export function useMessages() {
  const ctx = useContext(MessagesContext)
  if (!ctx) throw new Error('useMessages must be used within MessagesProvider')
  return ctx
}