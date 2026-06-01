package carrierpolling

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	mobilecarrier "github.com/darkwingdck/adora-coding-assessment/internal/services/mobile_carrier"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

const (
	pollInterval   = 5 * time.Minute
	workerPoolSize = 10
)

type Worker struct {
	store   store.Store
	carrier mobilecarrier.Service
}

func NewWorker(store store.Store, carrier mobilecarrier.Service) *Worker {
	return &Worker{store: store, carrier: carrier}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.poll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) poll(ctx context.Context) {
	entitlements, err := w.store.GetCarrierEntitlements(ctx)
	if err != nil {
		log.Printf("carrier polling: failed to fetch carrier users: %v", err)
		return
	}

	sem := make(chan struct{}, workerPoolSize)
	var wg sync.WaitGroup

	for _, e := range entitlements {
		wg.Add(1)
		sem <- struct{}{}
		go func(e *store.Entitlement) {
			defer wg.Done()
			defer func() { <-sem }()
			w.processUser(ctx, e)
		}(e)
	}

	wg.Wait()
}

func (w *Worker) processUser(ctx context.Context, entitlement *store.Entitlement) {
	result, err := w.carrier.GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{
		UserID: entitlement.UserID,
	})
	if err != nil {
		log.Printf("carrier polling: w.carrier.GetMobileCarrierUserStatus, userID: %s: %v", entitlement.UserID, err)
		return
	}

	switch result.Status {
	case dto.MobileCarrierUserStatusApiError:
		log.Printf("carrier polling: api_error for userId %s", entitlement.UserID)

	case dto.MobileCarrierUserStatusActive:
		if entitlement.Active {
			return
		}
		reason := dto.EntitlementReasonCarrierActive
		err := w.store.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
			UserID:          entitlement.UserID,
			Active:          true,
			Source:          dto.EntitlementSourceCarrier,
			Reason:          &reason,
			ExpiresAt:       entitlement.ExpiresAt,
			LastEventTimeMs: entitlement.LastEventTimeMs,
		})
		if err != nil {
			log.Printf("carrier polling: w.store.UpdateEntitlement for userID %s: %v", entitlement.UserID, err)
		}

	case dto.MobileCarrierUserStatusInactive:
		if !entitlement.Active {
			return
		}
		reason := dto.EntitlementReasonCarrierInactive
		err := w.store.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
			UserID:          entitlement.UserID,
			Active:          false,
			Source:          dto.EntitlementSourceCarrier,
			Reason:          &reason,
			ExpiresAt:       nil,
			LastEventTimeMs: entitlement.LastEventTimeMs,
		})
		if err != nil {
			log.Printf("carrier polling: w.store.UpdateEntitlement for userID %s: %v", entitlement.UserID, err)
		}
	}
}
