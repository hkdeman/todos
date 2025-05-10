package domain

import (
	"time"

	"github.com/google/uuid"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

type Todo struct {
	ID          uuid.UUID
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     *time.Time
	Priority    Priority
	Category    string
	Tags        []string
	Subtasks    []*Todo
	ParentID    *uuid.UUID
	AssignedTo  *uuid.UUID
	Comments    []Comment
	Recurring   *RecurringConfig
	Archived    bool
}

type Comment struct {
	ID        uuid.UUID
	Content   string
	CreatedAt time.Time
	UserID    uuid.UUID
}

type RecurringConfig struct {
	Frequency      string
	EndDate        *time.Time
	LastOccurrence time.Time
}

// NewTodo creates a new todo
func NewTodo(description string) *Todo {
	now := time.Now()
	return &Todo{
		ID:          uuid.New(),
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
		Priority:    PriorityMedium,
		Tags:        make([]string, 0),
		Subtasks:    make([]*Todo, 0),
		Comments:    make([]Comment, 0),
	}
}

// Update updates a todo
func (t *Todo) Update(completed bool, description string) {
	t.Completed = completed
	t.Description = description
	t.UpdatedAt = time.Now()
}

// AddSubtask adds a subtask to the todo
func (t *Todo) AddSubtask(subtask *Todo) {
	subtask.ParentID = &t.ID
	t.Subtasks = append(t.Subtasks, subtask)
	t.UpdatedAt = time.Now()
}

// AddComment adds a comment to the todo
func (t *Todo) AddComment(content string, userID uuid.UUID) {
	comment := Comment{
		ID:        uuid.New(),
		Content:   content,
		CreatedAt: time.Now(),
		UserID:    userID,
	}
	t.Comments = append(t.Comments, comment)
	t.UpdatedAt = time.Now()
}

// Archive marks the todo as archived
func (t *Todo) Archive() {
	t.Archived = true
	t.UpdatedAt = time.Now()
}

// SetRecurring sets the recurring configuration
func (t *Todo) SetRecurring(frequency string, endDate *time.Time) {
	t.Recurring = &RecurringConfig{
		Frequency:      frequency,
		EndDate:        endDate,
		LastOccurrence: time.Now(),
	}
	t.UpdatedAt = time.Now()
}
