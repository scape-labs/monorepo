-- TODO: initial schema for flow.
-- Use platform/migrator conventions: forward-only, no destructive changes.

CREATE TABLE IF NOT EXISTS flow_example (
    id           BIGINT PRIMARY KEY,
    tenant_id    TEXT      NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS flow_example_tenant_idx
    ON flow_example (tenant_id);
