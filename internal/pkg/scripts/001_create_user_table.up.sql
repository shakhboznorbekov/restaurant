alter table schema_migrations
    add column if not exists error text;

create table if not exists users
(
    id         bigserial primary key not null,
    name       text                  not null,
    birth_date date,
    phone      varchar(12),
    role       text,
    gender     text,
    password   text,
    status     text default 'inactive',
    rating real default 0.0,
    photo      text,
    address    text,
    attendance_status bool,
    service_percent bigint references service_percentage(id),
    branch_id bigint references branches (id),
    restaurant_id bigint references restaurants (id),
    created_by bigint references users (id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users (id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users (id)
);

insert into users (name, birth_date, phone, role, gender)
values ('super-admin', '12-12-2000'::date, '998905656436', 'SUPER-ADMIN', 'M');

create table if not exists credit_cards
(
    id         bigserial primary key not null,
    number     varchar(16)           not null,
    expires_at varchar(4)            not null, -- [ MMYY ]
    phone      varchar(12),
    is_main    boolean                  default 'true',
    user_id    bigint references users (id),
    created_by bigint references users (id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users (id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users (id)
    );