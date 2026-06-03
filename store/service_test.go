package store_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	testPool  *pgxpool.Pool
	testStore store.Store
	dbErr     error
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	get := func(key, def string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return def
	}
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		get("DB_HOST", "localhost"),
		get("DB_PORT", "5432"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		dbErr = fmt.Errorf("connect to test DB: %w", err)
	} else if err := pool.Ping(ctx); err != nil {
		dbErr = fmt.Errorf("ping test DB: %w", err)
	} else {
		testPool = pool
		testStore = store.NewStore(pool)
	}

	code := m.Run()
	if pool != nil {
		pool.Close()
	}
	os.Exit(code)
}

func requireDB(t *testing.T) {
	t.Helper()
	if dbErr != nil {
		t.Skipf("skipping DB test: %v", dbErr)
	}
}

// cleanup registers a t.Cleanup that deletes all rows for the given userIDs.
// Call it at the start of each test so data is removed even on failure.
func cleanup(t *testing.T, userIDs ...string) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		for _, uid := range userIDs {
			testPool.Exec(ctx, "DELETE FROM notifications WHERE user_id = $1", uid)
			testPool.Exec(ctx, "DELETE FROM entitlements WHERE user_id = $1", uid)
			testPool.Exec(ctx, "DELETE FROM store_events WHERE user_id = $1", uid)
		}
	})
}

func ptr[T any](v T) *T { return &v }

func Test_WithTransaction(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("commits on success", func(t *testing.T) {
		userID := "test_tx_commit"
		cleanup(t, userID)

		err := testStore.WithTransaction(ctx, func(tx store.Store) error {
			return tx.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
				Active: true,
				UserID: userID,
				Source: dto.EntitlementSourceStore,
				Reason: ptr(dto.EntitlementReasonInitialPurchase),
			})
		})
		require.NoError(t, err)

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
	})

	t.Run("rolls back on error", func(t *testing.T) {
		userID := "test_tx_rollback"
		cleanup(t, userID)

		err := testStore.WithTransaction(ctx, func(tx store.Store) error {
			if err := tx.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
				Active: true,
				UserID: userID,
				Source: dto.EntitlementSourceStore,
				Reason: ptr(dto.EntitlementReasonInitialPurchase),
			}); err != nil {
				return err
			}
			return errors.New("force rollback")
		})
		require.Error(t, err)

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.Nil(t, ent)
	})
}
