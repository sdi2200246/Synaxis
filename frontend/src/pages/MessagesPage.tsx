import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { ConversationList } from '../components/messages/ConversationList'
import { ChatBubble } from '../components/messages/ChatBubble'
import {
  getConversations,
  getConversationMessages,
  sendMessage,
  markConversationAsRead,
} from '../api/messages'
import type { ConversationWithParticipants, Message } from '../types'
import '../components/messages/Messages.css'

async function fetchUserName(userId: string): Promise<string> {
  const { default: api } = await import('../api/client')
  const res = await api.get<{ first_name: string; last_name: string }>(`/users/${userId}`)
  return `${res.data.first_name} ${res.data.last_name}`
}

export function MessagesPage() {
  const { conversationId } = useParams<{ conversationId?: string }>()
  const navigate = useNavigate()
  const { userId } = useAuth()

  const [conversations, setConversations] = useState<ConversationWithParticipants[]>([])
  const [messages, setMessages] = useState<Message[]>([])
  const [userNames, setUserNames] = useState<Record<string, string>>({})
  const [draft, setDraft] = useState('')
  const [sending, setSending] = useState(false)

  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Load conversations
  useEffect(() => {
    getConversations().then(setConversations)
  }, [])

  // Resolve participant names
  useEffect(() => {
    if (conversations.length === 0) return

    const ids = new Set<string>()
    for (const c of conversations) {
      for (const p of c.participants) {
        if (p.user_id !== userId && !userNames[p.user_id]) {
          ids.add(p.user_id)
        }
      }
    }

    if (ids.size === 0) return

    Promise.all(
      [...ids].map(async id => {
        const name = await fetchUserName(id)
        return [id, name] as const
      })
    ).then(pairs => {
      setUserNames(prev => {
        const next = { ...prev }
        for (const [id, name] of pairs) next[id] = name
        return next
      })
    })
  }, [conversations, userId])

  // Load messages when active conversation changes
  useEffect(() => {
    if (!conversationId) {
      setMessages([])
      return
    }

    getConversationMessages(conversationId).then(msgs => {
      setMessages(msgs)
      scrollToBottom()
    })

    // Mark as read
    const conv = conversations.find(c => c.conversation.id === conversationId)
    if (conv && conv.conversation.unseen_count > 0) {
      markConversationAsRead(conversationId).then(() => {
        setConversations(prev =>
          prev.map(c =>
            c.conversation.id === conversationId
              ? { ...c, conversation: { ...c.conversation, unseen_count: 0 } }
              : c
          )
        )
      })
    }
  }, [conversationId])

  function scrollToBottom() {
    setTimeout(() => messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' }), 50)
  }

  function handleSelect(id: string) {
    navigate(`/messages/${id}`)
  }

  async function handleUpdate(id: string, payload: { content?: string; delete?: number }) {
    const { updateMessage } = await import('../api/messages')
    await updateMessage(id, payload)
    if (conversationId) {
      const msgs = await getConversationMessages(conversationId)
      setMessages(msgs)
    }
  }

  async function handleSend() {
    if (!conversationId || !draft.trim() || sending) return

    setSending(true)
    try {
      await sendMessage(conversationId, draft.trim())
      setDraft('')
      const msgs = await getConversationMessages(conversationId)
      setMessages(msgs)
      scrollToBottom()
    } finally {
      setSending(false)
    }
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const activeConv = conversations.find(c => c.conversation.id === conversationId)
  const otherParticipant = activeConv?.participants.find(p => p.user_id !== userId)
  const otherName = otherParticipant ? (userNames[otherParticipant.user_id] ?? '...') : ''

  function getInitials(name: string): string {
    return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
  }

  return (
    <div className="card message-page">
      <ConversationList
        conversations={conversations}
        activeId={conversationId ?? null}
        currentUserId={userId!}
        userNames={userNames}
        onSelect={handleSelect}
      />

      <div className="message-page__chat">
        {!conversationId ? (
          <div className="chat-empty">Select a conversation</div>
        ) : (
          <>
            <div className="chat-header">
              <div className="avatar avatar--lg">
                {otherName ? getInitials(otherName) : '??'}
              </div>
              <div>
                <p className="chat-header__title">{otherName}</p>
                <p className="chat-header__sub">{activeConv?.event_title}</p>
              </div>
            </div>

            <div className="message-feed">
              {messages.map(msg => (
                <ChatBubble
                  key={msg.id}
                  message={msg}
                  isOutgoing={msg.sender_id === userId}
                  onUpdate={handleUpdate}
                />
              ))}
              <div ref={messagesEndRef} />
            </div>

            <div className="composer">
              <input
                className="composer__input"
                type="text"
                placeholder="Type a message..."
                value={draft}
                onChange={e => setDraft(e.target.value)}
                onKeyDown={handleKeyDown}
              />
              <button
                className="btn btn--primary"
                onClick={handleSend}
                disabled={!draft.trim() || sending}
              >
                Send
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}