package store

import (
	"context"
	"fmt"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) UpsertUser(ctx context.Context, cmd dto.UpsertUserCmd) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO users (id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
	`, cmd.UserID)
	if err != nil {
		return fmt.Errorf("s.pool.Exec: %w", err)
	}
	return nil
}
