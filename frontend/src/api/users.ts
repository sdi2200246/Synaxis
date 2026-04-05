import api from './client'
import type { User } from '../types'


export interface UserListResponse {
  count: number
  users: User[]
}

export async function getPendingUsers(): Promise<UserListResponse> {
  const response = await api.get<UserListResponse>('admin/users?status=pending')
  return response.data
}