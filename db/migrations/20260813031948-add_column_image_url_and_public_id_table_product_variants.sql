
-- +migrate Up
alter table product_variants
add column image_url text,
add column image_public_id text;

-- +migrate Down
alter table product_variants
drop column image_url,
drop column image_public_id;