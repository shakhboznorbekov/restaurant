CREATE TABLE IF NOT EXISTS branch_likes (
    id bigserial primary key not null ,
    user_id bigint references users(id),
    branch_id bigint references branches(id)
);

