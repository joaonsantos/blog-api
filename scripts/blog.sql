create table if not exists posts (
  slug   text         primary key,
  title  varchar(100) not null,
  body   text         not null,
  author varchar(25)  not null,
  date   timestamp    not null
);
