package store

import (
	"context"
	"fmt"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) CreateStoreEvent(ctx context.Context, cmd dto.CreateStoreEventCmd) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		INSERT
		INTO
			store_events (user_id, event_id, type, product_id, event_time_ms)
		VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT
			(event_id)
		DO NOTHING
	`, cmd.UserID, cmd.EventID, cmd.Type, cmd.ProductID, cmd.EventTimeMs)
	if err != nil {
		return false, fmt.Errorf("s.pool.Exec: %w", err)
	}

	return tag.RowsAffected() == 1, nil
}
