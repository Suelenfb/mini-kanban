import React, { useState } from 'react'

export default function TaskModal({ task, onClose, onSave }) {
  const [title, setTitle] = useState(task?.title || '')
  const [description, setDescription] = useState(task?.description || '')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState(null)

  const isEditing = Boolean(task)

  async function handleSubmit(e) {
    e.preventDefault()
    if (!title.trim()) {
      setError('O título é obrigatório.')
      return
    }

    try {
      setSaving(true)
      await onSave({ title: title.trim(), description: description.trim() })
    } catch (err) {
      setError(err.message || 'Erro ao salvar a tarefa.')
      setSaving(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>{isEditing ? 'Editar tarefa' : 'Nova tarefa'}</h2>

        <form onSubmit={handleSubmit}>
          <label htmlFor="title">Título</label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Ex: Revisar pull request"
            autoFocus
          />

          <label htmlFor="description">Descrição</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Detalhes da tarefa (opcional)"
            rows={4}
          />

          {error && <p className="form-error">{error}</p>}

          <div className="modal-actions">
            <button type="button" className="btn btn--ghost" onClick={onClose}>
              Cancelar
            </button>
            <button type="submit" className="btn btn--primary" disabled={saving}>
              {saving ? 'Salvando...' : 'Salvar'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
