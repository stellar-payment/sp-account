create table users (
    id uuid primary key,
    username varchar(255) not null,
    password varchar(255) not null,
    role_id smallint,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);
