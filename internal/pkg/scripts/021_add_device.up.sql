create table if not exists devices
(
    id           bigserial primary key not null,
    name         text,
    user_id      integer not null,
    is_log_out   bool default false,
    device_id    text unique,
    device_token text,
    device_lang  text default 'uz',
    created_at   timestamp with time zone default current_timestamp,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone
);