import type { ConversationWithParticipants } from '../../types'

interface Props {
  conversations: ConversationWithParticipants[]
  activeId: string | null
  currentUserId: string
  userNames: Record<string, string>
  onSelect: (conversationId: string) => void
}

function getInitials(fullName: string): string {
  return fullName
    .split(' ')
    .map(w => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

export function ConversationList({ conversations, activeId, currentUserId, userNames, onSelect }: Props) {
  return (
    <div className="message-page__sidebar">
      <div className="conversation-list__header">
        <h2 className="conversation-list__title">Messages</h2>
      </div>

      <div className="conversation-list">
        {conversations.length === 0 && (
          <p className="empty-state">No conversations yet</p>
        )}

        {conversations.map(conv => {
          const other = conv.participants.find(p => p.user_id !== currentUserId)
          const otherName = other ? (userNames[other.user_id] ?? '...') : 'Unknown'
          const isActive = conv.conversation.id === activeId
          const isUnread = conv.conversation.unseen_count > 0

          let className = 'conversation-item'
          if (isActive) className += ' is-active'
          if (isUnread && !isActive) className += ' is-unread'

          return (
            <div
              key={conv.conversation.id}
              className={className}
              onClick={() => onSelect(conv.conversation.id)}
            >
              <div className="avatar">{getInitials(otherName)}</div>

              <div className="conversation-item__body">
                <span className="conversation-item__event">{conv.event_title}</span>
                <span className="conversation-item__name">{otherName}</span>
              </div>

              {isUnread && !isActive && (
                <span className="conversation-item__unread-count">
                  {conv.conversation.unseen_count}
                </span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}