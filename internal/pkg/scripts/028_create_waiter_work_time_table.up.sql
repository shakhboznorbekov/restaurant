create table if not exists waiter_work_time (
    id serial primary key,
    waiter_id bigint references users(id),
    date date,
    periods jsonb
);