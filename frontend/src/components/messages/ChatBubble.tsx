import { useState, useRef, useEffect } from 'react'
import type { Message } from '../../types'
 
interface Props {
  message: Message
  isOutgoing: boolean
  onUpdate: (id: string, payload: { content?: string; delete?: number }) => Promise<void>
}
 
function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
  })
}
 
const MENU_WIDTH = 192
 
export function ChatBubble({ message, isOutgoing, onUpdate }: Props) {
  const [menuOpen, setMenuOpen] = useState(false)
  const [menuPos, setMenuPos] = useState({ x: 0, y: 0 })
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(message.content)
  const [loading, setLoading] = useState(false)
  const [deleteDialog, setDeleteDialog] = useState<1 | 2 | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
 
  const direction = isOutgoing ? 'out' : 'in'
 
  useEffect(() => {
    if (!menuOpen) return
    function handleClick(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [menuOpen])
 
  useEffect(() => {
    if (editing) inputRef.current?.focus()
  }, [editing])
 
  function handleContextMenu(e: React.MouseEvent) {
    if (!isOutgoing || message.is_deleted) return
    e.preventDefault()
    
    const x = Math.max(8, e.clientX - MENU_WIDTH)
    const y = e.clientY
    setMenuPos({ x, y })
    setMenuOpen(true)
  }
 
  function handleEdit() {
    setMenuOpen(false)
    setDraft(message.content)
    setEditing(true)
  }
 
  async function handleEditConfirm() {
    if (!draft.trim() || draft.trim() === message.content) {
      setEditing(false)
      return
    }
    setLoading(true)
    try {
      await onUpdate(message.id, { content: draft.trim() })
    } finally {
      setLoading(false)
      setEditing(false)
    }
  }
 
  function handleEditKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter') handleEditConfirm()
    if (e.key === 'Escape') setEditing(false)
  }
 
  async function handleDeleteConfirm() {
    if (deleteDialog === null) return
    setLoading(true)
    try {
      await onUpdate(message.id, { delete: deleteDialog })
    } finally {
      setLoading(false)
      setDeleteDialog(null)
    }
  }
 
  if (message.is_deleted) {
    return (
      <div className={`msg-bubble-row msg-bubble-row--${direction}`}>
        <div className="msg-bubble msg-bubble--deleted">
          <p className="msg-bubble-text">Message deleted</p>
          <p className="msg-bubble-time">{formatTime(message.sent_at)}</p>
        </div>
      </div>
    )
  }
 
  return (
    <>
      <div
        className={`msg-bubble-row msg-bubble-row--${direction}`}
        onContextMenu={handleContextMenu}
        style={{ opacity: loading ? 0.5 : 1 }}
      >
        <div className={`msg-bubble msg-bubble--${direction}`}>
          {editing ? (
            <div className="msg-edit-row">
              <input
                ref={inputRef}
                className="msg-edit-input"
                value={draft}
                onChange={e => setDraft(e.target.value)}
                onKeyDown={handleEditKeyDown}
              />
              <button className="msg-edit-confirm" onClick={handleEditConfirm}>✓</button>
              <button className="msg-edit-cancel" onClick={() => setEditing(false)}>✕</button>
            </div>
          ) : (
            <p className="msg-bubble-text">{message.content}</p>
          )}
          <p className="msg-bubble-time">{formatTime(message.sent_at)}</p>
        </div>
      </div>
 
      {/* Context menu */}
      {menuOpen && (
        <div
          ref={menuRef}
          className="msg-context-menu"
          style={{ top: menuPos.y, left: menuPos.x }}
        >
          <button className="msg-context-item" onClick={handleEdit}>
            Edit message
          </button>
          <button
            className="msg-context-item msg-context-item--danger"
            onClick={() => { setMenuOpen(false); setDeleteDialog(1) }}
          >
            Delete for me
          </button>
          <button
            className="msg-context-item msg-context-item--danger"
            onClick={() => { setMenuOpen(false); setDeleteDialog(2) }}
          >
            Delete for everyone
          </button>
        </div>
      )}
 
      {/* Delete confirmation dialog */}
      {deleteDialog !== null && (
        <div className="msg-dialog-overlay" onClick={() => setDeleteDialog(null)}>
          <div className="msg-dialog" onClick={e => e.stopPropagation()}>
            <h3 className="msg-dialog-title">Delete message?</h3>
            <p className="msg-dialog-body">
              {deleteDialog === 1
                ? 'This message will be removed from your view only. Others can still see it.'
                : 'This message will be permanently deleted for everyone in this conversation.'}
            </p>
            <div className="msg-dialog-actions">
              <button className="msg-dialog-cancel" onClick={() => setDeleteDialog(null)}>
                Cancel
              </button>
              <button className="msg-dialog-confirm" onClick={handleDeleteConfirm}>
                {deleteDialog === 1 ? 'Delete for me' : 'Delete for everyone'}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
 