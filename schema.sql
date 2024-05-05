create table t1 (
  id bigserial
  , val text not null unique
  , created_at timestamp with time zone not null
  , primary key(id)
);

create table duplicate_t1 (
  id bigserial
  , val text not null
  , created_at timestamp with time zone not null
  , primary key(id)
);
