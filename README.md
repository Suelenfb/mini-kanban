# Mini Kanban — React + Go

Aplicação de Kanban simples para gerenciamento de tarefas, com **CRUD completo** e **drag-and-drop fluido** entre as colunas `To Do`, `Doing` e `Done`.

- **Frontend:** React 18 + Vite (JavaScript puro, sem bibliotecas externas de drag-and-drop — implementado com a API nativa do HTML5)
- **Backend:** Go, usando apenas a biblioteca padrão (`net/http`), com dados em memória

---

## 📁 Estrutura do projeto

```
kanban-app/
├── backend/                 # API REST em Go
│   ├── main.go               # bootstrap do servidor, rotas e middleware CORS
│   ├── go.mod
│   ├── models/
│   │   └── task.go           # struct Task, Status e payloads de entrada
│   ├── storage/
│   │   └── storage.go        # camada de persistência em memória (thread-safe)
│   └── handlers/
│       └── task_handler.go   # handlers HTTP (list/create/get/update/delete/reorder)
│
├── frontend/                 # SPA em React
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── src/
│       ├── main.jsx
│       ├── App.jsx           # estado global das tarefas + orquestração das ações
│       ├── styles.css
│       ├── api/
│       │   └── api.js        # camada de comunicação com o backend (fetch)
│       └── components/
│           ├── Board.jsx     # monta as 3 colunas
│           ├── Column.jsx    # coluna: lista de cards + drop zone
│           ├── TaskCard.jsx  # card arrastável (draggable)
│           └── TaskModal.jsx # formulário de criação/edição
│
├── docs/
│   └── user-flow-diagram.svg # diagrama do fluxo de uso da aplicação
└── README.md
```

> No VS Code, abra a pasta `kanban-app/` como workspace: `backend/` e `frontend/` já ficam separados na árvore de arquivos, cada um com sua própria configuração, e você pode abrir dois terminais integrados (uma aba para cada) para rodar os dois serviços lado a lado.

---

## ▶️ Como rodar

### Pré-requisitos
- Go 1.22+
- Node.js 18+

### 1. Backend (porta `8080`)

```bash
cd backend
go run main.go
```

O servidor sobe em `http://localhost:8080` já com 3 tarefas de exemplo (uma em cada coluna) e libera CORS para o frontend.

### 2. Frontend (porta `5173`)

```bash
cd frontend
npm install
npm run dev
```

Acesse `http://localhost:5173`.

---

## 🔌 API REST

| Método | Rota                    | Descrição                                             |
|--------|--------------------------|--------------------------------------------------------|
| GET    | `/api/tasks`             | Lista todas as tarefas (ordenadas por coluna/posição)  |
| POST   | `/api/tasks`              | Cria uma tarefa `{ title, description, status? }`      |
| GET    | `/api/tasks/{id}`         | Busca uma tarefa                                       |
| PUT    | `/api/tasks/{id}`         | Atualiza título/descrição/status                        |
| DELETE | `/api/tasks/{id}`         | Remove uma tarefa                                       |
| PUT    | `/api/tasks/reorder`      | Reordena/move tarefas `{ status, taskIds: [...] }`      |

O status aceita apenas os valores `todo`, `doing` ou `done`.

---

## 🧠 Decisões de arquitetura

- **Backend em camadas** (`models` → `storage` → `handlers` → `main`): separa o domínio, a persistência e a exposição HTTP, facilitando trocar o storage em memória por um banco de dados real no futuro sem tocar nos handlers.
- **Storage thread-safe:** o `TaskStore` usa `sync.RWMutex` porque o `net/http` atende cada requisição em uma goroutine própria.
- **Posição explícita (`position`) por tarefa:** permite reordenar dentro da mesma coluna e mover entre colunas com uma única rota (`/reorder`), mantendo a ordem visual consistente com o backend.
- **Drag-and-drop nativo (sem dependências):** usa `draggable` + eventos `onDragStart/onDragOver/onDrop` do HTML5, evitando dependências externas.
- **Atualização otimista no frontend:** ao soltar um card, a UI já reordena localmente e envia a persistência em paralelo — a movimentação entre colunas fica instantânea; se a chamada falhar, o board é recarregado do backend.

---

## 🗺️ Diagrama de User Flow

Veja [`docs/user-flow-diagram.svg`](./docs/user-flow-diagram.svg).

Resumo do fluxo:

1. Usuário acessa o board → frontend busca as tarefas (`GET /api/tasks`) e renderiza as 3 colunas.
2. A partir daí, pode:
   - **Criar** tarefa → modal → `POST /api/tasks` → card aparece em `A Fazer`.
   - **Editar** tarefa → modal → `PUT /api/tasks/{id}` → card atualizado.
   - **Mover** tarefa (drag-and-drop) → `PUT /api/tasks/reorder` → card muda de coluna/posição.
   - **Excluir** tarefa → `DELETE /api/tasks/{id}` → card removido.
3. Após qualquer ação, o board permanece consistente e o usuário pode repetir o fluxo livremente.

---

## ✅ Escopo mínimo atendido

- [x] CRUD completo de tarefas (criar, ler, atualizar, excluir)
- [x] Fluidez entre colunas via drag-and-drop com persistência da ordem/posição
- [x] Código organizado em camadas (models/storage/handlers no backend; api/components no frontend)
- [x] UI simples, responsiva a interações (modal, hover, drag feedback, contadores por coluna)
- [x] README + diagrama de User Flow
