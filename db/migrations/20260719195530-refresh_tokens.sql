
-- +migrate Up
create table refresh_tokens (
    id uuid primary key,
    user_id uuid not null references users(id),
    token text not null,
    expired_at timestamp not null,
    revoked boolean default false,
    created_at timestamp default current_timestamp
);

-- +migrate Down
drop table refresh_tokens;