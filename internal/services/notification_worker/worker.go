package notificationworker

import (
	"context"
	"log"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

const pollInterval = time.Minute

type Worker struct {
	store store.Store
}

func NewWorker(store store.Store) *Worker {
	return &Worker{store: store}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.process(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	err := w.store.WithTransaction(ctx, func(tx store.Store) error {
		notifications, err := tx.GetPendingNotifications(ctx)
		if err != nil {
			return err
		}

		for _, n := range notifications {
			log.Printf("notification worker: sending %s for userId %s", n.Type, n.UserID)
			if err := tx.MarkNotificationSent(ctx, dto.MarkNotificationSentCmd{
				NotificationID: n.ID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("notification worker: process failed: %v", err)
	}
}
