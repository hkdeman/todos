package todos

import (
	"context"

	"github.com/google/uuid"
	"github.com/stackus/todos/internal/domain"
)

// NotificationService handles todo reminders and notifications
type NotificationService interface {
	ScheduleReminder(ctx context.Context, todo *domain.Todo)
	SendNotification(ctx context.Context, userID uuid.UUID, message string)
}

// noopNotificationService is a no-operation implementation of NotificationService
type noopNotificationService struct{}

func NewNoopNotificationService() NotificationService {
	return &noopNotificationService{}
}

func (s *noopNotificationService) ScheduleReminder(ctx context.Context, todo *domain.Todo) {
	// No-op implementation
}

func (s *noopNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, message string) {
	// No-op implementation
}
