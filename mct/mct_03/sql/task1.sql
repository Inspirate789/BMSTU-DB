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
      (4, 'E4', '2007-12-16', 'IT'),
      (5, 'E2', '1993-12-16', 'IT'),
      (6, 'E3', '1994-12-16', 'IT'),
      (7, 'E4', '2007-12-16', 'IT'),
      (8, 'E2', '1993-12-16', 'IT'),
      (9, 'E3', '1994-12-16', 'IT'),
      (10, 'E4', '2007-12-16', 'IT'),
      (11, 'E2', '1993-12-16', 'IT'),
      (12, 'E3', '1994-12-16', 'IT'),
      (13, 'E4', '2007-12-16', 'IT');

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
      (0, '2022-12-15', 'Четверг', '09:00', 1),
      (0, '2022-12-15', 'Четверг', '09:12', 2),
      (3, '2022-12-15', 'Четверг', '09:05', 1),
      (3, '2022-12-15', 'Четверг', '09:35', 2),
      (3, '2022-12-14', 'Среда', '09:07', 1),
      (3, '2022-12-14', 'Среда', '09:45', 2),
      (3, '2022-12-13', 'Вторник', '09:46', 1),
      (3, '2022-12-13', 'Вторник', '09:50', 2),
      (3, '2022-12-13', 'Вторник', '09:55', 1),
      (3, '2022-12-13', 'Вторник', '09:58', 2),
      (3, '2022-12-12', 'Понедельник', '09:46', 1),
      (3, '2022-12-12', 'Понедельник', '09:50', 2),
      (4, '2022-12-15', 'Четверг', '09:05', 1),
      (4, '2022-12-15', 'Четверг', '09:35', 2),
      (4, '2022-12-14', 'Среда', '09:07', 1),
      (4, '2022-12-14', 'Среда', '09:45', 2),
      (4, '2022-12-13', 'Вторник', '09:46', 1),
      (4, '2022-12-13', 'Вторник', '09:50', 2),
      (4, '2022-12-12', 'Понедельник', '09:46', 1),
      (4, '2022-12-12', 'Понедельник', '09:50', 2);

select * from times;



-- Задание 1 - скалярная функция
CREATE OR REPLACE FUNCTION task_1()
    RETURNS INT
AS
$$
select count(*)
from employees e
where id in (
    with numbered_times as (
        select employee_id, date, time, type, row_number() over (partition by employee_id, date, type order by time) as num
        from times t
        order by employee_id, date, time, type
    )
    select employee_id
    from numbered_times
    where type = 2
    group by employee_id
    having max(num) > 3
)
  and extract(year from current_date) - extract(year from date_of_birth) between 18 and 40;
$$ LANGUAGE sql;

-- SELECT * FROM task_1();