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
