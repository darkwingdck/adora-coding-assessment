# adora-coding-assessment

こんにちは!

Here you will find my solution to the coding assignment. It turned out to be quite large and interesting. Below I will describe my thoughts about it, the main architecture, and a few jokes!

## Main commands

Run project

```bash
cp .env.example .env && make run
```

Test project

```bash
make test
```

Generate mocks

```bash
make mocks
```

Generate Swagger documentation

```bash
make swag
```

## API examples

All endpoints can be tested through Swagger at:

http://localhost:8080/swagger

Ingest store webhook event

```bash
curl -X POST http://localhost:8080/webhooks/store \
  -H "Content-Type: application/json" \
  -d '{
    "eventId":    "evt_abc123",
    "userId":     "u_42",
    "type":       "INITIAL_PURCHASE",
    "eventTimeMs": 1716700000000,
    "productId":  "premium_monthly"
  }'
```

Marketplace bulk revoke

```bash
curl -X POST http://localhost:8080/webhooks/marketplace/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "userIds": ["u_42", "u_91", "u_133"]
  }'
```

Get current entitlement

```bash
curl http://localhost:8080/users/u_42/entitlement
```

Example response:

```json
{
  "active": true,
  "source": "STORE",
  "expiresAt": "2026-07-03T12:00:00Z",
  "lastChangedAt": "2026-06-03T10:23:00Z",
  "reason": "RENEWAL"
}
```

Seed test entitlements

Creates 30 test entitlements (10 per source) so you can test the marketplace and mobile carrier functional.

```bash
curl -X POST http://localhost:8080/test/seed
```

Example response:

```json
{
  "store":       ["user_store_test_1", "user_store_test_2", "..."],
  "carrier":     ["user_carrier_test_1", "user_carrier_test_2", "..."],
  "marketplace": ["user_marketplace_test_1", "user_marketplace_test_2", "..."]
}
```

After seeding you can immediately test marketplace revoke:

```bash
curl -X POST http://localhost:8080/webhooks/marketplace/revoke \
  -H "Content-Type: application/json" \
  -d '{"userIds": ["user_marketplace_test_1", "user_marketplace_test_2"]}'
```

## Architecture

The service is built using a three-layer architecture:

```text
api → service → store
```

1. The HTTP layer (`internal/api`) only parses requests and calls services.
2. Services (`internal/services`) contain all business logic: processing store events, deduplication through the `store_events` table, row-level locks (`FOR UPDATE`) to protect against race conditions during parallel webhook processing.
3. `store/` is a wrapper around PostgreSQL.

In addition to the HTTP server, two background workers are started in separate goroutines:

1. `carrier_polling` runs every 5 minutes, polls `service/mobile_carrier` for all users with `source = CARRIER`, and updates their status. Multiple worker instances can run simultaneously — this is safe because each user is processed independently.
2. `notification_worker` runs every minute and fetches due notifications using `FOR UPDATE SKIP LOCKED`. This mechanism guarantees that two worker instances cannot process the same notification twice.

## Trade-offs and main decisions

1. PostgreSQL vs SQLite. I chose PostgreSQL because of its richer feature set, locks, and enums.
2. Users table. Out of habit, I initially added a `users` table, but quickly realized that a subscription microservice should not manage users. If a message arrives with `user_id="u_42"`, it means that user already exists in another service and another database.
3. `store_events` — a table for tracking events from the in-app store. It helps track already processed events (thanks to the UNIQUE constraint on `event_id`) as well as outdated events (thanks to `event_time_ms`).
4. Goroutines. Initially, I wanted to implement notifications and carrier polling as endpoints and run them via cron jobs. Later I decided that, for the purposes of this assignment, two background workers implemented with goroutines would be sufficient.
5. `api -> service -> store` architecture. I try to keep all business logic inside services. API handlers simply call services, and services call the store layer (which knows nothing about business logic). For the same reason, I do not use database triggers — the less the database knows about business rules, the better.
6. Mock carrier endpoint. It is implemented simply as a method in a separate service. I could have exposed it as an actual endpoint, but decided not to. Instead, I accounted for the fact that in a real system the worker would make an HTTP request, which is why I did not wrap everything into a single transaction.

## Use of AI

1. Test generation. I approached this the same way I usually do: after implementing an atomic service, I write down all test cases that come to mind. Then I write two tests — a happy path and a test covering an error scenario. After that, I ask an AI assistant to generate the remaining test cases based on those examples, and then I review everything it generated.
2. Swagger annotations.
3. AI code completion in the editor.

## What I would do differently

Almost everything!

One of my team leads once said that a developer stops growing when they look at code they wrote a year ago and can no longer see any room for improvement.

Fortunately, I still see improvements in almost every part of this service.

The architecture should have been designed a little differently. Ideally, I would keep two tables:

1. An `entitlements` table — the source of truth from which we obtain the current subscription state of a user.
2. An `events` table — where all incoming events are stored. After each event, the entitlement for a specific user would be recalculated, taking `eventTimeMs` and statuses into account (an event state machine would fit nicely here).

A few more things I would add or do differently:

1. Component tests.
2. Validation in `internal/api`.
3. A separate PostgreSQL instance for tests (currently the production database container is effectively used for tests — don't try this at home!).
4. Golang Migrate — a convenient library for migration management and versioning.
5. Stretch work.
6. Grant method for Marketplace.

Thank you for the assignment, and I hope we will have a chance to discuss it further.

Talk to you soon!
