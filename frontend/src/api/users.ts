import api from './client'
import type {UserSummary} from '../types'


export interface UserListResponse {
  count: number
  users: UserSummary[]
}

export async function getPendingUsers(): Promise<UserListResponse> {
  const response = await api.get<UserListResponse>('admin/users?status=pending')
  console.log(response.data)
  return response.data
}
export async function getUsers(): Promise<UserListResponse> {
  const response = await api.get<UserListResponse>('admin/users')
  console.log(response.data)
  return response.data
}

export async function approveUser(id: string): Promise<void> {
  await api.post(`/admin/users/${id}/approve`)
}

export async function rejectUser(id: string): Promise<void> {
  await api.post(`/admin/users/${id}/reject`)
}