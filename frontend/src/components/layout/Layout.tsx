// frontend/src/components/layout/Layout.tsx
import { Link, Outlet } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { Sidebar } from './Sidebar'

export function Layout() {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated) {
    return (
      <div className="app-authenticated">
        <Sidebar />
        <main className="main-content">
          <Outlet />
        </main>
      </div>
    )
  }

  return (
    <div className="app">
      <nav className="navbar">
        <Link to="/home" className="brand">Synaxis</Link>
        
        <div className="nav-links">
          <Link to="/login">Login</Link>
          <Link to="/register">Register</Link>
        </div>
      </nav>

      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}