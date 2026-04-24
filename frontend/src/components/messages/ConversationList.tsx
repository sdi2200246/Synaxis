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
    <div className="msg-sidebar">
      <div className="msg-sidebar-header">
        <h2>Messages</h2>
      </div>
 
      <div className="msg-sidebar-list">
        {conversations.length === 0 && (
          <p className="msg-sidebar-empty">No conversations yet</p>
        )}
 
        {conversations.map(conv => {
          const other = conv.participants.find(p => p.user_id !== currentUserId)
          const otherName = other ? (userNames[other.user_id] ?? '...') : 'Unknown'
          const isActive = conv.conversation.id === activeId
          const isUnread = conv.conversation.unseen_count > 0
 
          let className = 'msg-conv'
          if (isActive) className += ' msg-conv--active'
          if (isUnread && !isActive) className += ' msg-conv--unread'
 
          return (
            <div
              key={conv.conversation.id}
              className={className}
              onClick={() => onSelect(conv.conversation.id)}
            >
              <div className="msg-avatar">{getInitials(otherName)}</div>
 
              <div className="msg-conv-body">
                <div className="msg-conv-top">
                  <span className="msg-conv-event">{conv.event_title}</span>
                  {isUnread && !isActive && (
                    <span className="msg-badge">{conv.conversation.unseen_count}</span>
                  )}
                </div>
                <span className="msg-conv-name">{otherName}</span>
                <span className="msg-conv-role">{other?.role}</span>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
 