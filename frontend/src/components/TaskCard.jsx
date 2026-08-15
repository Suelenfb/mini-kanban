import React, { useState } from 'react'

export default function TaskCard({ task, onEdit, onDelete, onDropBefore }) {
  const [isOver, setIsOver] = useState(false)

  function handleDragStart(e) {
    e.dataTransfer.setData('text/plain', task.id)
    e.dataTransfer.effectAllowed = 'move'
  }

  function handleDragOver(e) {
    e.preventDefault()
    e.stopPropagation()
    setIsOver(true)
  }

  function handleDragLeave(e) {
    e.stopPropagation()
    setIsOver(false)
  }

  function handleDrop(e) {
    e.preventDefault()
    e.stopPropagation()
    setIsOver(false)
    const taskId = e.dataTransfer.getData('text/plain')
    if (taskId && taskId !== task.id) onDropBefore(taskId)
  }

  return (
    <article
      className={`task-card ${isOver ? 'task-card--drag-over' : ''}`}
      draggable
      onDragStart={handleDragStart}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <div className="task-card-header">
        <h3>{task.title}</h3>
        <div className="task-actions">
          <button title="Editar" onClick={onEdit} className="icon-btn">
            ✎
          </button>
          <button title="Excluir" onClick={onDelete} className="icon-btn icon-btn--danger">
            ✕
          </button>
        </div>
      </div>
      {task.description && <p className="task-description">{task.description}</p>}
    </article>
  )
}
