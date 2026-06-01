package store

import (
	"context"
	"fmt"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) CreateNotification(ctx context.Context, cmd dto.CreateNotificationCmd) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO
			notifications (user_id, entitlement_id, type, scheduled_for)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT
			(user_id, type, entitlement_id)
		DO UPDATE SET
			scheduled_for = EXCLUDED.scheduled_for,
			sent_at = NULL
	`, cmd.UserID, cmd.EntitlementID, cmd.Type, cmd.ScheduledFor)
	if err != nil {
		return fmt.Errorf("s.pool.Exec: %w", err)
	}
	return nil
}

func (s *store) GetPendingNotifications(ctx context.Context) ([]*Notification, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			id, user_id, entitlement_id, type, scheduled_for, sent_at
		FROM
			notifications
		WHERE
			scheduled_for <= NOW()
		AND
			sent_at IS NULL
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return nil, fmt.Errorf("s.pool.Query: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.EntitlementID,
			&n.Type,
			&n.ScheduledFor,
			&n.SentAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		notifications = append(notifications, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return notifications, nil
}

func (s *store) MarkNotificationSent(ctx context.Context, cmd dto.MarkNotificationSentCmd) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE notifications
		SET sent_at = NOW()
		WHERE id = $1
	`, cmd.NotificationID)
	if err != nil {
		return fmt.Errorf("s.pool.Exec: %w", err)
	}
	return nil
}
