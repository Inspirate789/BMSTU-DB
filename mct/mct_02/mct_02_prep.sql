-- Создать хранимую процедуру с входным параметром, которая выводит
-- имена хранимых процедур, созданных с параметром WITH RECOMPILE, в
-- тексте которых на языке SQL встречается строка, задаваемая параметром
-- процедуры. Созданную хранимую процедуру протестировать. 

-- WITH RECOMPILE есть только в MSSQL, поэтому в Postgres выкручиваемся без него

-- 1 способ:
create or replace procedure info_routine(str VARCHAR(32)) as '
DECLARE
    elem RECORD;
BEGIN
    FOR elem in
        SELECT routine_name, routine_type
        FROM information_schema.routines -- Чтобы были наши схемы.
        WHERE specific_schema = ''public''
        AND (routine_type = ''PROCEDURE''
        OR (routine_type = ''FUNCTION'' AND data_type != ''record''))
        AND routine_definition LIKE CONCAT(''%'', str, ''%'')
    LOOP
        RAISE NOTICE ''elem: %'', elem;
    END LOOP;
END;
' language plpgsql;

call info_routine('SELECT');

-- 2 способ:
select *
from information_schema.routines
where specific_schema = 'public';





-- Создать хранимую процедуру без параметров, которая в текущей базе
--данных обновляет все статистики для таблиц в схеме 'dbo'. Созданную
--хранимую процедуру протестировать.
--help: Чтобы обновить статистику оптимизации запросов для таблицы в
--указанной базе данных, необходимо
--воспользоваться инструкцией UPDATE STATISTICS, которая в простейшем
--случае имеет следующий формат: 

-- UPDATE STATISTICS есть только в MSSQL, поэтому в Postgres выкручиваемся без него

-- Способ 1:
create or replace procedure upd_stat()
as $$
begin
	analyse verbose employee;
end;
$$ language plpgsql;

call upd_stat();



-- Способ 2:
create or replace procedure upd_stat_2(table_name_in VARCHAR(32))
as '
DECLARE
    elem RECORD;
BEGIN
    FOR elem in
        SELECT relname, last_analyze
        FROM pg_stat_all_tables
        WHERE relname = table_name_in
        LOOP
            RAISE NOTICE ''elem: %'', elem;
        END LOOP;
END;
' LANGUAGE plpgsql;

CALL upd_stat_2('employee');


-- Я не уверен, что сделал то, что надо, 
-- поэтому прикрепляю копипаст из других источников, мб вы что-то лучще накостылите...

-- Источник 1:
SELECT * FROM pg_catalog.pg_stat_activity psa;

SELECT * 
from information_schema.ROUTINES
WHERE specific_schema = 'public';

SELECT * FROM pg_stats WHERE tablename = 'tbl';
UPDATE pg_stats;



SELECT * FROM pg_stat_rate;



SELECT * FROM rate;

EXPLAIN ANALYSE SELECT count(*) FROM rate; 
UPDATE STATISTICS s1; 

CREATE STATISTICS s1 (dependencies) ON id,sale FROM rate;

SELECT * FROM pg_stat_user_tables;



-- Источник 2:
SELECT * FROM pg_catalog.pg_stat_activity psa 

select specific_catalog, specific_schema, specific_name, routine_definition
from information_schema.ROUTINES
WHERE specific_schema = 'public'

SELECT * FROM pg_stats WHERE tablename = 'tbl'
UPDATE pg_stats

EXECUTE 'SQL CONNECT TO "rk2_2"';

EXECUTE 'select * from rate'



-- Создать хранимую процедуру с выходным параметром, которая уничтожает 
-- все SQL DDL триггеры (триггеры типа 'TR') в текущей базе данных. Выходной
-- параметр возвращает количество уничтоженных триггеров. Созданную хранимую 
-- процедуру протестировать. 

create table if not exists rate(
	a int
)

SELECT * FROM information_schema.triggers t;

SELECT * FROM pg_catalog.pg_trigger pt;

SELECT * FROM pg_catalog.pg_event_trigger pet;

CREATE OR REPLACE FUNCTION update_trigger()
RETURNS TRIGGER 
AS '
BEGIN
	RAISE NOTICE ''New =  %'', new;
    RAISE NOTICE ''Old =  %'', old; RAISE NOTICE ''New =  %'', new;
	
	RETURN new;
END;
' LANGUAGE plpgsql;

CREATE TRIGGER update_my
AFTER UPDATE ON rate
FOR EACH ROW 
EXECUTE PROCEDURE update_trigger();

-- по факту ЭТО УДАЛИТЬ DDL триггеры
-- НОООО в посгрессе нет DDL триггеров
-- там только DML.....
-- ясен пень, что ничего не удалит
CREATE OR REPLACE PROCEDURE drop_trigger(count_ INOUT int)  
AS 
$$
DECLARE 
    tmp_trigger_name record;
    cursor_trigger_name CURSOR FOR
	    SELECT trigger_name, event_object_table 
	    FROM information_schema.triggers;
	    --WHERE event_manipulation = 'CREATE' OR 
	      --    event_manipulation = 'ALTER' OR 
	        --  event_manipulation = 'DROP';
BEGIN  
    OPEN cursor_trigger_name;
    LOOP 
    	
        FETCH cursor_trigger_name INTO tmp_trigger_name;
        EXIT WHEN NOT FOUND;
   		count_ = count_ + 1;
        EXECUTE 'DROP TRIGGER ' || tmp_trigger_name.trigger_name || ' ON ' || tmp_trigger_name.event_object_table;
        RAISE NOTICE 'Trigger "%" was deleted!', tmp_trigger_name.trigger_name;
    END LOOP;

    CLOSE cursor_trigger_name;
END;
$$    LANGUAGE plpgsql;

DROP TRIGGER update_my ON rate;
CALL drop_trigger(0);



