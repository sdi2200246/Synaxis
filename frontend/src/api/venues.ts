import api from './client'
import type { Venue } from '../types'

export async function getVenues(): Promise<Venue[]> {
  const response = await api.get<Venue[]>('/venues')
  return response.data
}