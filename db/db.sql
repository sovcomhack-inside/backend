CREATE UNLOGGED TABLE IF NOT EXISTS users (
  "id" bigserial PRIMARY KEY,
  "email" text NOT NULL UNIQUE,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "password_hash" text NOT NULL,
  "password_salt" text NOT NULL,
  "created_at" timestamptz not null default now(),
  "updated_at" timestamptz not null default now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS accounts (
    number uuid,
    user_id bigint references users,
    created_at timestamptz,
    currency varchar(3),
    cents bigint
);

CREATE UNIQUE INDEX CONCURRENTLY account_user_id_currency
ON accounts (user_id, currency);

create type operation_type as enum (
    'refill',
    'withdrawal',
    'transfer_incoming',
    'transfer_outgoing',
    'purchase',
    'sale'
);

CREATE UNLOGGED TABLE IF NOT EXISTS operations (
    id uuid not null primary key,
    operation_type operation_type not null,
    purpose text,
    time timestamptz,
    account_number varchar(20),
    account_amount_cents bigint,
    account_amount_currency varchar(3)
);
