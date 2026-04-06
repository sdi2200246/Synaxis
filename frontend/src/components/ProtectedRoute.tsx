import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

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
    return <Navigate to="/login" replace />
  }

  if (role !== userRole) {
    return <Navigate to="/home" replace />
  }

  return <>{children}</>
}
