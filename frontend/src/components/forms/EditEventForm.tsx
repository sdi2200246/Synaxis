import { useState } from 'react'
import { FiX } from 'react-icons/fi'
import { updateEvent, getEvent } from '../../api/events'
import { uploadEventMedia, deleteEventMedia } from '../../api/media'
import type { MediaItem } from '../../api/media'
import { useStaticData } from '../../context/StaticData'
import type { Event } from '../../types'

interface EditEventFormProps {
  event: Event
  onClose: () => void
  onSuccess: (updated?: Event) => void
}

export function EditEventForm({ event: initialEvent, onClose, onSuccess }: EditEventFormProps) {
  const { venues, categories, loading } = useStaticData()
  const [event, setEvent] = useState(initialEvent)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [photoMessage, setPhotoMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)
  const [photoUrl, setPhotoUrl] = useState<string | null>(event.media?.[0]?.url ?? null)
  const [photoId, setPhotoId] = useState<string | null>(event.media?.[0]?.id ?? null)
  const [photoUploading, setPhotoUploading] = useState(false)

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

  async function refreshEvent() {
    try {
      const updated = await getEvent(event.id)
      setEvent(updated)
      setPhotoUrl(updated.media?.[0]?.url ?? null)
      setPhotoId(updated.media?.[0]?.id ?? null)
    } catch {}
  }

  async function handlePhotoChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    e.target.value = ''

    setPhotoUploading(true)
    setPhotoMessage(null)
    try {
      if (photoId) {
        await deleteEventMedia(event.id, photoId)
      }
      const res: MediaItem = await uploadEventMedia(event.id, file)
      setPhotoUrl(res.url)
      setPhotoId(res.id)
      setPhotoMessage({ type: 'success', text: 'Photo uploaded successfully' })
    } catch (err: any) {
      setPhotoMessage({ type: 'error', text: err.response?.data?.error ?? 'Failed to upload photo' })
    } finally {
      setPhotoUploading(false)
    }
  }

  async function handlePhotoRemove() {
    if (!photoId) return
    setPhotoUploading(true)
    setPhotoMessage(null)
    try {
      await deleteEventMedia(event.id, photoId)
      setPhotoUrl(null)
      setPhotoId(null)
      setPhotoMessage({ type: 'success', text: 'Photo removed successfully' })
    } catch (err: any) {
      setPhotoMessage({ type: 'error', text: err.response?.data?.error ?? 'Failed to remove photo' })
    } finally {
      setPhotoUploading(false)
    }
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
      await refreshEvent()
      onSuccess(event)
    } catch (err: any) {
      setError(err.response?.data?.error ?? 'Failed to update event')
    } finally {
      setSubmitting(false)
    }
  }

  const formatDatetime = (dt: string) =>
    new Date(dt).toLocaleString('en-US', { dateStyle: 'medium', timeStyle: 'short' })

  return (
    <div className="overlay" onClick={onClose}>
      <div className="dialog dialog--wide" onClick={e => e.stopPropagation()}>
        <div className="dialog__header">
          <h2 className="dialog__title">Edit Event</h2>
          <button className="btn btn--icon" onClick={onClose} type="button">
            <FiX size={24} />
          </button>
        </div>

        <div className="dialog__body">
          <form onSubmit={handleSubmit} className="event-form">
            {error && <div className="alert alert--error">{error}</div>}

            <div className="form-grid">
              <div className="field">
                <label className="field__label" htmlFor="title">Event Title</label>
                <input
                  className="field__control" id="title" name="title" type="text"
                  value={form.title} onChange={handleChange} required disabled={submitting}
                />
              </div>

              <div className="field">
                <label className="field__label" htmlFor="event_type">Event Type</label>
                <input
                  className="field__control" id="event_type" name="event_type" type="text"
                  value={form.event_type} onChange={handleChange} required disabled={submitting}
                />
              </div>

              <div className="field">
                <label className="field__label" htmlFor="capacity">Capacity</label>
                <input className="field__control" id="capacity" type="number" value={event.capacity} disabled />
              </div>

              <div className="field">
                <label className="field__label" htmlFor="venue_id">Venue</label>
                <div className="venue-row">
                  <select
                    className="field__control" id="venue_id" name="venue_id"
                    value={form.venue_id} onChange={handleChange}
                    disabled={loading || submitting} required
                  >
                    <option value="">{loading ? 'Loading venues…' : 'Select a venue'}</option>
                    {venues.map(v => (
                      <option key={v.id} value={v.id}>
                        {v.name} — {v.city}, {v.country}
                        {v.capacity != null ? ` (cap. ${v.capacity})` : ''}
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

              <div className="field field--full">
                <label className="field__label" htmlFor="categories">Categories</label>
                <select
                  className="field__control" id="categories"
                  onChange={handleCategoryAdd} disabled={loading || submitting}
                >
                  <option value="">{loading ? 'Loading categories…' : 'Add a category'}</option>
                  {categories
                    .filter(c => !form.category_ids.includes(c.id))
                    .map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
                {form.category_ids.length > 0 && (
                  <div className="chip-list">
                    {form.category_ids.map(id => (
                      <div key={id} className="chip">
                        <span>{getCategoryName(id)}</span>
                        <button type="button" className="chip__remove" onClick={() => handleCategoryRemove(id)} disabled={submitting}>×</button>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="field field--full">
                <label className="field__label" htmlFor="description">Description</label>
                <textarea
                  className="field__control" id="description" name="description"
                  value={form.description} onChange={handleChange}
                  rows={4} required disabled={submitting}
                />
              </div>

              <div className="field">
                <label className="field__label">Start Date & Time</label>
                <input className="field__control" type="text" value={formatDatetime(event.start_datetime)} disabled />
              </div>

              <div className="field">
                <label className="field__label">End Date & Time</label>
                <input className="field__control" type="text" value={formatDatetime(event.end_datetime)} disabled />
              </div>

              {/* ── Photo ───────────────────────────────────────────── */}
              <div className="field field--full">
                <label className="field__label">Photo</label>

                {photoMessage && (
                  <div className={`alert alert--${photoMessage.type === 'success' ? 'success' : 'error'}`}>
                    {photoMessage.text}
                  </div>
                )}

                {photoUrl ? (
                  <div className="photo-edit">
                    <div className="photo-edit__preview">
                      <img className="photo-edit__img" src={photoUrl} alt="event photo" />
                      <button
                        type="button"
                        className="photo-edit__remove"
                        onClick={handlePhotoRemove}
                        disabled={photoUploading || !photoId}
                        aria-label="Remove photo"
                      >
                        <FiX size={14} />
                      </button>
                    </div>
                    <p className="field__hintp">
                      Only one photo per event. Remove the current photo to upload a new one.
                    </p>
                  </div>
                ) : (
                  <>
                    <label
                      className={`btn btn--ghost ${photoUploading ? 'is-disabled' : ''}`}
                      htmlFor="photo-upload"
                    >
                      {photoUploading ? 'Uploading…' : '+ Add Photo'}
                    </label>
                    <p className="field__hintp">
                      JPG, PNG, or WebP up to 5MB. One photo per event.
                    </p>
                  </>
                )}

                <input
                  id="photo-upload"
                  type="file"
                  accept=".jpg,.jpeg,.png,.webp"
                  style={{ display: 'none' }}
                  onChange={handlePhotoChange}
                  disabled={photoUploading || submitting}
                />
              </div>
            </div>

            <div className="dialog__actions dialog__actions--with-divider">
              <button type="button" className="btn btn--ghost" onClick={onClose} disabled={submitting}>
                Cancel
              </button>
              <button type="submit" className="btn btn--primary" disabled={loading || submitting}>
                {submitting ? 'Saving…' : 'Save Changes'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}