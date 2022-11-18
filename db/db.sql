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
  "id" tex REFERENCES users(id),
  "status" user_status DEFAULT 'pending_approve' NOT NULL
);
