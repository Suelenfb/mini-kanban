package models

import "time"

// Status representa a coluna do Kanban em que a tarefa se encontra.
type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

// IsValid verifica se o status recebido é um dos status suportados pelo Kanban.
func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusDoing, StatusDone:
		return true
	}
	return false
}

// Task representa uma tarefa do Kanban.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TaskInput é o payload aceito na criação/edição de uma tarefa.
type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

// ReorderInput é o payload usado para mover/reordenar tarefas entre colunas.
type ReorderInput struct {
	Status  Status   `json:"status"`
	TaskIDs []string `json:"taskIds"`
}
