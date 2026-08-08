
-- +migrate Up
create table product_variants (
    id uuid primary key,
    product_id uuid not null references products(id),
    sku varchar(50) not null unique,
    color varchar(50),
    size varchar(50),
    price numeric(15, 2) not null,
    weight numeric(10, 2) default 0,
    barcode varchar(100),
    is_active boolean default true,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create index idx_product_variants_product_id on product_variants (product_id);
create index idx_product_variants_sku on product_variants (sku);

-- +migrate Down
drop table product_variants;

drop index idx_product_variants_product_id;
drop index idx_product_variants_sku;