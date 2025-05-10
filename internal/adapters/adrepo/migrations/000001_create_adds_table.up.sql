create table adds (
    id serial primary key,
    title varchar(100) not null,
    text varchar(500) not null,
    author_id int not null,
    published bool default false,
    date_created timestamp default current_timestamp,
    date_updated timestamp default current_timestamp
);

create or replace function update_adds_timestamp()
    returns trigger as $$
begin
    new.date_updated = now();
    return new;
end;
$$ language plpgsql;

create trigger trigger_update_adds_timestamp
    before update on adds
    for each row
execute function update_adds_timestamp();