create table if not exists restaurant_category (
    id bigserial primary key not null ,
    name text not null,
    photo text,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
    );

create table if not exists restaurants (
    id bigserial primary key not null,
    name text not null,
    logo text,
    website_url text,
    mini_logo text,
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);