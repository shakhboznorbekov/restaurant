create table if not exists orders (
    id bigserial primary key not null ,
    table_id bigint references tables (id) ,
    user_id bigint references users (id) ,
    number integer,
    client_count int default 1,
    accepted_at timestamp,
    waiter_id bigint references users(id),
    status order_status default 'NEW',
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);

create table if not exists order_payment (
    id bigserial primary key not null ,
    order_id bigint references orders (id) ,
    status varchar default 'unpaid' , -- [unpaid, paid, cancelled]
    price real default 0,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);


drop type if exists order_status ;

create type order_status as enum ('NEW', 'PAID', 'CANCELLED', 'SERVED');
