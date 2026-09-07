import { useEffect, useState } from 'react'
import { listTasks, createTask, updateTask, deleteTask } from './api/tasks.js'
import TaskForm from './components/TaskForm.jsx'
import TaskList from './components/TaskList.jsx'

export default function App() {
  const [tasks, setTasks] = useState([])
  const [error, setError] = useState(null)

  useEffect(() => {
    listTasks().then(setTasks).catch((e) => setError(e.message))
  }, [])

  async function handleCreate({ title, description }) {
    const task = await createTask({ title, description })
    setTasks((prev) => [task, ...prev])
  }

  async function handleToggle(task) {
    const updated = await updateTask(task.id, { ...task, done: !task.done })
    setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
  }

  async function handleDelete(id) {
    await deleteTask(id)
    setTasks((prev) => prev.filter((t) => t.id !== id))
  }

  return (
    <main className="app">
      <h1>Habit / Task Tracker</h1>
      {error && <p className="error">{error}</p>}
      <TaskForm onCreate={handleCreate} />
      <TaskList tasks={tasks} onToggle={handleToggle} onDelete={handleDelete} />
    </main>
  )
}
