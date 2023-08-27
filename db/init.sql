create table users (
    id varchar(32) primary key
);

create table segments (
    slug varchar(128) primary key,
    outreach real
);

create table assigned_segments (
    user_id varchar(32) references users(id) on delete cascade,
    segment_slug varchar(128) references segments(slug) on delete cascade,
    expires_at timestamp
);

alter table assigned_segments enable row level security;

create policy filter_expired_assignments on assigned_segments
    for select
    using (expires_at > now() or expires_at is null);

-- create table audit (
--     operation operation,
--     user_id integer references USERS(ID) on delete cascade,
--     segment_slug slug,
--     timestamp timestamp
-- );
