// frontend/src/components/layout/Sidebar.tsx
import { NavLink, useNavigate } from 'react-router-dom'
import { FiHome, FiSearch, FiCalendar, FiCheckSquare, FiUser } from 'react-icons/fi'
import { useAuth } from '../../context/AuthContext'

export function Sidebar() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <h2>Synaxis</h2>
      </div>

      <nav className="sidebar-nav">
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
      </nav>

      <div className="sidebar-footer">
        <button onClick={handleLogout} className="logout-btn">
          Logout
        </button>
      </div>
    </aside>
  )
}