import { createContext, useContext, useState, useEffect} from 'react'
import type {ReactNode} from 'react'
import { login as apiLogin, register as apiRegister } from '../api/auth'
import type { LoginCredentials, RegisterPayload } from '../types'

interface AuthContextType {
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  userRole : "admin" | "user" | null,
  userId : string|null,
  login: (credentials: LoginCredentials) => Promise<void>
  register: (payload: RegisterPayload) => Promise<void>
  logout: () => void
}

interface TokenClaims {
  user_id: string
  role: "admin" | "user" | null
  exp: number
  iat: number
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Check for existing token on mount
  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setToken(storedToken)
    }
    setIsLoading(false)
  }, [])

  async function login(credentials: LoginCredentials) {
    const newToken = await apiLogin(credentials)
    localStorage.setItem('token', newToken)
    const claims = parseToken(newToken);
    console.log(claims?.role)

    setToken(newToken)
  }

  async function register(payload: RegisterPayload) {
    await apiRegister(payload)
    // Don't auto-login — user needs admin approval
  }

  function logout() {
    localStorage.removeItem('token')
    setToken(null)
  }

  function parseToken(token: string): TokenClaims | null {
    try {
      const payload = token.split('.')[1]
      return JSON.parse(atob(payload))
    } catch {
      return null
    }
  }

  const claims = token ? parseToken(token) : null

  return (
    <AuthContext.Provider
      value={{
        token,
        isAuthenticated: !!token,
        isLoading,
        userRole: claims?.role ?? null,
        userId: claims?.user_id ?? null,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
