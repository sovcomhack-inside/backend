CREATE UNLOGGED TABLE IF NOT EXISTS users (
  "id" bigserial PRIMARY KEY,
  "email" text UNIQUE,
  "image" text,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "password_hash" text NOT NULL,
  "password_salt" text NOT NULL,
  "created_at" timestamptz not null default now(),
  "updated_at" timestamptz not null default now()
);

CREATE TYPE user_status AS ENUM (
    'pending_approve',
    'approved',
    'banned'
);

CREATE UNLOGGED TABLE IF NOT EXISTS users_statuses (
    "id" bigserial PRIMARY KEY REFERENCES users(id),
    "status" user_status DEFAULT 'approved' NOT NULL
);

CREATE UNLOGGED TABLE IF NOT EXISTS accounts (
     number uuid not null primary key default gen_random_uuid(),
     user_id bigint,
     created_at timestamptz default now(),
     currency varchar(3) not null,
     balance bigint not null default 0
);

CREATE UNIQUE INDEX CONCURRENTLY account_user_id_currency
    ON accounts (user_id, currency);

create type operation_type as enum (
    'refill',
    'withdrawal',
    'transfer_incoming',
    'transfer_outgoing'
    );

CREATE UNLOGGED TABLE IF NOT EXISTS operations (
   id uuid                          not null primary key default gen_random_uuid(),
   purpose                          text,
   operation_type                   operation_type not null,
   time timestamptz                 default now(),
   receiver_account_number          uuid,
   receiver_account_amount_cents    bigint,
   receiver_account_currency        varchar(3),
   sender_account_number            uuid,
   sender_account_amount_cents      bigint,
   sender_account_currency          varchar(3),
   currencies_exchange_rate_ratio   float
);
