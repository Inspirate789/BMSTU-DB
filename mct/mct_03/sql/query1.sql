-- 1) Найти все отделы, в которых работает более 10 сотрудников.
select department
from employees
group by department
having count(id) > 10;
