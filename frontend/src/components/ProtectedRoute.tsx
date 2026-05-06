import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { GuestGate } from './GuestGate'

interface Props {
  children: React.ReactNode
  role?:'admin'|'user'
}

export function ProtectedRoute({ children , role = 'user' }: Props) {
  const { isAuthenticated, isLoading , userRole} = useAuth()
 
  if (isLoading) {
    return <div>Loading...</div>
  }
  

  if (!isAuthenticated) {
    return (
      <div className="guest-gate-page">
        <GuestGate
          loginHref="/login"
          registerHref="/register"
          title="Create an account to get access to this page"
          subtitle="You need an account before you can access"
        />
      </div>
    )
  }

  if (role !== userRole) {
    return <Navigate to="/home" replace />
  }

  return <>{children}</>
}
