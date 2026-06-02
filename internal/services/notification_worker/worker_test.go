package notificationworker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// withTx makes WithTransaction actually invoke the callback so inner store
// calls can be verified on the same mock.
func withTx(mockStore *mocks.MockStore) {
	mockStore.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(store.Store) error) error {
			return fn(mockStore)
		})
}

func Test_process(t *testing.T) {
	ctx := context.Background()

	t.Run("no pending notifications -> no MarkNotificationSent calls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		mockStore.EXPECT().GetPendingNotifications(ctx).Return([]*store.Notification{}, nil)
		// MarkNotificationSent must NOT be called

		NewWorker(mockStore).process(ctx)
	})

	t.Run("pending notifications -> marks each as sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		notif1 := &store.Notification{ID: uuid.New(), UserID: "u_1", ScheduledFor: time.Now()}
		notif2 := &store.Notification{ID: uuid.New(), UserID: "u_2", ScheduledFor: time.Now()}

		withTx(mockStore)
		mockStore.EXPECT().GetPendingNotifications(ctx).Return([]*store.Notification{notif1, notif2}, nil)
		mockStore.EXPECT().MarkNotificationSent(ctx, dto.MarkNotificationSentCmd{NotificationID: notif1.ID}).Return(nil)
		mockStore.EXPECT().MarkNotificationSent(ctx, dto.MarkNotificationSentCmd{NotificationID: notif2.ID}).Return(nil)

		NewWorker(mockStore).process(ctx)
	})

	t.Run("GetPendingNotifications error -> no MarkNotificationSent calls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		mockStore.EXPECT().GetPendingNotifications(ctx).Return(nil, errors.New("db error"))
		// MarkNotificationSent must NOT be called

		NewWorker(mockStore).process(ctx)
	})

	t.Run("MarkNotificationSent error -> stops processing remaining notifications", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		notif1 := &store.Notification{ID: uuid.New(), UserID: "u_1", ScheduledFor: time.Now()}
		notif2 := &store.Notification{ID: uuid.New(), UserID: "u_2", ScheduledFor: time.Now()}

		withTx(mockStore)
		mockStore.EXPECT().GetPendingNotifications(ctx).Return([]*store.Notification{notif1, notif2}, nil)
		mockStore.EXPECT().MarkNotificationSent(ctx, dto.MarkNotificationSentCmd{NotificationID: notif1.ID}).Return(errors.New("db error"))
		// second MarkNotificationSent must NOT be called — the loop returns on first error

		NewWorker(mockStore).process(ctx)
	})
}
