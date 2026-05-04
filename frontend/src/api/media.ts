import api from './client'

export interface MediaItem {
  id: string
  event_id: string
  filename: string
  url: string
  uploaded_at: string
}

export async function uploadEventMedia(eventId: string, file: File): Promise<MediaItem> {
  const fd = new FormData()
  fd.append('photo', file)
  const res = await api.post<MediaItem>(`/events/${eventId}/media`, fd)
  return res.data
}

export async function deleteEventMedia(eventId: string, mediaId: string): Promise<void> {
  await api.delete(`/events/${eventId}/media/${mediaId}`)
}