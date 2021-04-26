create table if not exists posts (
  id           text         not null,
  title        varchar(256) not null,
  body         text         not null,
  summary      varchar(256) not null,
  author       varchar(128) not null,
  readTime     integer      not null,
  createDate   integer      not null,
  constraint   posts_pkey   primary key (id)
);