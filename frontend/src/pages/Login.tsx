import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  
  const { login , userRole } = useAuth()
  const navigate = useNavigate()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setIsSubmitting(true)

    try {
      await login({ username, password })
      if (userRole == "admin"){
          navigate('/admin/users')
      }
      else{ 
        navigate('/browse')
      }

    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="auth-page">
      <h1 className="auth-page__title">Login</h1>

      <form className="auth-page__form" onSubmit={handleSubmit}>
        {error && <div className="alert alert--error">{error}</div>}

        <div className="field">
          <label className="field__label" htmlFor="username">
            Username
          </label>

          <input
            className="field__control"
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="password">
            Password
          </label>

          <input
            className="field__control"
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button className="btn btn--primary btn--block" type="submit">
          {isSubmitting ? "Logging in..." : "Login"}
        </button>
      </form>

        <div className="auth-page__divider">or</div>

        <Link to="/browse" className="btn btn--ghost btn--block">
          Continue as guest
        </Link>
        
      <p className="auth-page__footer">
        Don't have an account? <Link to="/register">Register</Link>
      </p>
    </div>
  )
}
