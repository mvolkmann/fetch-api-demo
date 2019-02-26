-- This assumes that "createdb survey" has already been run.

drop table if exists dog;

create table dog (
  id serial primary key,
  breed text,
  name text
);
