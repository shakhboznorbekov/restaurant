CREATE TYPE partner_type AS ENUM (
    'COUNTER-AGENT'
);

create table if not exists partners (
    id bigserial primary key not null ,
    name text not null ,
    type partner_type not null default 'COUNTER-AGENT',
    restaurant_id bigint references restaurants(id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists partner_contract (
    id bigserial primary key not null,
    name text not null,
    contact_date date default current_date,
    contract_number int,
    payment_type text default 'cash', -- [cash, card]
    partner_id bigint references partners (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);