import React from 'react'
import Column from './Column.jsx'
import { COLUMNS } from '../App.jsx'

export default function Board({ tasks, onAddTask, onEditTask, onDeleteTask, onDropTask }) {
  return (
    <div className="board">
      {COLUMNS.map((col) => {
        const columnTasks = tasks
          .filter((t) => t.status === col.status)
          .sort((a, b) => a.position - b.position)

        return (
          <Column
            key={col.status}
            status={col.status}
            title={col.title}
            tasks={columnTasks}
            onAddTask={onAddTask}
            onEditTask={onEditTask}
            onDeleteTask={onDeleteTask}
            onDropTask={onDropTask}
          />
        )
      })}
    </div>
  )
}
