-- +goose Up
create table media (
  id bigserial primary key,
  object_key text not null unique,
  mime text not null,
  size_bytes bigint not null,
  width int,
  height int,
  alt text not null default '',
  created_at timestamptz not null default now()
);

create table classes (
  id bigserial primary key, title text not null, age_range text not null,
  description text not null, schedule text not null, price text not null,
  emoji text not null, color char(7) not null, cover_media_id bigint references media(id) on delete set null,
  full_description text[] not null default '{}', highlights text[] not null default '{}',
  position int not null default 0, is_published boolean not null default true,
  created_at timestamptz not null default now(), updated_at timestamptz not null default now()
);
create table class_photos (
  class_id bigint not null references classes(id) on delete cascade,
  media_id bigint not null references media(id) on delete cascade, position int not null default 0,
  primary key (class_id, media_id)
);

create table news (
  id bigserial primary key, title text not null, published_on date not null, category text not null,
  excerpt text not null, emoji text not null, cover_media_id bigint references media(id) on delete set null,
  content text[] not null default '{}', is_published boolean not null default true,
  created_at timestamptz not null default now(), updated_at timestamptz not null default now()
);
create table news_photos (
  news_id bigint not null references news(id) on delete cascade,
  media_id bigint not null references media(id) on delete cascade, position int not null default 0,
  primary key (news_id, media_id)
);

create table services (
  id bigserial primary key, title text not null, tagline text not null, description text not null,
  features text[] not null default '{}', price text not null, emoji text not null, path text not null,
  color char(7) not null, cover_media_id bigint references media(id) on delete set null,
  detail jsonb, position int not null default 0, is_published boolean not null default true
);
create table service_photos (
  service_id bigint not null references services(id) on delete cascade,
  media_id bigint not null references media(id) on delete cascade, position int not null default 0,
  primary key (service_id, media_id)
);

create table staff (
  id bigserial primary key, name text not null, role text not null, bio text not null,
  initials text not null, color char(7) not null, photo_media_id bigint references media(id) on delete set null,
  position int not null default 0, is_published boolean not null default true
);
create table documents (
  id bigserial primary key, title text not null, description text not null,
  file_media_id bigint not null references media(id) on delete restrict,
  display_file_type text, display_file_size text,
  position int not null default 0, is_published boolean not null default true
);
create table gallery (
  id bigserial primary key, media_id bigint not null references media(id) on delete cascade,
  title text not null default '', position int not null default 0
);
create table reviews (
  id bigserial primary key, author text not null, role text not null, text text not null,
  rating smallint not null check (rating between 1 and 5), position int not null default 0,
  is_published boolean not null default true
);
create table site_settings (id boolean primary key default true check (id), data jsonb not null);
create table home_photos (
  media_id bigint primary key references media(id) on delete cascade,
  position int not null default 0
);

create index classes_public_order on classes (position, id) where is_published;
create index news_public_date on news (published_on desc) where is_published;

-- +goose Down
drop table if exists home_photos, site_settings, reviews, gallery, documents, staff, service_photos,
  services, news_photos, news, class_photos, classes, media;
