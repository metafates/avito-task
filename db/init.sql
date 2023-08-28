do $$
    begin
    if not exists (select 1 from pg_type where typname = 'audit_action') then
        create type audit_action as enum ('ASSIGN', 'DEPRIVE');
    end if;
end$$;

create table if not exists users (
    id integer primary key
);

create table if not exists segments (
    slug varchar(128) primary key,
    outreach real
);

create table if not exists assigned_segments (
    user_id integer references users(id) on delete cascade,
    segment_slug varchar(128) references segments(slug) on delete cascade,
    expires_at timestamptz
);

create table if not exists assignments_audit (
    user_id integer references users(id) on delete cascade,
    action audit_action not null,
    segment_slug varchar(128) not null,
    stamp timestamptz not null,
    expires_at timestamptz
);

create or replace function tg_assignments_audit() returns trigger
    language plpgsql
as
$$
begin
    if (TG_OP = 'INSERT') then
        insert into assignments_audit (user_id, segment_slug, action, stamp, expires_at)
        values (NEW.user_id, NEW.segment_slug, 'ASSIGN'::audit_action, now(), NEW.expires_at);
        return new;
    elsif (TG_OP = 'DELETE') then
        insert into assignments_audit (user_id, segment_slug, action, stamp)
        values (OLD.user_id, OLD.segment_slug, 'DEPRIVE'::audit_action, now());
        return old;
    end if;
    return null;
end;
$$;

drop trigger if exists assignments_audit on assigned_segments;

create trigger assignments_audit
after insert or delete on assigned_segments
for each row execute procedure tg_assignments_audit();
