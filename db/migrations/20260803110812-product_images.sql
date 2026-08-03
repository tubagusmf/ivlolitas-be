
-- +migrate Up
create table product_images (
    id uuid primary key,
    product_id uuid not null references products(id),
    image_url text not null,
    image_public_id text,
    is_primary boolean default false,
    sort_order integer default 1,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create index idx_product_images_product_id on product_images (product_id);
create index idx_product_images_image_url on product_images (image_url);

-- +migrate Down
drop table product_images;

drop index idx_product_images_product_id;
drop index idx_product_images_image_url;