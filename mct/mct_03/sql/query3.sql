-- 3) Найти все отделы, в которых есть сотрудники, опоздавшие в определённую дату.
select distinct department
from employees e
where id in  (
    select employee_id
    from times
    where type = 1 and date = $1
    group by employee_id
    having min(time) > '09:00'
);