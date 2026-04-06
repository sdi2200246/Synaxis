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
      className="user-card"
      onClick={clickable ? () => setExpanded(!expanded) : undefined}
      style={clickable ? { cursor: 'pointer' } : undefined}
    >
      <div className="user-info">
        <span className="user-name">{user.first_name} {user.last_name}</span>
        <span className="user-username">@{user.username}</span>
        {user.status && <span className={`status-dot dot-${user.status}`} />}
        {expanded && (
          <div className="user-details">
              <div className="detail-field">
                <label>Email</label>
                <span>{user.email}</span>
              </div>
              <div className="detail-field">
                <label>Phone</label>
                <span>{user.phone}</span>
              </div>
              <div className="detail-field">
                <label>Address</label>
                <span>{user.address}</span>
              </div>
              <div className="detail-field">
                <label>City</label>
                <span>{user.city}</span>
              </div>
              <div className="detail-field">
                <label>Country</label>
                <span>{user.country}</span>
              </div>
              <div className="detail-field">
                <label>Tax ID</label>
                <span>{user.tax_id}</span>
              </div>
              {user.status && (
                <div className="detail-field">
                  <label>Status</label>
                  <span>{user.status}</span>
                </div>
              )}
              {user.created_at && (
                <div className="detail-field">
                  <label>Registered</label>
                  <span>{new Date(user.created_at).toLocaleString('en-US', { dateStyle: 'long' })}</span>
                </div>
              )}
              {user.updated_at && (
              <div className="detail-field">
                <label>Last Updated</label>
                <span>{new Date(user.updated_at).toLocaleString('en-US', { dateStyle: 'long' })}</span>
              </div>
            )}
          </div>
        )}
      </div>
      {actions && <div className="user-actions" onClick={(e) => e.stopPropagation()}>{actions}</div>}
    </div>
  )
}
 