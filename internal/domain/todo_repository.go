package domain

import (
	"time"

	"github.com/google/uuid"
)

type TodoRepository interface {
	Add(description string) *Todo
	Remove(id uuid.UUID)
	Update(id uuid.UUID, completed bool, description string) *Todo
	Search(search string) []*Todo
	All() []*Todo
	Get(id uuid.UUID) *Todo
	Reorder(ids []uuid.UUID) []*Todo

	// New methods for enhanced features
	GetByCategory(category string) []*Todo
	GetByTag(tag string) []*Todo
	GetByPriority(priority Priority) []*Todo
	GetByDueDate(start, end time.Time) []*Todo
	GetByAssignee(userID uuid.UUID) []*Todo
	GetRecurring() []*Todo
	GetArchived() []*Todo
	GetSubtasks(parentID uuid.UUID) []*Todo
	GetOverdue() []*Todo
	GetUpcoming(days int) []*Todo
}
