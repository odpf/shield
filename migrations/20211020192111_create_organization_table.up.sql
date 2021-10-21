CREATE TABLE IF NOT EXISTS organizations
(
    id         uuid        PRIMARY KEY     DEFAULT uuid_generate_v4(),
    name       varchar     NOT NULL,
    slug       varchar     NOT NULL,
    metadata   jsonb,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);
