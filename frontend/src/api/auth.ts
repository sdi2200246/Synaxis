import api from './client'
import type { LoginCredentials, RegisterPayload, LoginResponse } from '../types'

export async function login(credentials: LoginCredentials): Promise<string> {
  const response = await api.post<LoginResponse>('/auth/login', credentials)
  return response.data.token
}

export async function register(payload: RegisterPayload): Promise<void> {
  await api.post('/users', payload)
}
