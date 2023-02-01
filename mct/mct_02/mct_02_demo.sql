-- Задание 1 (4 балла)
-- Создать базу данных RK2. Создать в ней структуру, соответствующую указанной на ER-диаграмме. 
-- Заполнить таблицы тестовыми значениями (не менее 10 в каждой таблице).

-- drop database RK2;
create database RK2;

drop table if exists employee CASCADE;
drop table if exists department;
drop table if exists medication CASCADE;
drop table if exists employee_medication;


 
CREATE TABLE IF NOT EXISTS employee
(
    id serial PRIMARY KEY,
    department_id INT,
    job VARCHAR(15),
    fio VARCHAR(40),
    salary INT
);

CREATE TABLE IF NOT EXISTS department
(
    id serial PRIMARY KEY,
    name VARCHAR(20),
    phone VARCHAR(12),
    manager_id INT
);

CREATE TABLE IF NOT EXISTS medication
(
    id serial PRIMARY KEY,
    name VARCHAR(20),
    instruction VARCHAR(50),
    price INT
);

CREATE TABLE IF NOT EXISTS employee_medication
(	
	id serial PRIMARY KEY,
    employee_id INT,
    medication_id INT
);

alter table employee
add constraint fk_department_id foreign key (department_id) references department(id);

alter table department
add constraint fk_manager foreign key (manager_id) references employee(id);

alter table employee_medication
add constraint fk_employee_id foreign key (employee_id) references employee(id),
add constraint fk_medication_id foreign key (medication_id) references medication(id);

INSERT INTO department (name, phone) 
VALUES
	('a', '01'),
	('b', '02'),
	('c', '03'),
	('d', '04'),
	('e', '05'),
	('f', '06'),
	('g', '07'),
	('h', '08'),
	('i', '09'),
	('j', '10'),
	('k', '11')
;


INSERT INTO employee (department_id, job, fio, salary) 
VALUES
	(1, 'washer1', 'qwe', 100),
	(2, 'washer2', 'rty', 101),
	(3, 'washer3', 'uio', 0),
	(4, 'washer4', 'pas', 1),
	(5, 'washer5', 'dfg', 50),
	(6, 'washer', 'hjk', 1),
	(7, 'washer', 'lzx', 100),
	(8, 'washer1', 'nmk', 4),
	(9, 'washer2', 'cvb', 1000),
	(10, 'washer2', 'qwe', 30),
	(11, 'washer4', 'rty', 1),
	(1, 'notboss1', 'uio', 100),
	(2, 'notboss2', 'dfg', 3)
;

UPDATE department SET manager_id = id;

INSERT INTO medication (name, instruction, price) 
VALUES
	('q', 'w', 100),
	('e', 'r', 1),
	('t', 'y', 4),
	('a', 's', 56),
	('g', 'r', 89),
	('w', 't', 09),
	('v', 'b', 12),
	('h', '2', 1300),
	('q', 'w', 5),
	('q1', 'w2', 100)
;

INSERT INTO employee_medication (employee_id, medication_id) 
VALUES
	(1, 10),
	(2, 9),
	(3, 8),
	(4, 7),
	(5, 6),
	(6, 5),
	(7, 4),
	(8, 3),
	(9, 2),
	(1, 9)
;



-- Задание 2 (6 баллов)
-- Написать к разработанной базе данных 3 запроса, в комментарии указать, что этот запрос делает:

-- 1) Инструкцию SELECT, использующую простое выражение CASE
-- Охарактеризовать зарплату сотрудника как большую, нормальную или маленькую
select fio, salary, case 
	when salary > 200 then 'big'
	when salary > 75 then 'normal'
	else 'small'
	end as is_salary_round
from employee; 


--select fio, salary, case 
--	salary % 10
--	when 0 then 'round'
--	when 5 then 'almost_round'
--	else 'strange'
--	end as is_salary_round
--from employee; 

-- 2) Инструкцию, использующую оконную функцию
-- Вывести данные из таблицы сотрудников и посчитать отклонение зарплаты каждого сотрудника от средней
select fio, job , salary, salary - avg(salary) over (partition by job) as salary_diff
from employee
order by job, salary_diff desc;

-- 3) Инструкцию SELECT, консолидирующую данные с помощью предложения GROUP BY и предложения HAVING
-- Выбрать такие отделы, в которых работает более 1 сотрудника
select department_id, COUNT(*)
from employee
group by department_id
having COUNT(*) > 1;



-- Задание 3 (10 баллов)
-- Создать хранимую процедуру с двумя входными параметрами – имя базы данных и имя таблицы, 
-- которая выводит сведения об индексах указанной таблицы в указанной базе данных.

SELECT * FROM pg_catalog.pg_class;
SELECT * FROM pg_catalog.pg_stat_activity;

CREATE OR REPLACE PROCEDURE index_info
(
    db_name_in VARCHAR(32),
    table_name_in VARCHAR(32)
)
AS '
DECLARE
    elem RECORD;
BEGIN
    FOR elem in
        SELECT *
        FROM pg_indexes
        WHERE tablename = table_name_in
        LOOP
            RAISE NOTICE ''elem: %'', elem;
        END LOOP;
END;
' LANGUAGE plpgsql;

CALL index_info('RK2', 'employee');

SELECT *
FROM pg_index JOIN pg_class ON pg_index.indrelid=pg_class.oid
WHERE relname='employee';

SELECT *
FROM pg_class;

SELECT * 
FROM information_schema."tables" t;

SELECT *
FROM pg_catalog.pg_database pd;









