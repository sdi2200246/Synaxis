import { createContext, useContext, useEffect, useState } from 'react'
import { getVenues } from '../api/venues'
import { getCategories } from '../api/category'
import type { Venue, Category } from '../types'
 
interface StaticDataContextType {
  venues: Venue[]
  categories: Category[]
  loading: boolean
  error: string
}
 
const StaticDataContext = createContext<StaticDataContextType>({
  venues: [],
  categories: [],
  loading: true,
  error: '',
})
 
export function StaticDataProvider({ children }: { children: React.ReactNode }) {
  const [venues, setVenues] = useState<Venue[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
 
  useEffect(() => {
    async function fetchStaticData() {
      try {
        const [venuesData, categoriesData] = await Promise.all([
          getVenues(),
          getCategories(),
        ])
        setVenues(venuesData)
        setCategories(categoriesData)
      } catch {
        setError('Failed to load application data')
      } finally {
        setLoading(false)
      }
    }
    fetchStaticData()
  }, [])
 
  return (
    <StaticDataContext.Provider value={{ venues, categories, loading, error }}>
      {children}
    </StaticDataContext.Provider>
  )
}
 
export function useStaticData() {
  return useContext(StaticDataContext)
}