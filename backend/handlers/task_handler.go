package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"kanban-backend/models"
	"kanban-backend/storage"
)

// TaskHandler agrupa as dependências necessárias para atender as
// requisições relacionadas a tarefas.
type TaskHandler struct {
	store *storage.TaskStore
}

func NewTaskHandler(store *storage.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// TasksCollection trata GET /api/tasks e POST /api/tasks
func (h *TaskHandler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}

// TaskItem trata GET/PUT/DELETE /api/tasks/{id}
func (h *TaskHandler) TaskItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id da tarefa não informado")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}

// Reorder trata PUT /api/tasks/reorder
func (h *TaskHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}

	var in models.ReorderInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if err := h.store.Reorder(in); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "uma ou mais tarefas não foram encontradas")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var in models.TaskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "o título da tarefa é obrigatório")
		return
	}

	task, err := h.store.Create(in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) get(w http.ResponseWriter, id string) {
	task, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var in models.TaskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	task, err := h.store.Update(id, in)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "tarefa não encontrada")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) delete(w http.ResponseWriter, id string) {
	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "tarefa não encontrada")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deletado"})
}
