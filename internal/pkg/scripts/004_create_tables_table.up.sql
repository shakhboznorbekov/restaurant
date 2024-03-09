create table tables (
    id bigserial primary key not null,
    number int default '1', -- [1, 2, 3, 4, 5, ...]
    capacity int default '4', -- [2, 4, 10, 20, ...]
    branch_id bigint references branches (id),
    status text default 'inactive',
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);