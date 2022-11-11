create table articles (
  id integer primary key,
  title text not null,
  content text not null,
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
) strict;

create index articles_created_idx on articles (created);

create trigger articles_updated_timestamp after update on articles begin
  update articles set updated = (strftime('%Y-%m-%dT%H:%M:%fZ')) where id = old.id;
end;
