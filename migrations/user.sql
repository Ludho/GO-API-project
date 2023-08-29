create table if not exists "user"(
    id varchar primary key,
    username varchar not null UNIQUE,
    pass varchar not null
);

create unique index if not exists index_users_id on "user" using btree(id);