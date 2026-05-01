import { useState, useEffect, useCallback } from 'react'
import { UserCard } from '../../components/UserCard'
import { getPendingUsers, approveUser, rejectUser } from '../../api/users'
import type { UserSummary } from '../../types'

export function PendingRegistrations() {
  const [users, setUsers] = useState<UserSummary[]>([])
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [actionInFlight, setActionInFlight] = useState<string | null>(null)

  const loadUsers = useCallback(async () => {
    try {
      const data = await getPendingUsers()
      setUsers(data.users)
      setCount(data.count)
    } catch {
      console.error('Failed to load pending users')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadUsers()
  }, [loadUsers])

  async function handleApprove(id: string) {
    setActionInFlight(id)
    try {
      await approveUser(id)
      setUsers((prev) => prev.filter((u) => u.id !== id))
      setCount((prev) => prev - 1)
    } catch {
      console.error('Failed to approve user')
    } finally {
      setActionInFlight(null)
    }
  }

  async function handleReject(id: string) {
    setActionInFlight(id)
    try {
      await rejectUser(id)
      setUsers((prev) => prev.filter((u) => u.id !== id))
      setCount((prev) => prev - 1)
    } catch {
      console.error('Failed to reject user')
    } finally {
      setActionInFlight(null)
    }
  }

  if (loading) {
    return <div className="page-narrow"><p>Loading...</p></div>
  }

  return (
    <div className="page-narrow">
      <div className="pending-header">
        <h1>Pending Registrations</h1>
        <span className="badge badge--neutral">{count} pending</span>
      </div>

      {users.length === 0 ? (
        <div className="empty-state">
          <p>No pending registrations — you're all caught up.</p>
        </div>
      ) : (
        <div className="list-stack">
          {users.map((user) => (
            <UserCard
              key={user.id}
              user={{
                id: user.id,
                username: user.username,
                first_name: user.first_name,
                last_name: user.last_name,
                email: user.email,
                phone: user.phone,
                city: user.city,
                country: user.country,
              }}
              variant='compact'
              actions={
                <>
                  <button
                    className="btn btn--success"
                    disabled={actionInFlight === user.id}
                    onClick={() => handleApprove(user.id)}
                  >
                    Approve
                  </button>
                  <button
                    className="btn btn--danger"
                    disabled={actionInFlight === user.id}
                    onClick={() => handleReject(user.id)}
                  >
                    Reject
                  </button>
                </>
              }
            />
          ))}
        </div>
      )}

      <footer className="pending-footer">
        <p>
          Approved users gain full access to create and attend events.
          Rejected accounts are notified by email and may re-apply.
        </p>
      </footer>
    </div>
  )
}