create table if not exists feedback (
    id bigserial primary key not null ,
    name jsonb
);

alter table order_food drop if exists u_order_food;

alter table order_food add constraint u_order_food unique (order_id, food_id);