create table if not exists printers (
    id bigserial primary key not null ,
    name       varchar,
    ip         varchar,
    port       varchar,
    warehouse_id  bigint references warehouses (id),
    branch_id  bigint references branches (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);