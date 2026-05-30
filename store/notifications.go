package store

import (
	"context"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) CreateNotification(ctx context.Context, cmd dto.CreateNotificationCmd) error {
	return nil
}

func (s *store) GetPendingNotifications(ctx context.Context) ([]*Notification, error) {
	return nil, nil
}

func (s *store) MarkNotificationSent(ctx context.Context, cmd dto.MarkNotificationSentCmd) error {
	return nil
}
