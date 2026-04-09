import api from './client'
import type { Venue } from '../types'


export interface VenueListResponse {
  count: number
  venues: Venue[]
}

export async function getVenues(): Promise<Venue[]> {
  const response = await api.get<VenueListResponse>('/venues')
  return response.data.venues
}