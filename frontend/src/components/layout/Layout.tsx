import { useState, useEffect } from 'react'
import {Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Topbar } from './Topbar'
import { useLocation } from "react-router-dom";

const SIDEBAR_KEY = 'synaxis-sidebar-collapsed'

export function Layout() {
  const { pathname } = useLocation();
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    return localStorage.getItem(SIDEBAR_KEY) === '1'
  })
 const hideSidebar = pathname === "/login" || pathname === "/register" || pathname === "/welcome";
  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0')
  }, [collapsed])

    return (
     <> 
        <Topbar/>
        <div className={`app-shell ${collapsed ? 'app-shell--collapsed' : ''}`}>
         {!hideSidebar && <Sidebar collapsed={collapsed} onToggle={() => setCollapsed(c => !c)} />}
          <main className="main-content">
            <Outlet />
          </main>
        </div>
      </>
    )
}