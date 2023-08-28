create type audit_action as enum ('ASSIGN', 'DEPRIVE');

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
    expires_at timestamptz
);

create table assignments_audit (
    user_id varchar(32) references users(id) on delete cascade,
    action audit_action not null,
    segment_slug varchar(128) not null,
    stamp timestamptz not null
);

create function tg_assignments_audit() returns trigger
    language plpgsql
as
$$
begin
    if (TG_OP = 'INSERT') then
        insert into assignments_audit (user_id, segment_slug, action, stamp)
        values (NEW.user_id, NEW.segment_slug, 'ASSIGN'::audit_action, now());
        return new;
    elsif (TG_OP = 'DELETE') then
        insert into assignments_audit (user_id, segment_slug, action, stamp)
        values (OLD.user_id, OLD.segment_slug, 'DEPRIVE'::audit_action, now());
        return old;
    end if;
    return null;
end;
$$;

create trigger assignments_audit
after insert or delete on assigned_segments
for each row execute procedure tg_assignments_audit();
