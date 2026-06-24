import { useNavigate } from 'react-router-dom'

export function WelcomePage() {
  const navigate = useNavigate()

  return (
    <div className="welcome">
      <div className="welcome__content">
        <h1 className="welcome__brand">Synaxis</h1>
        <p className="welcome__tagline">Welcome</p>

        <div className="welcome__actions">
          <button
            type="button"
            className="welcome__btn welcome__btn--primary"
            onClick={() => navigate('/register')}
          >
            Create account
          </button>
          <button
            type="button"
            className="welcome__btn welcome__btn--secondary"
            onClick={() => navigate('/login')}
          >
            Log in
          </button>
          <button
            type="button"
            className="welcome__btn welcome__btn--ghost"
            onClick={() => navigate('/browse')}
          >
            Continue as guest
          </button>
        </div>
      </div>
    </div>
  )
}