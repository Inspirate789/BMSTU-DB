create or replace function get_latest_reg_date()
returns date
    as $$
    select max(registration_date)
    from public.users
$$ language sql;
