package storage

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"kanban-backend/models"
)

// ErrNotFound é retornado quando uma tarefa não é encontrada no storage.
var ErrNotFound = errors.New("tarefa não encontrada")

// TaskStore mantém as tarefas em memória, protegidas por um mutex para
// permitir acesso concorrente seguro (o servidor HTTP atende requisições
// em goroutines diferentes).
type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[string]*models.Task
	nextID int
}

// NewTaskStore cria um storage já populado com algumas tarefas de exemplo,
// facilitando a avaliação da aplicação assim que ela sobe.
func NewTaskStore() *TaskStore {
	s := &TaskStore{
		tasks:  make(map[string]*models.Task),
		nextID: 1,
	}

	seed := []models.TaskInput{
		{Title: "Configurar repositório", Description: "Criar repositório e estrutura inicial do projeto", Status: models.StatusDone},
		{Title: "Modelar API REST", Description: "Definir endpoints de CRUD das tarefas", Status: models.StatusDoing},
		{Title: "Criar board no frontend", Description: "Montar colunas To Do / Doing / Done em React", Status: models.StatusTodo},
	}
	for _, in := range seed {
		_, _ = s.Create(in)
	}

	return s
}

func (s *TaskStore) newID() string {
	id := strconv.Itoa(s.nextID)
	s.nextID++
	return id
}

// List retorna todas as tarefas, ordenadas por status e depois por posição.
func (s *TaskStore) List() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Status != result[j].Status {
			return result[i].Status < result[j].Status
		}
		return result[i].Position < result[j].Position
	})

	return result
}

// Get retorna uma tarefa pelo ID.
func (s *TaskStore) Get(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// Create insere uma nova tarefa, calculando automaticamente sua posição
// (sempre no final da coluna de destino).
func (s *TaskStore) Create(in models.TaskInput) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status := in.Status
	if status == "" {
		status = models.StatusTodo
	}
	if !status.IsValid() {
		return nil, errors.New("status inválido")
	}

	position := s.nextPositionLocked(status)
	now := time.Now()

	task := &models.Task{
		ID:          s.newID(),
		Title:       in.Title,
		Description: in.Description,
		Status:      status,
		Position:    position,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task

	return task, nil
}

// nextPositionLocked calcula a próxima posição livre em uma coluna.
// Deve ser chamada com o mutex já travado.
func (s *TaskStore) nextPositionLocked(status models.Status) int {
	max := -1
	for _, t := range s.tasks {
		if t.Status == status && t.Position > max {
			max = t.Position
		}
	}
	return max + 1
}

// Update atualiza título, descrição e/ou status de uma tarefa existente.
// Ao mudar de status sem informar reordenação explícita, a tarefa vai
// para o final da nova coluna.
func (s *TaskStore) Update(id string, in models.TaskInput) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}

	if in.Title != "" {
		task.Title = in.Title
	}
	task.Description = in.Description

	if in.Status != "" && in.Status != task.Status {
		if !in.Status.IsValid() {
			return nil, errors.New("status inválido")
		}
		task.Status = in.Status
		task.Position = s.nextPositionLocked(in.Status)
	}

	task.UpdatedAt = time.Now()
	return task, nil
}

// Delete remove uma tarefa pelo ID.
func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}

// Reorder define a ordem (e, se necessário, o status) de um conjunto de
// tarefas de uma coluna, usado após um drag-and-drop no frontend.
func (s *TaskStore) Reorder(in models.ReorderInput) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !in.Status.IsValid() {
		return errors.New("status inválido")
	}

	for i, id := range in.TaskIDs {
		task, ok := s.tasks[id]
		if !ok {
			return ErrNotFound
		}
		task.Status = in.Status
		task.Position = i
		task.UpdatedAt = time.Now()
	}

	return nil
}
