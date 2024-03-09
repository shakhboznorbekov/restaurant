CREATE TABLE IF NOT EXISTS story_views (
    id bigserial primary key not null,
    story_id bigint references stories (id),
    created_by bigint references users(id),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

ALTER TABLE story_views DROP CONSTRAINT IF EXISTS u_story_views;

ALTER TABLE story_views ADD CONSTRAINT u_story_views UNIQUE (story_id, created_by);