import { useState, useEffect } from 'react'
import {Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'

const SIDEBAR_KEY = 'synaxis-sidebar-collapsed'

export function Layout() {
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    return localStorage.getItem(SIDEBAR_KEY) === '1'
  })

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0')
  }, [collapsed])

    return (
      <div className={`app-shell ${collapsed ? 'app-shell--collapsed' : ''}`}>
        <Sidebar collapsed={collapsed} onToggle={() => setCollapsed(c => !c)} />
        <main className="main-content">
          <Outlet />
        </main>
      </div>
    )
}