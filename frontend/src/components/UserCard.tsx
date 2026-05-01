import { useState } from 'react'
import type { UserSummary } from "../types"

type UserCardProps = {
  user: UserSummary
  variant?: 'full' | 'compact'
  actions?: React.ReactNode
}

export function UserCard({ user, variant = 'full', actions }: UserCardProps) {
  const [expanded, setExpanded] = useState(false)
  const clickable = variant === 'full'

  return (
    <div
      className="card user-card"
      onClick={clickable ? () => setExpanded(!expanded) : undefined}
      style={clickable ? { cursor: 'pointer' } : undefined}
    >
      <div className="card__header">
        <div className="card__main">
          <div className="card__title-row">
            <span className="card__title">{user.first_name} {user.last_name}</span>
            {user.status && <span className={`dot dot--${user.status}`} />}
          </div>
          <span className="card__meta">@{user.username}</span>
        </div>
        {actions && (
          <div className="card__actions" onClick={(e) => e.stopPropagation()}>
            {actions}
          </div>
        )}
      </div>

      {expanded && (
        <div className="card__body">
          <div className="form-grid">
            <div className="field">
              <span className="label">Email</span>
              <span>{user.email}</span>
            </div>
            <div className="field">
              <span className="label">Phone</span>
              <span>{user.phone}</span>
            </div>
            <div className="field">
              <span className="label">Address</span>
              <span>{user.address}</span>
            </div>
            <div className="field">
              <span className="label">City</span>
              <span>{user.city}</span>
            </div>
            <div className="field">
              <span className="label">Country</span>
              <span>{user.country}</span>
            </div>
            <div className="field">
              <span className="label">Tax ID</span>
              <span>{user.tax_id}</span>
            </div>
            {user.status && (
              <div className="field">
                <span className="label">Status</span>
                <span>{user.status}</span>
              </div>
            )}
            {user.created_at && (
              <div className="field">
                <span className="label">Registered</span>
                <span>{new Date(user.created_at).toLocaleString('en-US', { dateStyle: 'long' })}</span>
              </div>
            )}
            {user.updated_at && (
              <div className="field">
                <span className="label">Last Updated</span>
                <span>{new Date(user.updated_at).toLocaleString('en-US', { dateStyle: 'long' })}</span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}