import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { StaticDataProvider } from './context/StaticData'
import { MessagesProvider } from './context/MessagesContext'
import { Layout } from './components/layout/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { 
  LoginPage,
  RegisterPage,
  MyEventsPage,
  PendingRegistrations,
  Users,
  EventTicketsPage,
  BrowsePage,
  SearchPage,
  AttendingPage,
  EventBookingsPage ,
  MessagesPage,
  WelcomePage,
} from './pages'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <StaticDataProvider>
          <MessagesProvider>
              <Routes>
                <Route path="/" element={<WelcomePage/>} />
                <Route element={<Layout />}>
                  {/* Public routes */}
                  <Route path="/welcome" element = {<WelcomePage />} />
                  <Route path="/login" element={<LoginPage />} />
                  <Route path="/register" element={<RegisterPage />} />
                 
                  {/* Protected routes */}
          
                  <Route
                    path="/browse"
                    element={
                        <BrowsePage/>
                    }
                  />

                  <Route path="/search" element={<SearchPage />} />

                  <Route
                    path="/my-events"
                    element={
                      <ProtectedRoute>
                        <MyEventsPage />
                      </ProtectedRoute>
                    }
                  />
                  
                  <Route
                    path="/attending"
                    element={
                      <ProtectedRoute>
                        <AttendingPage />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/my-events/:id/bookings"
                    element={
                      <ProtectedRoute>
                        <EventBookingsPage />
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
                  <Route
                    path="/events/:id/tickets"
                    element={
                      <ProtectedRoute>
                        <EventTicketsPage />
                      </ProtectedRoute>
                    }
                  />

                  <Route path="/messages" element={<ProtectedRoute><MessagesPage /></ProtectedRoute>} />
                  <Route path="/messages/:conversationId" element={<ProtectedRoute><MessagesPage /></ProtectedRoute>} />

                  {/* Default redirect */}
                  <Route path="*" element={<WelcomePage />} />
                </Route>
              </Routes>
            </MessagesProvider>
          </StaticDataProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
