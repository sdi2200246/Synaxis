interface Props {
  title: string
  body: string
  confirmLabel?: string
  cancelLabel?: string
  loading?: boolean
  error?: string
  onConfirm: () => void
  onCancel: () => void
  overlayClassName?: string
  dialogClassName?: string
  confirmClassName?: string
  cancelClassName?: string
}

export function ConfirmDialog({
  title, body, confirmLabel = 'Confirm', cancelLabel = 'Cancel',
  loading = false, error, onConfirm, onCancel,
  overlayClassName = 'browse-detail-overlay',
  dialogClassName = 'browse-detail',
  confirmClassName = 'browse-detail__confirm-btn',
  cancelClassName = 'browse-detail__btn',
}: Props) {
  return (
    <div className={overlayClassName} onClick={() => !loading && onCancel()}>
      <div className={dialogClassName} style={{ maxWidth: '400px' }} onClick={e => e.stopPropagation()}>
        <div className="browse-detail__content">
          <h2 className="browse-detail__title">{title}</h2>
          <div className="browse-detail__confirm">
            <p className="browse-detail__confirm-text">{body}</p>
            {error && <p className="browse-detail__confirm-warning">{error}</p>}
            <div className="browse-detail__confirm-actions">
              <button className={cancelClassName} onClick={onCancel} disabled={loading}>
                {cancelLabel}
              </button>
              <button className={confirmClassName} onClick={onConfirm} disabled={loading}>
                {loading ? 'Working…' : confirmLabel}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}