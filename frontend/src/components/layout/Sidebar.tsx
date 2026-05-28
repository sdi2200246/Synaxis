import { useState, useEffect } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import {
  FiHome, FiSearch, FiCalendar, FiCheckSquare, FiUser, FiUsers, FiBell,
  FiMessageSquare, FiChevronLeft, FiChevronRight, FiLogOut,
} from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'
import { useMessages } from '../../context/MessagesContext'
import { getPendingUsers } from '../../api/users'

interface Props {
  collapsed: boolean
  onToggle: () => void
}

export function Sidebar({ collapsed, onToggle }: Props) {
  const [pendingUsersCount, setPendingCount] = useState<number>(0)
  const { logout, userRole } = useAuth()
  const navigate = useNavigate()
  const { hasUnread } = useMessages()

  function handleLogout() {
    logout()
    navigate('/login')
  }

  useEffect(() => {
    if (userRole != 'admin') return

    async function loadPendingUsers() {
      const data = await getPendingUsers()
      setPendingCount(data.count)
    }

    loadPendingUsers()
    const interval = setInterval(loadPendingUsers, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <>
      <aside className={`sidebar ${collapsed ? 'is-collapsed' : ''}`}>
        <nav className="sidebar__nav">
          {userRole === 'admin' ? (
            <>
              <NavLink to="/admin/users" className="sidebar__link" title="Users">
                <FiUsers size={20} />
                <span className="sidebar__label">Users</span>
              </NavLink>
              <NavLink to="/admin/registrations" className="sidebar__link" title="Registrations">
                <span className="icon-badge">
                  <FiBell size={20} />
                  {pendingUsersCount > 0 && <span className="icon-badge__dot" />}
                </span>
                <span className="sidebar__label">Registrations</span>
              </NavLink>
            </>
          ) : (
            <>
              <NavLink to="/browse" className="sidebar__link" title="Home">
                <FiHome size={20} />
                <span className="sidebar__label">Home</span>
              </NavLink>

              <NavLink to="/search" className="sidebar__link" title="Search">
                <FiSearch size={20} />
                <span className="sidebar__label">Search</span>
              </NavLink>

              <NavLink to="/my-events" className="sidebar__link" title="My Events">
                <FiCalendar size={20} />
                <span className="sidebar__label">My Events</span>
              </NavLink>

              <NavLink to="/attending" className="sidebar__link" title="Attending">
                <FiCheckSquare size={20} />
                <span className="sidebar__label">Attending</span>
              </NavLink>

              <NavLink to="/messages" className="sidebar__link" title="Messages">
                <span className="icon-badge">
                  <FiMessageSquare size={20} />
                  {hasUnread && <span className="icon-badge__dot" />}
                </span>
                <span className="sidebar__label">Messages</span>
              </NavLink>
            </>
          )}
        </nav>

        <div className="sidebar__footer">
          <button
            type="button"
            onClick={handleLogout}
            className="sidebar__link sidebar__logout"
            title="Logout"
          >
            <FiLogOut size={20} />
            <span className="sidebar__label"> {userRole ? "Logout" : "Login"} </span>
          </button>
        </div>
      </aside>

      <button
        type="button"
        className="sidebar-toggle-edge"
        onClick={onToggle}
        aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        title={collapsed ? 'Expand' : 'Collapse'}
      >
        {collapsed ? <FiChevronRight size={14} /> : <FiChevronLeft size={14} />}
      </button>
    </>
  )
}