import { useState, useEffect ,} from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { FiHome, FiSearch, FiCalendar, FiCheckSquare, FiUser, FiUsers, FiBell } from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'
import { getPendingUsers} from '../../api/users'

export function Sidebar() {
  const [pendingUsersCount, setPendingCount] = useState<number>(0)

  const { logout, userRole } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login')
  }
  
  useEffect(() => {
    async function loadPendingUsers() {
      const data = await getPendingUsers()
      setPendingCount(data.count)
    }

    loadPendingUsers()
      const interval = setInterval(loadPendingUsers, 30000)

    return () => clearInterval(interval)
  }, [])

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <h2>Synaxis</h2>
      </div>

      <nav className="sidebar-nav">
        {userRole == 'admin' ? (
          <>
            <NavLink to="/admin/users" className="sidebar-link">
              <FiUsers size={20} />
              <span>Users</span>
            </NavLink>
            <NavLink to="/admin/registrations" className="sidebar-link">
              <span className="icon-badge">
                <FiBell size={20} />
                {pendingUsersCount > 0 && <span className="red-dot" />}
              </span>
              <span>Registrations</span>
            </NavLink>
          </>
        ) : (
          <>
            <NavLink to="/home" className="sidebar-link">
              <FiHome size={20} />
              <span>Home</span>
            </NavLink>
            <NavLink to="/browse" className="sidebar-link">
              <FiSearch size={20} />
              <span>Browse</span>
            </NavLink>
            <NavLink to="/my-events" className="sidebar-link">
              <FiCalendar size={20} />
              <span>My Events</span>
            </NavLink>
            <NavLink to="/attending" className="sidebar-link">
              <FiCheckSquare size={20} />
              <span>Attending</span>
            </NavLink>
            <NavLink to="/profile" className="sidebar-link">
              <FiUser size={20} />
              <span>Profile</span>
            </NavLink>
          </>
        )}
      </nav>

      <div className="sidebar-footer">
        <button onClick={handleLogout} className="logout-btn">
          Logout
        </button>
      </div>
    </aside>
  )
}