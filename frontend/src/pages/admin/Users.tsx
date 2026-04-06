import { useState, useEffect, useCallback } from 'react'
import { UserCard } from '../../components/UserCard'
import { getUsers} from '../../api/users'
import type { UserSummary } from '../../types'

export function Users() {
  const [users, setUsers] = useState<UserSummary[]>([])
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)

  const loadUsers = useCallback(async () => {
    try {
      const data = await getUsers()
      setUsers(data.users)
      setCount(data.count)
    } catch {
      console.error('Failed to load users')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadUsers()
  }, [loadUsers])

  
  if (loading) {
    return <div className="pending-page"><p>Loading...</p></div>
  }

  return (
    <div className="pending-page">
      <div className="pending-header">
        <h1>All Users</h1>
        <span className="count">{count} users</span>
      </div>

      {users.length === 0 ? (
        <div className="pending-empty">
          <p>No user exists in the system</p>
        </div>
      ) : (
        <div className="pending-list">
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
                address:user.address,
                city: user.city,
                country: user.country,
                tax_id: user.tax_id,
                status: user.status,
                created_at :user.created_at,
                updated_at: user?.updated_at,
              }}
            />
          ))}
        </div>
      )}
    </div>
  )
}