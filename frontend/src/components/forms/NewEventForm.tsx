// frontend/src/components/CreateEventForm.tsx
import { useState} from 'react'
import { FiX } from 'react-icons/fi'
import { createEvent } from '../../api/events'
import { useStaticData } from '../../context/StaticData'

interface CreateEventFormProps {
  onClose: () => void
  onSuccess: () => void
}

export function CreateEventForm({ onClose, onSuccess }: CreateEventFormProps) {

  const { venues, categories, loading } = useStaticData()
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  
  const [form, setForm] = useState({
    title: '',
    event_type: '',
    venue_id: '',
    description: '',
    capacity: 0,
    start_datetime: '',
    end_datetime: '',
    category_ids: [] as string[],
  })

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) {
    setForm({ ...form, [e.target.name]: e.target.value })
    if (error) setError('')
  }


  function handleCategoryAdd(e: React.ChangeEvent<HTMLSelectElement>) {
    const categoryId = e.target.value
    if (categoryId && !form.category_ids.includes(categoryId)) {
      setForm(prev => ({
        ...prev,
        category_ids: [...prev.category_ids, categoryId]
      }))
    }
    // Reset select to placeholder
    e.target.value = ''
    if (error) setError('')
  }

  function handleCategoryRemove(categoryId: string) {
    setForm(prev => ({
      ...prev,
      category_ids: prev.category_ids.filter(id => id !== categoryId)
    }))
    if (error) setError('')
  }


  function getCategoryName(categoryId: string): string {
    return categories.find(c => c.id === categoryId)?.name || ''
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    
    if (form.category_ids.length === 0) {
      setError('Error : Please select at least one category')
      return
    }

    const selectedVenue = venues.find(v => v.id === form.venue_id)
    if (selectedVenue?.capacity != null && selectedVenue.capacity < Number(form.capacity)) {
      setError(`Error:Venue capacity (${selectedVenue.capacity}) is less than event capacity (${form.capacity})`)
      return
    }


    setError('')
    setSubmitting(true)

    try {
      await createEvent(form)
      onSuccess()
    } catch (err: any) {
     const errorMessage = err.response?.data?.error || err.response?.data?.details || 'Failed to create event';
      setError(errorMessage)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="overlay" onClick={onClose}>
      <div
        className="dialog dialog--wide"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="dialog__header">
          <h2 className="dialog__title">Create New Event</h2>

          <button
            className="btn btn--icon"
            onClick={onClose}
            type="button"
          >
            <FiX size={24} />
          </button>
        </div>

        <div className="dialog__body">
          <form onSubmit={handleSubmit} className="event-form">
            {error && (
              <div className="alert alert--error">
                {error}
              </div>
            )}

            <div className="form-grid">

              <div className="field">
                <label
                  className="field__label"
                  htmlFor="title"
                >
                  Event Title
                </label>

                <input
                  className="field__control"
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
                <label
                  className="field__label"
                  htmlFor="event_type"
                >
                  Event Type
                </label>

                <input
                  className="field__control"
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
                <label
                  className="field__label"
                  htmlFor="capacity"
                >
                  Capacity
                </label>

                <input
                  className="field__control"
                  id="capacity"
                  name="capacity"
                  type="number"
                  value={form.capacity}
                  onChange={handleChange}
                  placeholder="e.g., 100"
                  min="1"
                  required
                  disabled={submitting}
                />
              </div>

              <div className="field">
                <label
                  className="field__label"
                  htmlFor="venue_id"
                >
                  Venue
                </label>

                <div className="venue-row">
                  <select
                    className="field__control"
                    id="venue_id"
                    name="venue_id"
                    value={form.venue_id}
                    onChange={handleChange}
                    disabled={loading || submitting}
                    required
                  >
                    <option value="">
                      {loading
                        ? "Loading venues…"
                        : "Select a venue"}
                    </option>

                    {venues.map((v) => (
                      <option key={v.id} value={v.id}>
                        {v.name} — {v.city}, {v.country}
                        {v.capacity != null
                          ? ` (cap. ${v.capacity})`
                          : ""}
                      </option>
                    ))}
                  </select>

                  <span
                    className="venue-info-icon"
                    title="Venue capacity must be ≥ your event capacity"
                  >
                    <svg
                      width="15"
                      height="15"
                      viewBox="0 0 15 15"
                      fill="none"
                    >
                      <circle
                        cx="7.5"
                        cy="7.5"
                        r="6.5"
                        stroke="currentColor"
                        strokeWidth="1"
                      />
                      <path
                        d="M7.5 6.5v4M7.5 4.5v.5"
                        stroke="currentColor"
                        strokeWidth="1.2"
                        strokeLinecap="round"
                      />
                    </svg>
                  </span>
                </div>
              </div>

              <div className="field field--full">
                <label
                  className="field__label"
                  htmlFor="categories"
                >
                  Categories
                </label>

                <select
                  className="field__control"
                  id="categories"
                  onChange={handleCategoryAdd}
                  disabled={loading || submitting}
                >
                  <option value="">
                    {loading
                      ? "Loading categories..."
                      : "Add a category"}
                  </option>

                  {categories
                    .filter(
                      (cat) =>
                        !form.category_ids.includes(cat.id)
                    )
                    .map((category) => (
                      <option
                        key={category.id}
                        value={category.id}
                      >
                        {category.name}
                      </option>
                    ))}
                </select>

                {form.category_ids.length > 0 && (
                  <div className="chip-list">
                    {form.category_ids.map((catId) => (
                      <div
                        key={catId}
                        className="chip"
                      >
                        <span>
                          {getCategoryName(catId)}
                        </span>

                        <button
                          type="button"
                          className="chip__remove"
                          onClick={() =>
                            handleCategoryRemove(catId)
                          }
                          disabled={submitting}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="field field--full">
                <label
                  className="field__label"
                  htmlFor="description"
                >
                  Description
                </label>

                <textarea
                  className="field__control"
                  id="description"
                  name="description"
                  value={form.description}
                  onChange={handleChange}
                  placeholder="Describe your event..."
                  rows={4}
                  required
                  disabled={submitting}
                />
              </div>

              <div className="field">
                <label
                  className="field__label"
                  htmlFor="start_datetime"
                >
                  Start Date & Time
                </label>

                <input
                  className="field__control"
                  id="start_datetime"
                  name="start_datetime"
                  type="datetime-local"
                  value={form.start_datetime}
                  onChange={handleChange}
                  required
                  disabled={submitting}
                />
              </div>

              <div className="field">
                <label
                  className="field__label"
                  htmlFor="end_datetime"
                >
                  End Date & Time
                </label>

                <input
                  className="field__control"
                  id="end_datetime"
                  name="end_datetime"
                  type="datetime-local"
                  value={form.end_datetime}
                  onChange={handleChange}
                  required
                  disabled={submitting}
                />
              </div>

            </div>

            <div className="dialog__actions dialog__actions--with-divider">
              <button
                type="button"
                className="btn btn--ghost"
                onClick={onClose}
                disabled={submitting}
              >
                Cancel
              </button>

              <button
                type="submit"
                className="btn btn--primary"
                disabled={loading || submitting}
              >
                {submitting
                  ? "Creating..."
                  : "Create Draft Event"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}