package store

import (
	"context"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) CreateStoreEvent(ctx context.Context, cmd dto.CreateStoreEventCmd) (bool, error) {
	return false, nil
}
