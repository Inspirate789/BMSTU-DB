-- 2) Найти сотрудников, которые не выходят с рабочего места в течение всего рабочего дня.
select distinct employee_id
from times
where type = 2 and date = '2022-12-13'
group by employee_id, date
having count(*) = 1;
