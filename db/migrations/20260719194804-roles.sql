
-- +migrate Up
create table roles (
    id serial primary key,
    name varchar(100),
    description text,
    created_at timestamp default current_timestamp,
    update_at timestamp default current_timestamp,
    deleted_at timestamp default null
);

-- +migrate Down
drop table roles;
