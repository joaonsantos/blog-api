create table if not exists posts (
  id    serial primary key,
  title text not null,
  body  text not null
);