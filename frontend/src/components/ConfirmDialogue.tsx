interface Props {
  title: string
  body: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'primary' | 'danger'
  loading?: boolean
  error?: string
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmDialog({
  title,
  body,
  confirmLabel = 'Confirm',
  cancelLabel = 'Cancel',
  variant = 'primary',
  loading = false,
  error,
  onConfirm,
  onCancel,
}: Props) {
  return (
    <div className="overlay" onClick={() => !loading && onCancel()}>
      <div className="dialog" onClick={e => e.stopPropagation()}>
        <h3 className="dialog__title">{title}</h3>
        <div className="dialog__body">{body}</div>
        {error && <div className="alert alert--error" style={{ marginBottom: 12 }}>{error}</div>}
        <div className="dialog__actions">
          <button
            className="btn btn--ghost"
            onClick={onCancel}
            disabled={loading}
          >
            {cancelLabel}
          </button>
          <button
            className={`btn btn--${variant}`}
            onClick={onConfirm}
            disabled={loading}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}