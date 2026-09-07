const API_BASE = 'http://localhost:8080/api'

async function request(path, options) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `request failed: ${res.status}`)
  }

  if (res.status === 204) return null
  return res.json()
}

export function listTasks() {
  return request('/tasks')
}

export function createTask({ title, description }) {
  return request('/tasks', {
    method: 'POST',
    body: JSON.stringify({ title, description }),
  })
}

export function updateTask(id, { title, description, done }) {
  return request(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ title, description, done }),
  })
}

export function deleteTask(id) {
  return request(`/tasks/${id}`, { method: 'DELETE' })
}
