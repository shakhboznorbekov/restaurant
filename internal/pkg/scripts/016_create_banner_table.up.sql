create table if not exists banners (
    id bigserial primary key not null ,
    title jsonb,
    description jsonb,
    photo text not null,
    price real,
    old_price real,
    status text default 'DRAFT', -- [approved, draft, cancelled]
    expired_at timestamp with time zone default current_timestamp,
    menu_ids bigint[] default '{}',
    branch_id bigint references branches(id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);