create table cats (
    id bigserial,
    name text not null
);

insert into cats (name) values ('foo'), ('bar');
