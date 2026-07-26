
-- +migrate Up
create table users (
    id uuid primary key,
    role_id integer not null references roles(id),
    full_name varchar(150) not null,
    email varchar(150) not null unique,
    password text not null,
    address text not null,
    phone_number varchar(20) not null,
    is_active boolean default true,
    last_login timestamp default null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp default null
);

create index idx_users_full_name on users (full_name);
create index idx_users_email on users (email);
create index idx_users_role_id on users (role_id);

-- +migrate Down
drop table users;

drop index idx_users_full_name;
drop index idx_users_email;
drop index idx_users_role_id;
