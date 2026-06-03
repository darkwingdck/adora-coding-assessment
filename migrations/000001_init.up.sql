CREATE TYPE entitlement_source AS ENUM (
    'STORE',
    'CARRIER',
    'MARKETPLACE',
    'NONE'
);

CREATE TYPE event_type AS ENUM (
    'INITIAL_PURCHASE',
    'RENEWAL',
    'CANCELLATION',
    'BILLING_ISSUE',
    'EXPIRATION',
    'UN_CANCELLATION'
);

CREATE TYPE entitlement_reason AS ENUM (
    'INITIAL_PURCHASE',
    'RENEWAL',
    'CANCELLATION',
    'BILLING_ISSUE',
    'EXPIRATION',
    'UN_CANCELLATION',
    'CARRIER_INACTIVE',
    'MARKETPLACE_REVOKE'
);

CREATE TABLE entitlements (
    id                  UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             TEXT                NOT NULL UNIQUE,
    active              BOOLEAN             NOT NULL DEFAULT FALSE,
    source              entitlement_source  NOT NULL DEFAULT 'NONE',
    reason              entitlement_reason,
    expires_at          TIMESTAMPTZ,
    last_changed_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    last_event_time_ms  BIGINT              NOT NULL DEFAULT 0
);

CREATE TABLE store_events (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         TEXT        NOT NULL,
    event_id        TEXT        NOT NULL UNIQUE,
    type            event_type  NOT NULL,
    product_id      TEXT        NOT NULL,
    event_time_ms   BIGINT      NOT NULL,
    processed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         TEXT        NOT NULL,
    entitlement_id  UUID        NOT NULL REFERENCES entitlements(id),
    type            TEXT        NOT NULL,
    scheduled_for   TIMESTAMPTZ NOT NULL,
    sent_at         TIMESTAMPTZ,
    UNIQUE (user_id, type, entitlement_id)
);

CREATE INDEX idx_entitlements_source ON entitlements (source);

CREATE INDEX idx_notifications_pending ON notifications (scheduled_for) WHERE sent_at IS NULL;

CREATE INDEX idx_notifications_entitlement_id ON notifications (entitlement_id);
