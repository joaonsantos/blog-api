create table if not exists posts (
  id    text primary key,
  title text not null,
  body  text not null,
  date  timestamp not null
);
