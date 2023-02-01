drop table if exists public.nifi_users;
create table if not exists public.nifi_users(
    id int generated always as identity primary key,
    name text not null,
    registration_date date not null check (first_addition_date >= registration_date),
    first_addition_date date not null check (last_addition_date >= first_addition_date),
    last_addition_date date not null check (terms_count >= 0),
    terms_count int default 0
);

copy public.nifi_users(name, registration_date, first_addition_date, last_addition_date, terms_count) from '/home/data/users.csv' delimiter ',';
select * from public.nifi_users;

COPY
(
    SELECT row_to_json(u) result 
    FROM (SELECT name, registration_date, first_addition_date, last_addition_date, terms_count
    	  FROM public.nifi_users) u
)
TO '/home/data/users.json';