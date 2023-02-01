-- Скалярная функция
drop function get_latest_reg_date();
create or replace function get_latest_reg_date()
returns date
as $$
	select max(registration_date)
	from public.users
$$ language sql;



select *
from public.users u
where u.registration_date = get_latest_reg_date()



-- Подставляемая табличная функция
drop function get_user(user_id int);
create or replace function get_user(user_id int)
returns public.users
as $$
	select *
	from public.users
	where public.users.id = user_id
$$ language sql;



select *
from get_user(557)



-- Многооператорная табличная функция
drop function get_terms_by_date(date_beg date, date_end date);
create or replace function get_terms_by_date(date_beg date, date_end date)
returns table 
(
	id int,
	model_id int,
	class int,
	words_count int,
	registration_date date,
	text text
)
as $$
begin
	return query
	select *
	from public.terms_ru
	where public.terms_ru.registration_date between date_beg and date_end;
	
	return query
	select *
	from public.terms_en
	where public.terms_ru.registration_date between date_beg and date_end;
end;
$$ language plpgsql;



select *
from get_terms_ru_by_date('2022-05-01', '2022-05-31')



-- Рекурсивная функция или функция с рекурсивным ОТВ
drop function get_terms_by_level(req_level int);
create or replace function get_terms_by_level(req_level int)
returns table 
(
	res_id int,
	res_text text,
	res_registration_date date,
	res_level int
)
as $$
BEGIN
	RETURN QUERY
	SELECT id, text, registration_date, level
	FROM (
		SELECT id, text, registration_date, level, ROW_NUMBER() OVER(PARTITION BY id) as row_number
		FROM (
			-- Определение ОТВ
			WITH RECURSIVE long_models(id, text, registration_date, level) AS
			(
				-- Определение закрепленного элемента
				SELECT m.id, m.text, m.registration_date, 0 AS level
				FROM public.models AS m
				WHERE words_count = 1
				UNION 
				-- Определение рекурсивного элемента
				SELECT m.id, m.text, m.registration_date, level + 1
				FROM public.models AS m JOIN long_models AS l
					 ON (m.text LIKE CONCAT('%+', l.text, '%')) OR (m.text LIKE CONCAT('%', l.text, '+%'))
			) 
			-- Инструкция, использующая ОТВ
			SELECT id, text, registration_date, level
			FROM long_models
			ORDER BY id ASC, level DESC
		) AS ordered_models
	) AS numbered_models
	WHERE row_number = 1 and level = req_level;
END;
$$ language plpgsql;



select *
from get_terms_by_level(3);



-- Хранимая процедура без параметров или с параметрами
drop procedure insert_term_en;
create or replace procedure insert_term_en
(
    id int,
	model_id int,
	class int,
	words_count int,
	registration_date date,
	text text
)
as $$
begin
    insert into terms_en
    values (id, model_id, class, words_count, registration_date, text);
end;
$$ language plpgsql;

call insert_term_en(5000, 500, 2, 9, '2022-07-07', 'Andreys fireplug businessmen telltales Botswanas embank detective shambles silkscreens');



-- Рекурсивная хранимая процедура или хранимая процедура с рекурсивным ОТВ
drop procedure get_models_tree(root_id int, level int);
create or replace procedure get_models_tree(root_id int, level int)
as $$
declare
	cur_text text;
	child record;
begin
	cur_text = (select text
			    from public.models
			    where public.models.id = root_id
			    limit 1);

	raise notice 'level: %, model:  %', level, cur_text;
    FOR child IN 
		(select id, text
	    from public.models m
	    where (m.text LIKE CONCAT('%+', cur_text, '%')) 
	   		   OR (m.text LIKE CONCAT('%', cur_text, '+%')))
        LOOP
            CALL get_models_tree(child.id, level + 1);
        END LOOP;
    -- CALL fib_proc_index(res, index - 1, second, first + second);
	--raise notice 'level: %, model:  %', level, cur_text;
end;
$$ language plpgsql;

call get_models_tree(500, 0);



-- Хранимая процедура с курсором
select *
into user_tmp_cursor
from public.models;

-- Меняет class всех моделей с req_class на new_class.
drop PROCEDURE proc_update_cursor;
CREATE OR REPLACE PROCEDURE proc_update_cursor
(
    req_class INT,
    new_class INT
)
AS $$
DECLARE
    my_cursor CURSOR FOR
        SELECT *
        FROM user_tmp_cursor
        WHERE class = req_class;
    tmp user_tmp_cursor;
BEGIN
    OPEN my_cursor;
   
    LOOP
        -- FETCH - Получает следующую строку из курсора
        -- И присваевает в переменную, которая стоит после INTO.
        -- Если строка не найдена (конец), то присваевается значение NULL.
        FETCH my_cursor
        INTO tmp;
        -- Выходим из цикла, если нет больше строк (Т.е. конец).
        EXIT WHEN NOT FOUND;
        UPDATE user_tmp_cursor
        SET class = new_class
        WHERE user_tmp_cursor.class = tmp.class;
        RAISE NOTICE 'Elem =  %', tmp; -- Debug print
    END LOOP;
   
    CLOSE my_cursor;
END;
$$ LANGUAGE  plpgsql;



CALL proc_update_cursor(3, 4);

SELECT * FROM user_tmp_cursor;



-- Хранимая процедура доступа к метаданным

-- Получаем название атрибутов и их тип.
drop PROCEDURE metadata(name VARCHAR);
CREATE OR REPLACE PROCEDURE metadata(name VARCHAR) -- Получает название таблицы.
AS $$
DECLARE
    myCursor CURSOR FOR
        SELECT column_name, data_type
        FROM information_schema.columns -- INFORMATION_SCHEMA обеспечивает доступ к метаданным о базе данных. columns - данные о столбацых.
        WHERE table_name = name;
    tmp RECORD; -- RECORD - переменная, которая подстравивается под любой тип.
BEGIN
        OPEN myCursor;
        LOOP
            FETCH myCursor
            INTO tmp;
            EXIT WHEN NOT FOUND;
            RAISE NOTICE 'column name = %; data type = %', tmp.column_name, tmp.data_type;
        END LOOP;
        CLOSE myCursor;
END;
$$ LANGUAGE plpgsql;

CALL metadata('users');



-- Триггер AFTER
drop FUNCTION update_trigger();
CREATE OR REPLACE FUNCTION update_trigger()
RETURNS TRIGGER
AS $$
BEGIN
    --RAISE NOTICE 'Old =  %', old; 
   	--RAISE NOTICE 'New =  %', new;
    UPDATE public.users
    SET registration_date = new.registration_date;
    RETURN new; -- Для операций INSERT и UPDATE возвращаемым значением должно быть NEW.
END;
$$ LANGUAGE plpgsql;

-- AFTER - определяет, что заданная цункция будет вызываться после события.
drop trigger if exists log_update on public.users;
CREATE OR REPLACE TRIGGER log_update
AFTER INSERT ON public.users
FOR EACH ROW
EXECUTE PROCEDURE update_trigger();
-- Триггер с пометкой FOR EACH ROW вызывается один раз для каждой строки,
-- изменяемой в процессе операции.

-- UPDATE user_tmp
-- SET id = 20
-- WHERE id = 1;
insert into public.users values (5009, 'Andrey Sapozhkov', '2015-06-01', '2021-01-28', '2022-01-03', 7);

SELECT * FROM users;



-- Триггер INSTEAD OF

-- VIEW - представление
-- Необходимо создавать, так как у обычной таблицы не может быть триггера INSTEAD OF
drop VIEW tmp_users;
CREATE VIEW tmp_users AS
SELECT *
FROM users
WHERE id < 19;



drop FUNCTION delete_trigger();
CREATE OR REPLACE FUNCTION delete_trigger()
RETURNS TRIGGER
AS $$
BEGIN
    --RAISE NOTICE 'Old =  %', old; 
   	--RAISE NOTICE 'New =  %', new;
    UPDATE tmp_users
    SET name = 'deleted_user'
    WHERE tmp_users.id = old.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

drop trigger if exists user_delete on tmp_users;
CREATE OR REPLACE TRIGGER user_delete
INSTEAD OF DELETE ON tmp_users
FOR EACH ROW
EXECUTE PROCEDURE delete_trigger();



delete from tmp_users
where id = 7;

SELECT * FROM tmp_users;



-- Защита:
drop FUNCTION update_trigger();
CREATE OR REPLACE FUNCTION update_trigger()
RETURNS TRIGGER
AS $$
BEGIN
    --RAISE NOTICE 'Old =  %', old; 
   	--RAISE NOTICE 'New =  %', new;
    UPDATE public.users
    SET last_addition_date = new.last_addition_date;
    RETURN new; -- Для операций INSERT и UPDATE возвращаемым значением должно быть NEW.
END;
$$ LANGUAGE plpgsql;

-- AFTER - определяет, что заданная цункция будет вызываться после события.
drop trigger if exists log_update on public.users;
CREATE OR REPLACE TRIGGER log_update
BEFORE INSERT ON public.users
FOR EACH ROW
EXECUTE PROCEDURE update_trigger();
-- Триггер с пометкой FOR EACH ROW вызывается один раз для каждой строки,
-- изменяемой в процессе операции.

-- UPDATE user_tmp
-- SET id = 20
-- WHERE id = 1;
insert into public.users values (5017, 'Andrey Sapozhkov', '2015-06-01', '2021-01-28', '2022-01-03', 7);

SELECT * FROM users;
