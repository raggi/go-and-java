-- select id, email, created_at, updated_at, name, admin, active from users where apikey = $1
create table users (
       id text not null,
       apikey text not null,
       name text not null,
       email text not null,
       password_hash text not null,
       admin boolean not null default false,
       active boolean not null default true,
       created_at timestamp with time zone not null default localtimestamp,
       updated_at timestamp with time zone not null default localtimestamp
);
