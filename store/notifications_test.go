package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_CreateNotification(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("creates notification", func(t *testing.T) {
		userID := "test_create_notif"
		cleanup(t, userID)

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID: userID,
			Source: dto.EntitlementSourceStore,
			Reason: ptr(dto.EntitlementReasonInitialPurchase),
		}))
		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)

		err = testStore.CreateNotification(ctx, dto.CreateNotificationCmd{
			UserID:        userID,
			EntitlementID: ent.ID,
			Type:          dto.NotificationTypePremiumExpiresSoon,
			ScheduledFor:  time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)

		var count int
		require.NoError(t, testPool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id = $1", userID).Scan(&count))
		require.Equal(t, 1, count)
	})

	t.Run("upserts on duplicate — does not create a second row", func(t *testing.T) {
		userID := "test_create_notif_dup"
		cleanup(t, userID)

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID: userID,
			Source: dto.EntitlementSourceStore,
			Reason: ptr(dto.EntitlementReasonInitialPurchase),
		}))
		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)

		cmd := dto.CreateNotificationCmd{
			UserID:        userID,
			EntitlementID: ent.ID,
			Type:          dto.NotificationTypePremiumExpiresSoon,
			ScheduledFor:  time.Now().Add(24 * time.Hour),
		}
		require.NoError(t, testStore.CreateNotification(ctx, cmd))
		require.NoError(t, testStore.CreateNotification(ctx, cmd))

		var count int
		require.NoError(t, testPool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id = $1", userID).Scan(&count))
		require.Equal(t, 1, count)
	})
}

func Test_GetPendingNotifications(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	pastUser := "test_pending_past"
	futureUser := "test_pending_future"
	cleanup(t, pastUser, futureUser)

	// past notification — should appear as pending
	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: pastUser,
		Source: dto.EntitlementSourceStore,
		Reason: ptr(dto.EntitlementReasonInitialPurchase),
	}))
	pastEnt, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: pastUser})
	require.NoError(t, err)
	require.NoError(t, testStore.CreateNotification(ctx, dto.CreateNotificationCmd{
		UserID:        pastUser,
		EntitlementID: pastEnt.ID,
		Type:          dto.NotificationTypePremiumExpiresSoon,
		ScheduledFor:  time.Now().Add(-1 * time.Hour),
	}))

	// future notification — should not appear
	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: futureUser,
		Source: dto.EntitlementSourceStore,
		Reason: ptr(dto.EntitlementReasonInitialPurchase),
	}))
	futureEnt, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: futureUser})
	require.NoError(t, err)
	require.NoError(t, testStore.CreateNotification(ctx, dto.CreateNotificationCmd{
		UserID:        futureUser,
		EntitlementID: futureEnt.ID,
		Type:          dto.NotificationTypePremiumExpiresSoon,
		ScheduledFor:  time.Now().Add(24 * time.Hour),
	}))

	notifications, err := testStore.GetPendingNotifications(ctx)
	require.NoError(t, err)

	var pastFound, futureFound bool
	for _, n := range notifications {
		if n.UserID == pastUser {
			pastFound = true
		}
		if n.UserID == futureUser {
			futureFound = true
		}
	}
	require.True(t, pastFound, "past notification should be pending")
	require.False(t, futureFound, "future notification should not be pending")
}

func Test_MarkNotificationSent(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	userID := "test_mark_sent"
	cleanup(t, userID)

	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: userID,
		Source: dto.EntitlementSourceStore,
		Reason: ptr(dto.EntitlementReasonInitialPurchase),
	}))
	ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
	require.NoError(t, err)

	require.NoError(t, testStore.CreateNotification(ctx, dto.CreateNotificationCmd{
		UserID:        userID,
		EntitlementID: ent.ID,
		Type:          dto.NotificationTypePremiumExpiresSoon,
		ScheduledFor:  time.Now().Add(-1 * time.Hour),
	}))

	var notifID uuid.UUID
	require.NoError(t, testPool.QueryRow(ctx, "SELECT id FROM notifications WHERE user_id = $1", userID).Scan(&notifID))

	require.NoError(t, testStore.MarkNotificationSent(ctx, dto.MarkNotificationSentCmd{NotificationID: notifID}))

	var sentAt *time.Time
	require.NoError(t, testPool.QueryRow(ctx, "SELECT sent_at FROM notifications WHERE id = $1", notifID).Scan(&sentAt))
	require.NotNil(t, sentAt)
}
