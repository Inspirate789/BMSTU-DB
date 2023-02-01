drop table if exists times cascade;
drop table if exists employees cascade;

create table if not exists employees(
    id INT primary key,
    name VARCHAR(64),
    date_of_birth DATE,
    department VARCHAR(64)
);

insert into employees(
	id,
	name,
	date_of_birth, 
	department
) values 
	(0, 'E0', '1991-12-16', 'FIN'),
	(1, 'E1', '1992-12-16', 'FIN'),
	(2, 'E2', '1993-12-16', 'IT'),
	(3, 'E3', '1994-12-16', 'IT'),
	(4, 'E4', '1995-12-16', 'IT');

select * from employees;



create table if not exists times(
    employee_id INT,
    foreign key (employee_id) references employees(id),
    date DATE default CURRENT_DATE,
    day_of_week VARCHAR(64),
    time TIME default CURRENT_TIME,
    type INT check (type = 1 or type = 2)
);

insert into times(
	employee_id, 
	date, 
	day_of_week, 
	time, 
	type
) values
	(0, '2022-12-16', 'Пятница', '09:00', 1),
	(0, '2022-12-16', 'Пятница', '09:12', 2),
	(3, '2022-12-16', 'Четверг', '09:05', 1),
	(3, '2022-12-16', 'Четверг', '09:35', 2),
	(3, '2022-12-16', 'Среда', '09:07', 1),
	(3, '2022-12-16', 'Среда', '09:45', 2),
	(3, '2022-12-16', 'Вторник', '09:46', 1),
	(3, '2022-12-16', 'Вторник', '09:50', 2),
	(3, '2022-12-16', 'Понедельник', '09:46', 1),
	(3, '2022-12-16', 'Понедельник', '09:50', 2),
	(4, '2022-12-16', 'Четверг', '09:05', 1),
	(4, '2022-12-16', 'Четверг', '09:35', 2),
	(4, '2022-12-16', 'Среда', '09:07', 1),
	(4, '2022-12-16', 'Среда', '09:45', 2),
	(4, '2022-12-16', 'Вторник', '09:46', 1),
	(4, '2022-12-16', 'Вторник', '09:50', 2),
	(4, '2022-12-16', 'Понедельник', '09:46', 1),
	(4, '2022-12-16', 'Понедельник', '09:50', 2);

select * from times;

CREATE OR REPLACE FUNCTION Visit(dt DATE)
RETURNS TABLE
(
    minutes double precision,
    employee_qty int
)
AS
$$
    SELECT EXTRACT (HOURS FROM time - '09:00:00') * 60 + EXTRACT (MINUTES FROM time - '09:00:00'), COUNT(*) AS employee_qty
    FROM times
    WHERE date = dt
    AND time > '09:00:00'
    AND type = 1
    GROUP BY time - '09:00:00'
$$ LANGUAGE SQL;

SELECT * FROM Visit('2020-11-17');



-- Написать функцию, возвращающую количество опоздавших сотрудников. Дата опоздания передается в качестве параметра. 
create or replace function late(dt DATE)
returns INT
as
$$
    select COUNT(*) as employee_qty -- , extract (HOURS FROM time - '09:00:00') * 60 + EXTRACT (MINUTES FROM time - '09:00:00')
    from times
    where date = dt
    and time > '09:00:00'
    and type = 1
    -- group by time - '09:00:00'
$$ language sql;

select * from late('2022-12-17');


-- Найти отделы, в которых хоть один сотрудник опаздывает больше 3-х раз в неделю.
SELECT date_part, COUNT(employee_id) AS cnt, department 
FROM
(
	select employee_id, EXTRACT(WEEK FROM date) AS date_part, COUNT(*) as late_cnt
	from times t
	where type = 1 
	group by employee_id, date_part
	having min(time) > '9:00'
) AS tmp join employees e on tmp.employee_id = e.id 
where late_cnt > 3
GROUP BY date_part, department;
--HAVING COUNT(employee_id) > 3;



-- Найти средний возраст сотрудников, не находящихся на рабочем месте 8 часов в день.
select avg(EXTRACT(YEAR FROM CURRENT_DATE) - EXTRACT(YEAR FROM date_of_birth))
from employees join
	(
	select distinct on (employee_id, date) employee_id, date, sum(tmp_dur) over (partition by employee_id, date) as day_dur
	from
		(
		select employee_id, date, time, 
			type, 
			lag(time) over (partition by employee_id, date order by time) as prev_time, 
			time-lag(time) over (partition by employee_id, date order by time) as tmp_dur
		from times t 
		order by employee_id, date, time
		) as small_durations
	) as day_durations
on employees.id = day_durations.employee_id
where day_durations.day_dur < '08:00:00';



-- Вывести все отделы и количество сотрудников хоть раз опоздавших за всю историю учета.
SELECT e.department, COUNT(distinct e.id)
FROM employees e
INNER JOIN times t ON (t.employee_id = e.id)
WHERE ((t.time > '09:00:00') AND (t.type = 1))
GROUP by e.department;


