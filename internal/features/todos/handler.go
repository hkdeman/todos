package todos

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/stackus/todos/internal/domain"
	"github.com/stackus/todos/internal/templates/pages"
	"github.com/stackus/todos/internal/templates/partials"
)

type (
	Handler interface {
		// Search : GET /todos
		Search(w http.ResponseWriter, r *http.Request)
		// Create : POST /todos
		Create(w http.ResponseWriter, r *http.Request)
		// Update : PATCH /todos/{todoId}
		// Update : POST /todos/{todoId}/edit
		Update(w http.ResponseWriter, r *http.Request)
		// Get : GET /todos/{todoId}
		Get(w http.ResponseWriter, r *http.Request)
		// Delete : DELETE /todos/{todoId}
		// Delete : POST /todos/{todoId}/delete
		Delete(w http.ResponseWriter, r *http.Request)
		// Sort : POST /todos/sort
		Sort(w http.ResponseWriter, r *http.Request)
		// CreateTodo : POST /todos/create
		CreateTodo(w http.ResponseWriter, r *http.Request)
		// AddSubtask : POST /todos/add-subtask
		AddSubtask(w http.ResponseWriter, r *http.Request)
		// AddComment : POST /todos/add-comment
		AddComment(w http.ResponseWriter, r *http.Request)
	}

	handler struct {
		service Service
	}

	// New request/response types
	CreateTodoRequest struct {
		Description string     `json:"description"`
		DueDate     *time.Time `json:"dueDate,omitempty"`
		Priority    int        `json:"priority"`
		Category    string     `json:"category,omitempty"`
		Tags        []string   `json:"tags,omitempty"`
	}

	CommentRequest struct {
		Content string `json:"content"`
		UserID  string `json:"userId"`
	}
)

func NewHandler(svc Service) Handler {
	return &handler{service: svc}
}

func Mount(r chi.Router, h Handler) {
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.Search)
		r.Post("/", h.Create)
		r.Route("/{todoId}", func(r chi.Router) {
			r.Patch("/", h.Update)
			r.Post("/edit", h.Update)
			r.Get("/", h.Get)
			r.Delete("/", h.Delete)
			r.Post("/delete", h.Delete)
		})
		r.Post("/sort", h.Sort)
		r.Post("/create", h.CreateTodo)
		r.Post("/add-subtask", h.AddSubtask)
		r.Post("/add-comment", h.AddComment)
	})
}

func (h handler) Sort(w http.ResponseWriter, r *http.Request) {
	var todoIDs []uuid.UUID
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, id := range r.Form["id"] {
		var todoID uuid.UUID
		var err error
		if todoID, err = uuid.Parse(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todoIDs = append(todoIDs, todoID)
	}
	if err := h.service.Sort(r.Context(), todoIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h handler) Search(w http.ResponseWriter, r *http.Request) {
	var search = r.URL.Query().Get("search")
	todos, err := h.service.Search(r.Context(), search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodos(todos).Render(r.Context(), w)
	default:
		err = pages.TodosPage(todos, search).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var description = r.Form.Get("description")

	todo, err := h.service.Add(r.Context(), description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodo(todo).Render(r.Context(), w)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var completed = r.Form.Get("completed") == "true"
	var description = r.Form.Get("description")

	todo, err := h.service.Update(r.Context(), todoID, completed, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodo(todo).Render(r.Context(), w)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := h.service.Get(r.Context(), todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.EditTodoForm(todo).Render(r.Context(), w)
	default:
		err = pages.TodoPage(todo).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Remove(r.Context(), todoID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		_, err = w.Write([]byte(""))
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.service.AddWithDetails(r.Context(), req.Description, req.DueDate,
		domain.Priority(req.Priority), req.Category, req.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func (h handler) AddSubtask(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parentId")
	parentUUID, err := uuid.Parse(parentID)
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subtask, err := h.service.AddSubtask(r.Context(), parentUUID, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(subtask)
}

func (h handler) AddComment(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("todoId")
	todoUUID, err := uuid.Parse(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.service.AddComment(r.Context(), todoUUID, req.Content, userUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isHTMX(r *http.Request) bool {
	// Check for "HX-Request" header
	if r.Header.Get("HX-Request") != "" {
		return true
	}

	return false
}
