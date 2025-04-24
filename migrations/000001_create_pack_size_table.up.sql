create table pack_size
(
    size integer
);

create unique index pack_size_size_uindex
    on pack_size (size);

