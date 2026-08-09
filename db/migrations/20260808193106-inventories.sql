
-- +migrate Up
create table inventories (
    id uuid primary key,
    product_variant_id uuid not null references product_variants(id),
    stock integer not null default 0,
    reserved_stock integer not null default 0,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

-- +migrate Down
drop table inventories;