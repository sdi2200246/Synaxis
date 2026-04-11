import { useState} from 'react'
import { FiX } from 'react-icons/fi'
import { updateEvent } from '../../api/events'
import { useStaticData } from '../../context/StaticData'
import type {Event } from '../../types'
 
interface EditEventFormProps {
  event: Event
  onClose: () => void
  onSuccess: () => void
}
 
export function EditEventForm({ event, onClose, onSuccess }: EditEventFormProps) {

const { venues, categories, loading } = useStaticData()
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [photos, setPhotos] = useState<File[]>([])
  const [photoPreviews, setPhotoPreviews] = useState<string[]>([])
 
  const [form, setForm] = useState({
    title: event.title,
    event_type: event.event_type,
    venue_id: event.venue.id,
    description: event.description,
    category_ids: event.categories?.map(c => c.id) ?? [],
  })
 
 
  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) {
    setForm({ ...form, [e.target.name]: e.target.value })
    if (error) setError('')
  }
 
  function handleCategoryAdd(e: React.ChangeEvent<HTMLSelectElement>) {
    const id = e.target.value
    if (id && !form.category_ids.includes(id)) {
      setForm(prev => ({ ...prev, category_ids: [...prev.category_ids, id] }))
    }
    e.target.value = ''
    if (error) setError('')
  }
 
  function handleCategoryRemove(id: string) {
    setForm(prev => ({ ...prev, category_ids: prev.category_ids.filter(c => c !== id) }))
    if (error) setError('')
  }
 
  function getCategoryName(id: string) {
    return categories.find(c => c.id === id)?.name ?? ''
  }
 
  function handlePhotoAdd(e: React.ChangeEvent<HTMLInputElement>) {
    const files = Array.from(e.target.files ?? [])
    if (!files.length) return
    setPhotos(prev => [...prev, ...files])
    files.forEach(file => {
      const reader = new FileReader()
      reader.onload = ev => {
        setPhotoPreviews(prev => [...prev, ev.target?.result as string])
      }
      reader.readAsDataURL(file)
    })
    e.target.value = ''
  }
 
  function handlePhotoRemove(index: number) {
    setPhotos(prev => prev.filter((_, i) => i !== index))
    setPhotoPreviews(prev => prev.filter((_, i) => i !== index))
  }
 
  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
 
    if (form.category_ids.length === 0) {
      setError('Please select at least one category')
      return
    }
 
    const selectedVenue = venues.find(v => v.id === form.venue_id)
    if (selectedVenue?.capacity != null && selectedVenue.capacity < event.capacity) {
      setError(`Venue capacity (${selectedVenue.capacity}) is less than event capacity (${event.capacity})`)
      return
    }
 
    setError('')
    setSubmitting(true)
 
    try {
      await updateEvent(event.id, {
        title: form.title,
        event_type: form.event_type,
        venue_id: form.venue_id,
        description: form.description,
        category_ids: form.category_ids,
      })
      onSuccess()
    } catch (err: any) {
      setError(err.response?.data?.error ?? 'Failed to update event')
    } finally {
      setSubmitting(false)
    }
  }
 
  const formatDatetime = (dt: string) =>
    new Date(dt).toLocaleString('en-US', { dateStyle: 'medium', timeStyle: 'short' })
 
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Edit Event</h2>
          <button className="close-btn" onClick={onClose} type="button">
            <FiX size={24} />
          </button>
        </div>
 
        <form onSubmit={handleSubmit} className="event-form">
          {error && <div className="error-message">{error}</div>}
 
          <div className="form-grid">
            <div className="field">
              <label htmlFor="title">Event Title</label>
              <input
                id="title"
                name="title"
                type="text"
                value={form.title}
                onChange={handleChange}
                required
                disabled={submitting}
              />
            </div>
 
            <div className="field">
              <label htmlFor="event_type">Event Type</label>
              <input
                id="event_type"
                name="event_type"
                type="text"
                value={form.event_type}
                onChange={handleChange}
                required
                disabled={submitting}
              />
            </div>
 
            <div className="field">
              <label htmlFor="capacity">Capacity</label>
              <input
                id="capacity"
                type="number"
                value={event.capacity}
                disabled
              />
            </div>
 
            <div className="field">
              <label htmlFor="venue_id">Venue</label>
              <div className="venue-row">
                <select
                  id="venue_id"
                  name="venue_id"
                  value={form.venue_id}
                  onChange={handleChange}
                  disabled={loading || submitting}
                  required
                >
                  <option value="">{loading ? 'Loading venues…' : 'Select a venue'}</option>
                  {venues.map(v => (
                    <option key={v.id} value={v.id}>
                      {v.name} — {v.city}, {v.country}{v.capacity != null ? ` (cap. ${v.capacity})` : ''}
                    </option>
                  ))}
                </select>
                <span className="venue-info-icon" title="Venue capacity must be ≥ event capacity">
                  <svg width="15" height="15" viewBox="0 0 15 15" fill="none">
                    <circle cx="7.5" cy="7.5" r="6.5" stroke="currentColor" strokeWidth="1" />
                    <path d="M7.5 6.5v4M7.5 4.5v.5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
                  </svg>
                </span>
              </div>
            </div>
 
            <div className="field full-width">
              <label htmlFor="categories">Categories</label>
              <select
                id="categories"
                onChange={handleCategoryAdd}
                disabled={loading || submitting}
              >
                <option value="">{loading ? 'Loading categories…' : 'Add a category'}</option>
                {categories
                  .filter(c => !form.category_ids.includes(c.id))
                  .map(c => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
              </select>
              {form.category_ids.length > 0 && (
                <div className="selected-categories">
                  {form.category_ids.map(id => (
                    <div key={id} className="category-tag">
                      <span>{getCategoryName(id)}</span>
                      <button
                        type="button"
                        className="remove-tag"
                        onClick={() => handleCategoryRemove(id)}
                        disabled={submitting}
                      >×</button>
                    </div>
                  ))}
                </div>
              )}
            </div>
 
            <div className="field full-width">
              <label htmlFor="description">Description</label>
              <textarea
                id="description"
                name="description"
                value={form.description}
                onChange={handleChange}
                rows={4}
                required
                disabled={submitting}
              />
            </div>
 
            <div className="field">
              <label>Start Date & Time</label>
              <input type="text" value={formatDatetime(event.start_datetime)} disabled />
            </div>
 
            <div className="field">
              <label>End Date & Time</label>
              <input type="text" value={formatDatetime(event.end_datetime)} disabled />
            </div>
 
            <div className="field full-width">
              <label>Photos</label>
              <label className="media-upload-zone" htmlFor="photo-upload">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                  <path d="M12 16V8M8 12l4-4 4 4" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
                  <rect x="3" y="3" width="18" height="18" rx="4" stroke="currentColor" strokeWidth="1.2" />
                </svg>
                <span>Click to upload photos</span>
                <small>PNG, JPG up to 10MB</small>
              </label>
              <input
                id="photo-upload"
                type="file"
                accept="image/*"
                multiple
                style={{ display: 'none' }}
                onChange={handlePhotoAdd}
                disabled={submitting}
              />
              {photoPreviews.length > 0 && (
                <div className="media-previews">
                  {photoPreviews.map((src, i) => (
                    <div key={i} className="media-thumb">
                      <img src={src} alt={`photo ${i + 1}`} />
                      <button
                        type="button"
                        className="media-remove"
                        onClick={() => handlePhotoRemove(i)}
                        disabled={submitting}
                      >×</button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
 
          <div className="form-actions">
            <button type="button" className="btn-cancel" onClick={onClose} disabled={submitting}>
              Cancel
            </button>
            <button type="submit" className="btn-submit" disabled={loading || submitting}>
              {submitting ? 'Saving…' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}