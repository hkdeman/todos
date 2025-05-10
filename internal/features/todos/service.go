package todos

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stackus/todos/internal/domain"
)

type (
	Service interface {
		// Add adds a todo to the list
		Add(ctx context.Context, description string) (*domain.Todo, error)
		// Remove removes a todo from the list
		Remove(ctx context.Context, id uuid.UUID) error
		// Update updates a todo in the list
		Update(ctx context.Context, id uuid.UUID, completed bool, description string) (*domain.Todo, error)
		// Search returns a list of todos that match the search string
		Search(ctx context.Context, search string) ([]*domain.Todo, error)
		// Get returns a todo by id
		Get(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
		// Sort sorts the todos by the given ids
		Sort(ctx context.Context, ids []uuid.UUID) error

		// New methods for enhanced features
		AddWithDetails(ctx context.Context, description string, dueDate *time.Time, priority domain.Priority, category string, tags []string) (*domain.Todo, error)
		AddSubtask(ctx context.Context, parentID uuid.UUID, description string) (*domain.Todo, error)
		AddComment(ctx context.Context, todoID uuid.UUID, content string, userID uuid.UUID) error
		SetRecurring(ctx context.Context, id uuid.UUID, frequency string, endDate *time.Time) error
		Archive(ctx context.Context, id uuid.UUID) error
		Assign(ctx context.Context, todoID uuid.UUID, userID uuid.UUID) error

		// Query methods
		GetByCategory(ctx context.Context, category string) ([]*domain.Todo, error)
		GetByTag(ctx context.Context, tag string) ([]*domain.Todo, error)
		GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Todo, error)
		GetByDueDate(ctx context.Context, start, end time.Time) ([]*domain.Todo, error)
		GetByAssignee(ctx context.Context, userID uuid.UUID) ([]*domain.Todo, error)
		GetRecurring(ctx context.Context) ([]*domain.Todo, error)
		GetArchived(ctx context.Context) ([]*domain.Todo, error)
		GetSubtasks(ctx context.Context, parentID uuid.UUID) ([]*domain.Todo, error)
		GetOverdue(ctx context.Context) ([]*domain.Todo, error)
		GetUpcoming(ctx context.Context, days int) ([]*domain.Todo, error)
	}

	service struct {
		todos         domain.TodoRepository
		notifications NotificationService
	}
)

func NewService(todos domain.TodoRepository) Service {
	return &service{
		todos:         todos,
		notifications: NewNoopNotificationService(),
	}
}

func (s service) Add(_ context.Context, description string) (*domain.Todo, error) {
	todo := s.todos.Add(description)

	return todo, nil
}

func (s service) Remove(_ context.Context, id uuid.UUID) error {
	s.todos.Remove(id)

	return nil
}

func (s service) Update(_ context.Context, id uuid.UUID, completed bool, description string) (*domain.Todo, error) {
	todo := s.todos.Update(id, completed, description)

	return todo, nil
}

func (s service) Search(_ context.Context, search string) ([]*domain.Todo, error) {
	todos := s.todos.Search(search)

	return todos, nil
}

func (s service) Get(_ context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo := s.todos.Get(id)

	return todo, nil
}

func (s service) Sort(_ context.Context, ids []uuid.UUID) error {
	s.todos.Reorder(ids)

	return nil
}

func (s *service) AddWithDetails(ctx context.Context, description string, dueDate *time.Time, priority domain.Priority, category string, tags []string) (*domain.Todo, error) {
	if description == "" {
		return nil, ErrInvalidInput
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		return nil, ErrInvalidDate
	}

	if priority < domain.PriorityLow || priority > domain.PriorityHigh {
		return nil, ErrInvalidPriority
	}

	todo := s.todos.Add(description)
	todo.DueDate = dueDate
	todo.Priority = priority
	todo.Category = category
	todo.Tags = tags
	todo.UpdatedAt = time.Now()

	if dueDate != nil {
		s.notifications.ScheduleReminder(ctx, todo)
	}

	return todo, nil
}

func (s *service) AddSubtask(ctx context.Context, parentID uuid.UUID, description string) (*domain.Todo, error) {
	if description == "" {
		return nil, ErrInvalidInput
	}

	parent := s.todos.Get(parentID)
	if parent == nil {
		return nil, ErrTodoNotFound
	}

	if parent.Archived {
		return nil, ErrInvalidInput
	}

	subtask := s.todos.Add(description)
	parent.AddSubtask(subtask)
	return subtask, nil
}

func (s *service) AddComment(ctx context.Context, todoID uuid.UUID, content string, userID uuid.UUID) error {
	if content == "" {
		return ErrInvalidInput
	}

	todo := s.todos.Get(todoID)
	if todo == nil {
		return ErrTodoNotFound
	}

	if todo.Archived {
		return ErrInvalidInput
	}

	todo.AddComment(content, userID)
	return nil
}

func (s *service) SetRecurring(ctx context.Context, id uuid.UUID, frequency string, endDate *time.Time) error {
	todo := s.todos.Get(id)
	if todo == nil {
		return ErrTodoNotFound
	}

	todo.SetRecurring(frequency, endDate)
	return nil
}

func (s *service) Archive(ctx context.Context, id uuid.UUID) error {
	todo := s.todos.Get(id)
	if todo == nil {
		return ErrTodoNotFound
	}

	todo.Archive()
	return nil
}

func (s *service) Assign(ctx context.Context, todoID uuid.UUID, userID uuid.UUID) error {
	todo := s.todos.Get(todoID)
	if todo == nil {
		return ErrTodoNotFound
	}

	todo.AssignedTo = &userID
	todo.UpdatedAt = time.Now()
	return nil
}

func (s *service) GetByCategory(ctx context.Context, category string) ([]*domain.Todo, error) {
	return s.todos.GetByCategory(category), nil
}

func (s *service) GetByTag(ctx context.Context, tag string) ([]*domain.Todo, error) {
	return s.todos.GetByTag(tag), nil
}

func (s *service) GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Todo, error) {
	return s.todos.GetByPriority(priority), nil
}

func (s *service) GetByDueDate(ctx context.Context, start, end time.Time) ([]*domain.Todo, error) {
	return s.todos.GetByDueDate(start, end), nil
}

func (s *service) GetByAssignee(ctx context.Context, userID uuid.UUID) ([]*domain.Todo, error) {
	return s.todos.GetByAssignee(userID), nil
}

func (s *service) GetRecurring(ctx context.Context) ([]*domain.Todo, error) {
	return s.todos.GetRecurring(), nil
}

func (s *service) GetArchived(ctx context.Context) ([]*domain.Todo, error) {
	return s.todos.GetArchived(), nil
}

func (s *service) GetSubtasks(ctx context.Context, parentID uuid.UUID) ([]*domain.Todo, error) {
	return s.todos.GetSubtasks(parentID), nil
}

func (s *service) GetOverdue(ctx context.Context) ([]*domain.Todo, error) {
	return s.todos.GetOverdue(), nil
}

func (s *service) GetUpcoming(ctx context.Context, days int) ([]*domain.Todo, error) {
	return s.todos.GetUpcoming(days), nil
}
