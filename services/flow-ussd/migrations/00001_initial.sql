-- TODO: initial schema for flow-ussd.
-- Use platform/migrator conventions: forward-only, no destructive changes.

CREATE TABLE IF NOT EXISTS flow-ussd_example (
    id           BIGINT PRIMARY KEY,
    tenant_id    TEXT      NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS flow-ussd_example_tenant_idx
    ON flow-ussd_example (tenant_id);
