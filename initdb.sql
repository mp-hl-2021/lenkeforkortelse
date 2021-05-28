drop table if exists accounts cascade;
create table accounts
(
    id        serial primary key,
    login     varchar(255) not null,
    password  varchar(255) not null,

    createdAt timestamp without time zone default now(),
    updatedAt timestamp without time zone default now(),

    unique (login)
);

drop table if exists links cascade;
create table links
(
    linkId varchar(255) primary key,
    link text,
    linkStatus int default 0,
    accountId varchar(255)
);