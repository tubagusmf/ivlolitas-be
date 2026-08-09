
-- +migrate Up
create type inventory_transaction_type as enum ('RESTOCK', 'SALE', 'RETURN', 'DANGER', 'ADJUSTMENT', 'RELEASE');

create table inventory_transactions (
    id uuid primary key,
    product_variant_id uuid not null references product_variants(id),
    transaction_type inventory_transaction_type not null,
    quantity integer not null,
    stock_before integer not null,
    stock_after integer not null,
    reserved_stock_before integer not null,
    reserved_stock_after integer not null,
    created_by uuid not null references users(id),
    created_at timestamp default current_timestamp
);

create index idx_inventory_transactions_product_variant_id on inventory_transactions (product_variant_id);

-- +migrate Down
drop table inventory_transactions;

drop type inventory_transaction_type;

drop index idx_inventory_transactions_product_variant_id;