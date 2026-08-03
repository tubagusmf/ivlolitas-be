
-- +migrate Up
create table products (
    id uuid primary key,
    category_id integer not null references categories(id),
    sku varchar(50) not null unique,
    name varchar(200) not null,
    slug varchar(200) not null unique,
    short_description text,
    description text,
    price numeric(15, 2) not null,
    weight numeric(10, 2) default 0,
    stock integer default 0,
    is_active boolean default true,
    created_by uuid not null references users(id),
    updated_by uuid not null references users(id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create index idx_products_category_id on products (category_id);
create index idx_products_sku on products (sku);
create index idx_products_name on products (name);
create index idx_products_created_by on products (created_by);
create index idx_products_updated_by on products (updated_by);

-- +migrate Down
drop table products;

drop index idx_products_category_id;
drop index idx_products_sku;
drop index idx_products_name;
drop index idx_products_created_by;
drop index idx_products_updated_by;