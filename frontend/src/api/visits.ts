import api from './client'

export async function recordVisit(id:string) : Promise<void>{
    await api.post(`/events/${id}/visits`)
}