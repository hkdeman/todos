package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Todos is a list of Todo
type Todos []*Todo

// NewTodos creates a new list of todos
func NewTodos() *Todos {
	return &Todos{}
}

// Add adds a todo to the list
func (l *Todos) Add(description string) *Todo {
	todo := NewTodo(description)
	*l = append(*l, todo)
	return todo
}

// Remove removes a todo from the list
func (l *Todos) Remove(id uuid.UUID) {
	index := l.indexOf(id)
	if index == -1 {
		return
	}
	*l = append((*l)[:index], (*l)[index+1:]...)
}

// Update updates a todo in the list
func (l *Todos) Update(id uuid.UUID, completed bool, description string) *Todo {
	index := l.indexOf(id)
	if index == -1 {
		return nil
	}
	todo := (*l)[index]
	todo.Update(completed, description)
	return todo
}

// Search returns a list of todos that match the search string
func (l *Todos) Search(search string) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if strings.Contains(todo.Description, search) {
			list = append(list, todo)
		}
	}
	return list
}

// All returns a copy of the list of todos
func (l *Todos) All() []*Todo {
	list := make([]*Todo, len(*l))
	copy(list, *l)
	return list
}

// Get returns a todo by id
func (l *Todos) Get(id uuid.UUID) *Todo {
	index := l.indexOf(id)
	if index == -1 {
		return nil
	}
	return (*l)[index]
}

// Reorder reorders the list of todos
func (l *Todos) Reorder(ids []uuid.UUID) []*Todo {
	newTodos := make([]*Todo, len(ids))
	for i, id := range ids {
		newTodos[i] = (*l)[l.indexOf(id)]
	}
	copy(*l, newTodos)
	return newTodos
}

// GetByCategory returns todos in the specified category
func (l *Todos) GetByCategory(category string) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.Category == category {
			list = append(list, todo)
		}
	}
	return list
}

// GetByTag returns todos with the specified tag
func (l *Todos) GetByTag(tag string) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		for _, t := range todo.Tags {
			if t == tag {
				list = append(list, todo)
				break
			}
		}
	}
	return list
}

// GetByPriority returns todos with the specified priority
func (l *Todos) GetByPriority(priority Priority) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.Priority == priority {
			list = append(list, todo)
		}
	}
	return list
}

// GetByDueDate returns todos due between start and end dates
func (l *Todos) GetByDueDate(start, end time.Time) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.DueDate != nil && !todo.DueDate.Before(start) && !todo.DueDate.After(end) {
			list = append(list, todo)
		}
	}
	return list
}

// GetByAssignee returns todos assigned to the specified user
func (l *Todos) GetByAssignee(userID uuid.UUID) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.AssignedTo != nil && *todo.AssignedTo == userID {
			list = append(list, todo)
		}
	}
	return list
}

// GetRecurring returns all recurring todos
func (l *Todos) GetRecurring() []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.Recurring != nil {
			list = append(list, todo)
		}
	}
	return list
}

// GetArchived returns all archived todos
func (l *Todos) GetArchived() []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.Archived {
			list = append(list, todo)
		}
	}
	return list
}

// GetSubtasks returns all subtasks for a given parent todo
func (l *Todos) GetSubtasks(parentID uuid.UUID) []*Todo {
	list := make([]*Todo, 0)
	for _, todo := range *l {
		if todo.ParentID != nil && *todo.ParentID == parentID {
			list = append(list, todo)
		}
	}
	return list
}

// GetOverdue returns all overdue todos
func (l *Todos) GetOverdue() []*Todo {
	list := make([]*Todo, 0)
	now := time.Now()
	for _, todo := range *l {
		if todo.DueDate != nil && todo.DueDate.Before(now) && !todo.Completed {
			list = append(list, todo)
		}
	}
	return list
}

// GetUpcoming returns todos due in the next specified number of days
func (l *Todos) GetUpcoming(days int) []*Todo {
	list := make([]*Todo, 0)
	now := time.Now()
	end := now.AddDate(0, 0, days)
	for _, todo := range *l {
		if todo.DueDate != nil && !todo.DueDate.Before(now) && !todo.DueDate.After(end) {
			list = append(list, todo)
		}
	}
	return list
}

// indexOf returns the index of the todo with the given id or -1 if not found
func (l *Todos) indexOf(id uuid.UUID) int {
	for i, todo := range *l {
		if todo.ID == id {
			return i
		}
	}
	return -1
}
