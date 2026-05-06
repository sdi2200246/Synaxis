import React from 'react'
import { FiLock } from 'react-icons/fi'

export type GuestGateProps = {
  title?: string
  subtitle?: string
  loginHref: string
  registerHref: string
  loginLabel?: string
  registerLabel?: string
  compact?: boolean
  className?: string
  icon?: React.ReactNode
}

export function GuestGate({
  title = 'Sign in to continue',
  subtitle = 'This action is available to registered users only. Log in or create an account to keep going.',
  loginHref,
  registerHref,
  loginLabel = 'Log in',
  registerLabel = 'Register',
  compact = false,
  className = '',
  icon,
}: GuestGateProps) {
  return (
    <div className={["guest-gate", compact ? 'guest-gate--compact' : '', className].filter(Boolean).join(' ')}>
      <div className="guest-gate__icon" aria-hidden="true">
        {icon ?? <FiLock size={20} />}
      </div>

      <div className="guest-gate__copy">
        <h3 className="guest-gate__title">{title}</h3>
        <p className="guest-gate__subtitle">{subtitle}</p>
      </div>

      <div className="guest-gate__actions">
        <a className="btn btn--secondary" href={loginHref}>
          {loginLabel}
        </a>
        <a className="btn btn--primary" href={registerHref}>
          {registerLabel}
        </a>
      </div>
    </div>
  )
}