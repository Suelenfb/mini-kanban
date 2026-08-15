import React, { useEffect, useState, useCallback } from 'react'
import Board from './components/Board.jsx'
import TaskModal from './components/TaskModal.jsx'
import { api } from './api/api.js'

export const COLUMNS = [
  { status: 'todo', title: 'A Fazer' },
  { status: 'doing', title: 'Em Andamento' },
  { status: 'done', title: 'Concluído' },
]

export default function App() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [modalOpen, setModalOpen] = useState(false)
  const [editingTask, setEditingTask] = useState(null)
  const [defaultStatus, setDefaultStatus] = useState('todo')

  const loadTasks = useCallback(async () => {
    try {
      setLoading(true)
      const data = await api.listTasks()
      setTasks(data || [])
      setError(null)
    } catch (err) {
      setError('Não foi possível conectar ao backend. Verifique se ele está rodando em :8080.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadTasks()
  }, [loadTasks])

  function openCreateModal(status) {
    setEditingTask(null)
    setDefaultStatus(status)
    setModalOpen(true)
  }

  function openEditModal(task) {
    setEditingTask(task)
    setModalOpen(true)
  }

  async function handleSaveTask(formData) {
    if (editingTask) {
      const updated = await api.updateTask(editingTask.id, formData)
      setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
    } else {
      const created = await api.createTask({ ...formData, status: defaultStatus })
      setTasks((prev) => [...prev, created])
    }
    setModalOpen(false)
  }

  async function handleDeleteTask(id) {
    await api.deleteTask(id)
    setTasks((prev) => prev.filter((t) => t.id !== id))
  }

  // Chamado quando um card é solto em uma coluna (ou reordenado dentro dela).
  async function handleDropTask(taskId, targetStatus, targetIndex) {
    setTasks((prev) => {
      const moving = prev.find((t) => t.id === taskId)
      if (!moving) return prev

      const withoutMoving = prev.filter((t) => t.id !== taskId)
      const columnTasks = withoutMoving
        .filter((t) => t.status === targetStatus)
        .sort((a, b) => a.position - b.position)

      const updatedMoving = { ...moving, status: targetStatus }
      columnTasks.splice(targetIndex, 0, updatedMoving)

      const columnTasksWithPosition = columnTasks.map((t, idx) => ({ ...t, position: idx }))
      const otherTasks = withoutMoving.filter((t) => t.status !== targetStatus)

      // Persiste a nova ordem no backend (otimista: UI já foi atualizada acima).
      api
        .reorder(
          targetStatus,
          columnTasksWithPosition.map((t) => t.id),
        )
        .catch(() => loadTasks())

      return [...otherTasks, ...columnTasksWithPosition]
    })
  }

  return (
    <div className="app">
      <header className="app-header">
        <div>
          <h1>Mini Kanban</h1>
          <p className="subtitle">Organize suas tarefas arrastando entre as colunas</p>
        </div>
      </header>

      {error && <div className="error-banner">{error}</div>}

      {loading ? (
        <div className="loading">Carregando tarefas...</div>
      ) : (
        <Board
          tasks={tasks}
          onAddTask={openCreateModal}
          onEditTask={openEditModal}
          onDeleteTask={handleDeleteTask}
          onDropTask={handleDropTask}
        />
      )}

      {modalOpen && (
        <TaskModal
          task={editingTask}
          onClose={() => setModalOpen(false)}
          onSave={handleSaveTask}
        />
      )}
    </div>
  )
}
