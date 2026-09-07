import { useState } from 'react'

export default function TaskForm({ onCreate }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e) {
    e.preventDefault()
    if (!title.trim()) return

    setSubmitting(true)
    try {
      await onCreate({ title, description })
      setTitle('')
      setDescription('')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="task-form" onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Task title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <input
        type="text"
        placeholder="Description (optional)"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <button type="submit" disabled={submitting || !title.trim()}>
        Add task
      </button>
    </form>
  )
}
