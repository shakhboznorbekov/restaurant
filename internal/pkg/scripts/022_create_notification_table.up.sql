create type notification_type as enum ('NEW', 'SENT', 'CANCELLED');

create table if not exists notifications (
    id bigserial primary key not null,
    title jsonb,
    description jsonb,
    photo varchar,
    status notification_type not null default 'NEW',
    restaurant_id bigint references restaurants (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    updated_by bigint references users(id),
    deleted_at timestamp with time zone,
    deleted_by bigint references users(id)
);

create table if not exists notification_views (
    id bigserial primary key ,
    created_by bigint references users (id),
    notification_id bigint references notifications (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp
    );

alter table notification_views drop constraint if exists u_notification_view;

alter table notification_views add constraint u_notification_view unique (notification_id, created_by);