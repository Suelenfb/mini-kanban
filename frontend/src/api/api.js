const API_URL = 'http://localhost:8080/api'

async function request(path, options = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `Erro na requisição (${res.status})`)
  }

  if (res.status === 204) return null
  return res.json()
}

export const api = {
  listTasks: () => request('/tasks'),

  createTask: (task) =>
    request('/tasks', {
      method: 'POST',
      body: JSON.stringify(task),
    }),

  updateTask: (id, task) =>
    request(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(task),
    }),

  deleteTask: (id) =>
    request(`/tasks/${id}`, {
      method: 'DELETE',
    }),

  reorder: (status, taskIds) =>
    request('/tasks/reorder', {
      method: 'PUT',
      body: JSON.stringify({ status, taskIds }),
    }),
}
