create or replace function get_user(user_id int)
returns public.users
as $$
    select *
    from public.users
    where public.users.id = user_id
$$ language sql;
