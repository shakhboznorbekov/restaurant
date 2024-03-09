create table if not exists measure_unit (
    id bigserial primary key not null ,
    name text not null ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);

create table if not exists products (
    id bigserial primary key not null,
    name text not null,
    measure_unit_id bigint references measure_unit (id),
    restaurant_id bigint references restaurants (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists food_recipe (
    id bigserial primary key not null ,
    amount float8 not null ,
    product_id bigint references products (id) ,
    food_id bigint references foods (id) ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);