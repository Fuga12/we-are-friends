-- +goose Up
create table summer_shifts (
  id bigserial primary key, title text not null, age_range text not null,
  description text not null, schedule text not null, price text not null,
  emoji text not null, color char(7) not null, cover_media_id bigint references media(id) on delete set null,
  full_description text[] not null default '{}', highlights text[] not null default '{}',
  position int not null default 0, is_published boolean not null default true,
  publication_status text not null default 'published' check (publication_status in ('draft','published','scheduled')),
  publish_at timestamptz,
  created_at timestamptz not null default now(), updated_at timestamptz not null default now()
);
create table summer_shift_photos (
  summer_shift_id bigint not null references summer_shifts(id) on delete cascade,
  media_id bigint not null references media(id) on delete cascade, position int not null default 0,
  primary key (summer_shift_id, media_id)
);
create index summer_shifts_public_order on summer_shifts (position, id) where is_published;

insert into summer_shifts (title, age_range, description, schedule, price, emoji, color, full_description, highlights, position) values
('Школа магии: Тайная миссия', '7–12 лет', 'Ребята становятся учениками школы волшебства: квесты, опыты и магические изделия.', '13–24 июля', 'по запросу', '🪄', '#ffa513', array['Захватывающий квест, удивительные опыты и магические изделия на столярном мастер-классе. Каждый участник выращивает собственную микрозелень и забирает садик домой.'], array['Квест и магические испытания','Опыты и эксперименты','Столярный мастер-класс','Своя микрозелень'], 0),
('Деффчонки', '7–12 лет', 'Творческие мастер-классы, фотосессия, стиль, походы и пижамная вечеринка.', '27 июля – 7 августа', 'по запросу', '💖', '#ffa513', array['Самая нежная смена: изделия своими руками и профессиональная фотосессия. Уход за собой, стиль и уверенность, а также поход, боулинг, лазертаг и пижамная вечеринка.'], array['Профессиональная фотосессия','Творческие мастер-классы','Стиль и уход за собой','Поход, боулинг, лазертаг'], 1),
('Майнкрафт: от игры к роботу', '7–13 лет', 'Робототехника, программирование, мини-роботы, бассейн, походы и кино.', '10–28 августа', 'по запросу', '🤖', '#ffa513', array['Смена для любителей Minecraft: знакомство с робототехникой и программированием, создание мини-роботов и гонки роботов «Формула 1». Можно выбрать участие на одну, две или три недели.'], array['Робототехника и программирование','Гонки роботов «Формула 1»','Бассейн и походы','Участие на 1–3 недели'], 2);

-- +goose Down
drop table if exists summer_shift_photos, summer_shifts;
