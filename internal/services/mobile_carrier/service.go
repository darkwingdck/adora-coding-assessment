package mobilecarrier

import (
	"context"
	"math/rand/v2"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

type Service interface {
	GetMobileCarrierUserStatus(
		ctx context.Context,
		cmd dto.GetMobileCarrierUserStatusCmd,
	) (*dto.GetMobileCarrierUserStatusResult, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GetMobileCarrierUserStatus(
	ctx context.Context,
	cmd dto.GetMobileCarrierUserStatusCmd,
) (*dto.GetMobileCarrierUserStatusResult, error) {
	randomNumber := rand.IntN(100)

	if randomNumber >= 95 {
		return &dto.GetMobileCarrierUserStatusResult{
			Status: dto.MobileCarrierUserStatusApiError,
		}, nil
	}

	if randomNumber >= 85 {
		return &dto.GetMobileCarrierUserStatusResult{
			Status: dto.MobileCarrierUserStatusInactive,
		}, nil
	}

	return &dto.GetMobileCarrierUserStatusResult{
		Status: dto.MobileCarrierUserStatusActive,
	}, nil
}
