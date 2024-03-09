create table if not exists product_recipe (
    id bigserial primary key not null ,
    amount float8 not null ,
    product_id bigint references products (id) ,
    recipe_id bigint references products (id) ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);