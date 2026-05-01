import { useState, useEffect } from 'react'
import { Link, Outlet } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { Sidebar } from './Sidebar'

const SIDEBAR_KEY = 'synaxis-sidebar-collapsed'

export function Layout() {
  const { isAuthenticated } = useAuth()
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    return localStorage.getItem(SIDEBAR_KEY) === '1'
  })

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0')
  }, [collapsed])

  if (isAuthenticated) {
    return (
      <div className={`app-shell ${collapsed ? 'app-shell--collapsed' : ''}`}>
        <Sidebar collapsed={collapsed} onToggle={() => setCollapsed(c => !c)} />
        <main className="main-content">
          <Outlet />
        </main>
      </div>
    )
  }

  return (
    <div className="app">
      <nav className="navbar">
        <Link to="/home" className="navbar__brand">Synaxis</Link>

        <div className="navbar__links">
          <Link to="/login" className="navbar__link">Login</Link>
          <Link to="/register" className="navbar__link">Register</Link>
        </div>
      </nav>

      <main className="public-content">
        <Outlet />
      </main>
    </div>
  )
}