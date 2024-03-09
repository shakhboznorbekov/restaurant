create table if not exists service_percentage (
    id bigserial primary key,
    percent real default 0.0,
    branch_id bigint references branches(id),
    created_by bigint references users (id),
    created_at timestamp default current_timestamp,
    deleted_by bigint references users (id),
    deleted_at timestamp
);