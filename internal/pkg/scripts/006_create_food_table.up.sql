create table if not exists food_category (
    id bigserial primary key not null,
    name text not null,
    logo text,
    main bool not null default false,
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists foods (
    id bigserial primary key not null,
    name text not null,
    photos text[] default '{}',
    category_id bigint references food_category (id),
    restaurant_id bigint references restaurants (id),
    branch_id int references branches(id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists food_price (
    id bigserial primary key not null,
    price float8 not null,
    set_date date default current_date,
    food_id bigint references foods (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists order_menu (
    id bigserial primary key not null,
    count int default ,
    menu_id bigint references menus (id),
    order_id bigint references orders (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);
alter table order_menu add constraint u_order_menu unique (order_id, menu_id);