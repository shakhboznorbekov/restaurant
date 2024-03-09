create table if not exists stories (
    id bigserial primary key not null ,
    name text,
    file text not null,
    type text not null,
    duration real not null,
    status text default 'DRAFT', -- [approved, draft, cancelled]
    expired_at timestamp with time zone default current_timestamp,
    restaurant_id int not null references restaurants(id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);