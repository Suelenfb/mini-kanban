import React, { useState } from 'react'
import TaskCard from './TaskCard.jsx'

export default function Column({ status, title, tasks, onAddTask, onEditTask, onDeleteTask, onDropTask }) {
  const [isDragOver, setIsDragOver] = useState(false)

  function handleDragOver(e) {
    e.preventDefault()
    setIsDragOver(true)
  }

  function handleDragLeave() {
    setIsDragOver(false)
  }

  // Solto na área vazia da coluna (ou depois do último card): vai para o final.
  function handleDropOnColumn(e) {
    e.preventDefault()
    setIsDragOver(false)
    const taskId = e.dataTransfer.getData('text/plain')
    if (taskId) onDropTask(taskId, status, tasks.length)
  }

  // Solto sobre um card específico: insere antes dele.
  function handleDropOnCard(taskId, index) {
    setIsDragOver(false)
    onDropTask(taskId, status, index)
  }

  return (
    <section
      className={`column ${isDragOver ? 'column--drag-over' : ''}`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDropOnColumn}
    >
      <div className="column-header">
        <h2>{title}</h2>
        <span className="task-count">{tasks.length}</span>
      </div>

      <div className="column-body">
        {tasks.map((task, index) => (
          <TaskCard
            key={task.id}
            task={task}
            onEdit={() => onEditTask(task)}
            onDelete={() => onDeleteTask(task.id)}
            onDropBefore={(taskId) => handleDropOnCard(taskId, index)}
          />
        ))}

        {tasks.length === 0 && <div className="empty-hint">Nenhuma tarefa aqui ainda</div>}
      </div>

      <button className="add-task-btn" onClick={() => onAddTask(status)}>
        + Nova tarefa
      </button>
    </section>
  )
}
