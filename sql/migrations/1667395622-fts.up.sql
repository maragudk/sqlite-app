create virtual table articles_fts
  using fts5(title, content, tokenize = porter, content = 'articles', content_rowid = 'id');

create trigger articles_after_insert after insert on articles begin
  insert into articles_fts (rowid, title, content) values (new.id, new.title, new.content);
end;

create trigger articles_fts_after_update after update on articles begin
  insert into articles_fts (articles_fts, rowid, title, content) values('delete', old.id, old.title, old.content);
  insert into articles_fts (rowid, title, content) values (new.id, new.title, new.content);
end;

create trigger articles_fts_after_delete after delete on articles begin
  insert into articles_fts (articles_fts, rowid, title, content) values('delete', old.id, old.title, old.content);
end;

