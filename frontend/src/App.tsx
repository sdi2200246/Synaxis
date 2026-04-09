import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { StaticDataProvider } from './context/StaticData'
import { Layout } from './components/layout/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { 
  LoginPage,
  RegisterPage,
  HomePage,
  MyEventsPage,
  PendingRegistrations,
  Users
} from './pages'
import './styles.css'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <StaticDataProvider>
          <Routes>
            <Route element={<Layout />}>
              {/* Public routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />

              {/* Protected routes */}
      
              <Route
                path="/home"
                element={
                  <ProtectedRoute>
                    <HomePage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/my-events"
                element={
                  <ProtectedRoute>
                    <MyEventsPage />
                  </ProtectedRoute>
                }
              />
              
              <Route
                path="/admin/registrations"
                element={
                  <ProtectedRoute role="admin">
                    <PendingRegistrations/>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/admin/users"
                element={
                  <ProtectedRoute role="admin">
                    <Users/>
                  </ProtectedRoute>
                }
              />

              {/* Default redirect */}
              <Route path="*" element={<LoginPage />} />
            </Route>
          </Routes>
          </StaticDataProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
