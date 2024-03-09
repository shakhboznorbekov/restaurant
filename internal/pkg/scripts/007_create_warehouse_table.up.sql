CREATE TYPE warehouse_type AS ENUM (
    'REAL',
    'VIRTUAL'
);

create table if not exists warehouses (
    id bigserial primary key not null ,
    name text not null ,
    location jsonb not null ,
    warehouse_type not null default 'REAL',
    branch_id bigint references branches (id) ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);
