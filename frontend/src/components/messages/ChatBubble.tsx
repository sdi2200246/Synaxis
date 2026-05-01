import { useState, useRef, useEffect } from 'react'
import type { Message } from '../../types'
import { ConfirmDialog } from '../ConfirmDialogue'

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
      <div className={`message-row message-row--${direction}`}>
        <div className="message-bubble message-bubble--deleted">
          <p className="message-bubble__text">Message deleted</p>
          <p className="message-bubble__time">{formatTime(message.sent_at)}</p>
        </div>
      </div>
    )
  }

  return (
    <>
      <div
        className={`message-row message-row--${direction}`}
        onContextMenu={handleContextMenu}
        style={{ opacity: loading ? 0.5 : 1 }}
      >
        <div className={`message-bubble message-bubble--${direction}`}>
          {editing ? (
            <div className="message-edit">
              <input
                ref={inputRef}
                className="message-edit__input"
                value={draft}
                onChange={e => setDraft(e.target.value)}
                onKeyDown={handleEditKeyDown}
              />
              <button className="btn btn--icon btn--icon-success" onClick={handleEditConfirm}>✓</button>
              <button className="btn btn--icon" onClick={() => setEditing(false)}>✕</button>
            </div>
          ) : (
            <p className="message-bubble__text">{message.content}</p>
          )}
          <p className="message-bubble__time">{formatTime(message.sent_at)}</p>
        </div>
      </div>

      {/* Context menu */}
      {menuOpen && (
        <div
          ref={menuRef}
          className="context-menu"
          style={{ top: menuPos.y, left: menuPos.x }}
        >
          <button className="context-menu__item" onClick={handleEdit}>
            Edit message
          </button>
          <button
            className="context-menu__item context-menu__item--danger"
            onClick={() => { setMenuOpen(false); setDeleteDialog(1) }}
          >
            Delete for me
          </button>
          <button
            className="context-menu__item context-menu__item--danger"
            onClick={() => { setMenuOpen(false); setDeleteDialog(2) }}
          >
            Delete for everyone
          </button>
        </div>
      )}

      {/* Delete confirmation dialog */}
      {deleteDialog !== null && (
        <ConfirmDialog
          title="Delete message?"
          body={deleteDialog === 1
            ? 'This message will be removed from your view only. Others can still see it.'
            : 'This message will be permanently deleted for everyone in this conversation.'}
          confirmLabel={deleteDialog === 1 ? 'Delete for me' : 'Delete for everyone'}
          variant="danger"
          loading={loading}
          onConfirm={handleDeleteConfirm}
          onCancel={() => setDeleteDialog(null)}
        />
      )}
    </>
  )
}