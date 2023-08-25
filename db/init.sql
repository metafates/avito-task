CREATE TYPE IF NOT EXISTS OPERATION AS ENUM ('add', 'remove');
CREATE DOMAIN IF NOT EXISTS SLUG AS VARCHAR(128);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
);

CREATE TABLE IF NOT EXISTS segments (
    slug SLUG PRIMARY KEY,
    outreach REAL
);

CREATE TABLE IF NOT EXISTS assigned_segments (
    user INTEGER REFERENCES users(id),
    segment SLUG REFERENCES segments(slug),
    expires_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS changes (
    user_id INTEGER REFERENCES users(id),
    segment_slug SLUG,
    operation OPERATION,
    timestamp TIMESTAMP
);
