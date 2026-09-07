export default function TaskItem({ task, onToggle, onDelete }) {
  return (
    <li className="task-item">
      <label>
        <input
          type="checkbox"
          checked={task.done}
          onChange={() => onToggle(task)}
        />
        <span className={task.done ? 'done' : ''}>{task.title}</span>
      </label>
      {task.description && <p className="task-description">{task.description}</p>}
      <button onClick={() => onDelete(task.id)}>Delete</button>
    </li>
  )
}
