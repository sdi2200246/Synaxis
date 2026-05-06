import api from './client'
import type {UserSummary} from '../types'


export interface UserListResponse {
  count: number
  users: UserSummary[]
}

export type ExportFormat = 'xml' | 'json'

const EXPORT_ACCEPT: Record<ExportFormat, string> = {
  xml: 'application/xml',
  json: 'application/json',
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

export async function exportUserEvents(userID: string, format: ExportFormat): Promise<void> {
  const res = await api.get(`admin/events?organizer_id=${userID}`, {
    headers: { Accept: EXPORT_ACCEPT[format] },
    responseType: 'blob',
  })

  const disposition = res.headers['content-disposition'] as string | undefined
  const match = disposition?.match(/filename="?([^"]+)"?/)
  const filename = match?.[1] ?? `events-${userID}.${format}`

  const url = URL.createObjectURL(res.data)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}