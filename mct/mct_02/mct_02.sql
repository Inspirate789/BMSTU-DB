-- Задание 1

drop database RK2;
create database RK2;

drop table if exists students CASCADE;
drop table if exists teachers CASCADE;
drop table if exists marks CASCADE;
drop table if exists topics CASCADE;
 
CREATE TABLE IF NOT EXISTS students
(
    id serial PRIMARY KEY,
    t_id INT,
    m_num INT,
    name VARCHAR(50),
    dep VARCHAR(50),
    grp INT
);

CREATE TABLE IF NOT EXISTS teachers
(
    id serial PRIMARY KEY,
    name VARCHAR(50),
    deg VARCHAR(50),
    kaf VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS marks
(
    id serial PRIMARY KEY,
    m_num INT,
    m_gov INT,
    m_dip INT
);

CREATE TABLE IF NOT EXISTS topics
(	
	id serial PRIMARY KEY,
    t_id INT,
    title VARCHAR(50)
);

ALTER TABLE students
ADD CONSTRAINT s_t_id FOREIGN KEY (t_id) REFERENCES teachers(id);

ALTER TABLE students
ADD CONSTRAINT m_id UNIQUE (m_num);

ALTER TABLE marks
ADD CONSTRAINT m_s_id FOREIGN KEY (m_num) REFERENCES students(m_num);

ALTER TABLE topics
ADD CONSTRAINT t_t_id FOREIGN KEY (t_id) REFERENCES teachers(id);

INSERT INTO students (m_num, name, dep, grp) 
VALUES
	(1, 'Student1', 'A', 1),
	(2, 'Student2', 'A', 2),
	(3, 'Student3', 'B', 3),
	(4, 'Student4', 'C', 2),
	(5, 'Student5', 'A', 1),
	(6, 'Student6', 'B', 5),
	(7, 'Student7', 'B', 5),
	(8, 'Student8', 'C', 1),
	(9, 'Student9', 'C', 6),
	(10, 'Student10', 'B', 1)
;

INSERT INTO teachers (name, deg, kaf) 
VALUES
	('Teacher1', '1', 'A1'),
	('Teacher2', '1', 'A1'),
	('Teacher3', '3', 'B1'),
	('Teacher4', '2', 'B2'),
	('Teacher5', '2', 'C1'),
	('Teacher6', '1', 'B1'),
	('Teacher7', '3', 'C2'),
	('Teacher8', '3', 'B3'),
	('Teacher9', '2', 'C2'),
	('Teacher10', '1', 'B1')
;

UPDATE students SET t_id = id;

INSERT INTO marks (m_num, m_gov, m_dip) 
VALUES
	(1, 5, 5),
	(2, 4, 5),
	(3, 3, 4),
	(4, 2, 5),
	(5, 5, 4),
	(6, 4, 3),
	(7, 3, 5),
	(8, 5, 2),
	(9, 2, 2),
	(10, 5, 5)
;

INSERT INTO topics (t_id, title) 
VALUES
	(1, 'A'),
	(1, 'B'),
	(2, 'A'),
	(2, 'B'),
	(2, 'B'),
	(3, 'A'),
	(5, 'C'),
	(7, 'D'),
	(8, 'E'),
	(8, 'F')
;

SELECT * FROM students;
SELECT * FROM teachers;
SELECT * FROM marks;
SELECT * FROM topics;



-- Задание 2

-- 1) Вывести данные о студентах с лучшими оценками за госы
SELECT students.id, students.name
FROM students JOIN marks ON students.m_num = marks.m_num
WHERE marks.m_gov  >= ALL (SELECT marks.m_gov
						   FROM marks);

-- 2) Посчитать среднюю оценку студентов за дипломы
SELECT AVG(marks.m_dip) AS Actual_AVG
FROM students JOIN marks ON students.m_num = marks.m_num;

-- 3) Создать временную локальную таблицу студентов-отличников
DROP TABLE IF EXISTS excellent_students;

SELECT students.id, t_id, students.m_num, name, dep, grp
INTO excellent_students
FROM students JOIN marks ON students.m_num = marks.m_num
WHERE marks.m_gov = 5 AND marks.m_dip = 5;

SELECT * from excellent_students;



-- Задание 3

create or replace function ufn_my_func()
returns int
as $$
	select 5
$$ language sql;



-- Вариант 1:
 CREATE OR REPLACE PROCEDURE info_function (count_ INOUT int)
AS '
DECLARE
	i int;
    elem RECORD;
BEGIN
    FOR elem in
        SELECT *
        FROM pg_proc
        WHERE prokind = ''f'' and pronargs > 0
    LOOP
		count_ = count_ + 1;
        RAISE NOTICE ''Name: %, count: %, name_param: %, type_param: %'', 
						elem.proname, elem.pronargs, elem.proargnames, 
						elem.proallargtypes;
    END LOOP;
END;
' LANGUAGE plpgsql;

CALL info_function(0);



-- Вариант 2:
create or replace procedure show_functions() as
$$
declare 
	cur cursor
	for select proname, proargtypes
	from (
		select proname, pronargs, prorettype, proargtypes
		from pg_proc
		where pronargs > 0 and proname LIKE 'ufn%'
	) AS tmp;
	row record;
begin
	open cur;
	loop
		fetch cur into row;
		exit when not found;
		raise notice 'func_name : %, args : %', row.proname, row.proargtypes;
	end loop;
	close cur;
end
$$ language plpgsql;

call show_functions();









