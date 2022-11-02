create table articles (
  id integer primary key,
  title text not null,
  content text not null,
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) strict;

create index articles_created_idx on articles (created);
