
-- +migrate Up
create table categories (
    id serial primary key,
    name varchar(150) not null,
    slug varchar(150) not null unique,
    description text,
    image_url text,
    is_active boolean default true,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create index idx_categories_name on categories (name);
create index idx_categories_slug on categories (slug);

-- +migrate Down
drop table categories;

drop index idx_categories_name;
drop index idx_categories_slug;