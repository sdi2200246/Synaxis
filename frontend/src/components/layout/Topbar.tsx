import { useEffect, useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import api from '../../api/client'
import type { User } from '../../types'

export function Topbar() {
  const { userId } = useAuth()
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    api
      .get<User>(`/users/${userId}`)
      .then((res) => setUser(res.data))
      .catch(() => setUser(null))
  }, [userId])

  const initial = user?.username?.charAt(0).toUpperCase() ?? '?'

  return (
    <header className="topbar">
      <div className="topbar__brand">Synaxis</div>

      <div className="topbar__user" tabIndex={0} aria-label="Account menu">
        <div className="topbar__avatar">{initial}</div>

        <div className="topbar__menu" role="menu">
          {user ? (
            <>
              <div className="topbar__menu-header">
                <div className="topbar__menu-name">
                  {user.first_name} {user.last_name}
                </div>
                <div className="topbar__menu-username">@{user.username}</div>
              </div>

              <div className="topbar__menu-info">
                <div className="topbar__menu-row">
                  <span className="topbar__menu-label">Email</span>
                  <span className="topbar__menu-value">{user.email}</span>
                </div>
                <div className="topbar__menu-row">
                  <span className="topbar__menu-label">Phone</span>
                  <span className="topbar__menu-value">{user.phone}</span>
                </div>
                <div className="topbar__menu-row">
                  <span className="topbar__menu-label">Role</span>
                  <span className="topbar__menu-value">{user.role}</span>
                </div>
              </div>
            </>
          ) : (
            <div className="topbar__menu-loading">Synaxis Welcomes you</div>
          )}
        </div>
      </div>
    </header>
  )
}