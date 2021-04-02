create table if not exists posts (
  id               text         not null,
  title            text         not null,
  body             text         not null,
  summary          text         not null,
  author           text         not null,
  readTime         integer      not null,
  dateModified     integer      not null,
  constraint posts_pkey primary key (id)
);