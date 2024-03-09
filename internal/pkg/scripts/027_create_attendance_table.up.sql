create table if not exists attendances (
    id bigserial primary key,
    user_id bigint references users(id),
    came_at timestamp default current_timestamp,
    gone_at timestamp
);