create table if not exists regions (
    id bigserial primary key not null,
    name jsonb not null,
    code int,
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists districts (
    id bigserial primary key not null,
    name jsonb not null,
    region_id bigint references regions (id),
    code int,
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);
