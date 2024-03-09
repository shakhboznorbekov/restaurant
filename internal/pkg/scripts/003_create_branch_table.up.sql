create table if not exists branches (
    id bigserial primary key not null,
    location jsonb not null,
    photos text[] default '{}' ,
    status text default 'inactive', -- [inactive, active]
    work_time jsonb,
    name varchar not null,
    rate real default 0.0,
    token_expired_at timestamp,
    token text,
    menu_names text,
    category_id bigint references restaurant_category (id),
    restaurant_id bigint references restaurants (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

/* work_time:
   {
        "monday": {
            "from": "07:00",
            "to": "23:00"
        },
        "tuesday": {
            "from": "09:00",
            "to": "22:00"
        },
        ...
   }
*/

create table if not exists branch_reviews (
    id bigserial primary key not null ,
    point int , -- [1, 2, 3, 4, 5]
    comment text default '-' ,
    user_id bigint references users (id) ,
    branch_id bigint references branches (id) ,
    created_by bigint references users(id) ,
    created_at timestamp with time zone default current_timestamp ,
    updated_at timestamp with time zone ,
    updated_by bigint references users(id) ,
    deleted_at timestamp with time zone ,
    deleted_by bigint references users(id)
);
