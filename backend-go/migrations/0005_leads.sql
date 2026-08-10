-- +goose Up
create table leads (
  id bigserial primary key,
  name text not null,
  phone text not null,
  message text not null default '',
  is_processed boolean not null default false,
  created_at timestamptz not null default now()
);
create index leads_created_at on leads (created_at desc);

-- +goose Down
drop table if exists leads;
