drop type if exists transaction_type;

create type transaction_type as enum ('INCOME', 'OUTCOME');

create table if not exists warehouse_transaction_history (
    id bigserial primary key not null,
    amount float8 not null,
    type transaction_type not null,
    warehouse_id bigint references warehouses (id),
    product_id bigint references products (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists warehouse_product (
    id bigserial primary key not null,
    amount float8 default 0,
    price float8 not null,
    product_id bigint references products (id),
    warehouse_transaction_id bigint references warehouse_transaction_history (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists warehouse_state (
    id bigserial primary key not null,
    amount float8 not null,
    average_price float8 not null,
    product_id bigint references products (id),
    warehouse_id bigint references warehouses (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
    );