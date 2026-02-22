-- 監視ソース
create table sources (
  id          uuid primary key default gen_random_uuid(),
  name        text not null,
  url         text not null,
  type        text not null check (type in ('rss', 'scrape')),
  enabled     boolean default true,
  keywords    text[],
  created_at  timestamptz default now()
);

-- 既読URL管理
create table seen_items (
  id         uuid primary key default gen_random_uuid(),
  source_id  uuid references sources(id) on delete cascade,
  url        text not null unique,
  title      text,
  found_at   timestamptz default now()
);

-- 通知先設定
create table notification_channels (
  id          uuid primary key default gen_random_uuid(),
  type        text not null check (type in ('line', 'google_calendar')),
  config      jsonb not null,
  enabled     boolean default true,
  created_at  timestamptz default now()
);

-- 初期データ：マクドナルドのRSSを登録
insert into sources (name, url, type, enabled, keywords) values
  ('マクドナルド', 'https://www.mcdonalds.co.jp/rss/news.rss', 'rss', true, array['期間限定','新発売','季節限定']),
  ('スターバックス', 'https://www.starbucks.co.jp/press_release/', 'scrape', true, array['期間限定','新発売','季節限定']),
  ('ケンタッキー', 'https://www.kfc.co.jp/news/', 'scrape', false, array['期間限定','新発売','季節限定']);

-- 初期データ：LINE通知チャネルを登録
insert into notification_channels (type, config, enabled) values (
    'line',
    '{"token": "vAoUSz/NZpfT12+8cRO4wta00QTMFo3oTJqpLFFgseFY18vtQh1CA2/eJbDsuFmUZQQ3SpL99zvmgXPi001KWhynEwVDBLyhRhyfMK1KxHdbiR5ZBekZ42jppDof6n4q+nuve8sMCJUGlUjN1CraZgdB04t89/1O/w1cDnyilFU=", "user_id": "Uf97c115e5c8fef94c41b90e5807899d6"}',
    true
  );