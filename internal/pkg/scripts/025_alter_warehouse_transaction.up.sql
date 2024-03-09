ALTER TABLE warehouse_transaction_history
    RENAME TO warehouse_transactions;

ALTER TABLE warehouse_transactions
    ADD COLUMN if not exists total_price double precision,
    DROP COLUMN type,
    DROP COLUMN if exists warehouse_id,
    ADD COLUMN if not exists from_warehouse_id bigint references warehouses(id),
    ADD COLUMN if not exists from_partner_id bigint references partners(id),
    ADD COLUMN if not exists to_warehouse_id bigint references warehouses(id),
    ADD COLUMN if not exists to_partner_id bigint references partners(id);

DROP TABLE warehouse_product;

create table if not exists warehouse_state_history (
    id bigserial primary key not null ,
    amount float8 default 0 ,
    average_price double precision not null default 0,
    warehouse_state_id bigint references warehouse_state (id) ,
    warehouse_transaction_id bigint references warehouse_transactions (id) ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);
