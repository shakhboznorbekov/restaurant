drop type if exists menu_status;

CREATE TYPE menu_status AS ENUM ('inactive', 'active');

create table if not exists menus (
    id bigserial primary key not null ,
    food_id bigint references foods (id),
    branch_id bigint references branches (id),
    status menu_status default 'inactive', -- [inactive, active]
    new_price real,
    old_price real,
    description jsonb,
    printer_id int references printers (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);
