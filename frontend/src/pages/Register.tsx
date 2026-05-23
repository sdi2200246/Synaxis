import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function RegisterPage() {
  const [form, setForm] = useState({
    username: '',
    password: '',
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    country: '',
    tax_id: '',
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const { register } = useAuth()
  const navigate = useNavigate()

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setIsSubmitting(true)

    try {
      await register(form)
      setSuccess(true)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (success) {
    return (
      <div className="auth-page">
        <h1 className="auth-page__title">Registration Submitted</h1>
        <p className="auth-page__footer">
          Your account is pending admin approval.
        </p>
        <Link to="/login" className="btn btn--primary btn--block">
          Back to Login
        </Link>
      </div>
    )
  }

  return (
    <div className="auth-page">
      <h1 className="auth-page__title">Register</h1>

      <form className="auth-page__form" onSubmit={handleSubmit}>
        {error && <div className="alert alert--error">{error}</div>}

        <div className="field">
          <label className="field__label" htmlFor="username">Username</label>
          <input
            className="field__control"
            id="username"
            name="username"
            value={form.username}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="password">Password</label>
          <input
            className="field__control"
            id="password"
            name="password"
            type="password"
            value={form.password}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="first_name">First Name</label>
          <input
            className="field__control"
            id="first_name"
            name="first_name"
            value={form.first_name}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="last_name">Last Name</label>
          <input
            className="field__control"
            id="last_name"
            name="last_name"
            value={form.last_name}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="email">Email</label>
          <input
            className="field__control"
            id="email"
            name="email"
            type="email"
            value={form.email}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="phone">Phone</label>
          <input
            className="field__control"
            id="phone"
            name="phone"
            value={form.phone}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="address">Address</label>
          <input
            className="field__control"
            id="address"
            name="address"
            value={form.address}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="city">City</label>
          <input
            className="field__control"
            id="city"
            name="city"
            value={form.city}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="country">Country</label>
          <input
            className="field__control"
            id="country"
            name="country"
            value={form.country}
            onChange={handleChange}
            required
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="tax_id">Tax ID</label>
          <input
            className="field__control"
            id="tax_id"
            name="tax_id"
            value={form.tax_id}
            onChange={handleChange}
            required
          />
        </div>

        <button
          className="btn btn--primary btn--block"
          type="submit"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Registering...' : 'Register'}
        </button>
      </form>

      <p className="auth-page__footer">
        Already have an account? <Link to="/login">Login</Link>
      </p>
    </div>
  )
}