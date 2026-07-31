-- +goose Up
create type user_role as enum ('administrator', 'editor');
create table admin_users (
  id bigserial primary key,
  login text not null unique,
  password_hash text not null,
  role user_role not null default 'editor',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create table admin_sessions (
  id uuid primary key,
  user_id bigint not null references admin_users(id) on delete cascade,
  token_hash char(64) not null unique,
  expires_at timestamptz not null,
  revoked_at timestamptz,
  created_at timestamptz not null default now()
);
create index admin_sessions_active on admin_sessions (token_hash, expires_at) where revoked_at is null;

alter table classes add column publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  add column publish_at timestamptz;
alter table news add column publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  add column publish_at timestamptz;
alter table services add column publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  add column publish_at timestamptz;
alter table staff add column publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  add column publish_at timestamptz;
alter table documents add column publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  add column publish_at timestamptz;

-- +goose Down
alter table documents drop column publish_at, drop column publication_status;
alter table staff drop column publish_at, drop column publication_status;
alter table services drop column publish_at, drop column publication_status;
alter table news drop column publish_at, drop column publication_status;
alter table classes drop column publish_at, drop column publication_status;
drop table if exists admin_sessions;
drop table if exists admin_users;
drop type if exists user_role;
