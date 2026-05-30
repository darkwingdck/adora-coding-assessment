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

CREATE TABLE users (
    id          TEXT        PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE entitlements (
    id                  UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             TEXT                NOT NULL UNIQUE REFERENCES users(id),
    active              BOOLEAN             NOT NULL DEFAULT FALSE,
    source              entitlement_source  NOT NULL DEFAULT 'NONE',
    reason              entitlement_reason,
    expires_at          TIMESTAMPTZ,
    last_changed_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    last_event_time_ms  BIGINT              NOT NULL DEFAULT 0
);

CREATE TABLE store_events (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         TEXT        NOT NULL REFERENCES users(id),
    event_id        TEXT        NOT NULL UNIQUE,
    type            event_type  NOT NULL,
    product_id      TEXT        NOT NULL,
    event_time_ms   BIGINT      NOT NULL,
    processed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         TEXT        NOT NULL REFERENCES users(id),
    entitlement_id  UUID        NOT NULL REFERENCES entitlements(id),
    type            TEXT        NOT NULL,
    scheduled_for   TIMESTAMPTZ NOT NULL,
    sent_at         TIMESTAMPTZ,
    UNIQUE (user_id, type, entitlement_id)
);
