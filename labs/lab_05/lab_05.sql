-- 1. Из таблиц базы данных, созданной в первой
-- лабораторной работе, извлечь данные в JSON.

-- Функция row_to_json - Возвращает кортеж в виде объекта JSON.
SELECT row_to_json(u) result FROM public.users u;
SELECT row_to_json(tr) result FROM public.terms_ru tr;
SELECT row_to_json(te) result FROM public.terms_en te;



-- 2. Выполнить загрузку и сохранение JSON файла в таблицу.
-- Созданная таблица после всех манипуляций должна соответствовать таблице
-- базы данных, созданной в первой лабораторной работе.

-- Создаем новую таблицу, чтобы сравнить ее со старой.
-- Да и вообще, чтобы не дропать старую таблицу...
create table if not exists public.users_copy
(
	id int primary key,
	name text,
	registration_date date,
	first_addition_date date,
	last_addition_date date,
	terms_count int default 0
);

-- Копируем данные из таблицы users в файл users.json
-- (В начале нужно поставить \COPY).
COPY
(
    SELECT row_to_json(u) result FROM public.users u
)
TO '/home/json/users.json';

-- Подготовка данных завершена.
-- Собственно далее само задание.

-- Помещаем файл в таблицу БД.
-- Создаем таблицу, которая будет содержать json кортежи.
CREATE TABLE IF NOT EXISTS users_import
(
	doc json
);

-- Теперь копируем данные в созданную таблицу.
-- (Но опять же делаем это с помощью \COPY).
COPY users_import FROM '/home/json/users.json';

SELECT * FROM users_import;

-- В принципе можно было сделать так, но т.к. в условии написано
-- Выгрузить из файла, так что нужно использовтаь copy.
-- CREATE TABLE IF NOT EXISTS users_tmp(doc json);
-- INSERT INTO users_tmp
-- SELECT row_to_json(u) result FROM users u;
-- SELECT * FROM users_tmp;

-- Данный запрос преобразует данные из строки в формате json
-- В табличное представление. Т.е. разворачивает объект из json в табличную строку.
SELECT json_populate_record(null::users_copy, doc) FROM users_import;
-- Или так:
SELECT * FROM users_import, json_populate_record(null::users_copy, doc);

-- Преобразование одного типа в другой null::users_copy
SELECT json_populate_record(CAST(null AS users_copy ), doc) FROM users_import;
-- Или так:
SELECT * FROM users_import, json_populate_record(CAST(null AS users_copy ), doc);

-- Загружаем в таблицу сконвертированные данные из формата json из таблицы users_import.
INSERT INTO users_copy
SELECT id, name, registration_date, first_addition_date, last_addition_date, terms_count
FROM users_import, json_populate_record(null::users_copy, doc);

SELECT * FROM users_copy;



-- 3. Создать таблицу, в которой будет атрибут(-ы) с типом JSON, или
-- добавить атрибут с типом JSON к уже существующей таблице.
-- Заполнить атрибут правдоподобными данными с помощью команд INSERT или UPDATE

-- Создаем таблицу, которая будет содержать
-- Нарушителей в json формате.
CREATE TABLE IF NOT EXISTS mvp_users_json
(
    data json
);

-- Вставляем в нее json строку.
-- json_object - формирует объект JSON.
INSERT INTO mvp_users_json
SELECT * FROM json_object('{user_id, username, terms_count}', '{0, "Andrey", 100500}');

SELECT * FROM mvp_users_json;



-- 4. Выполнить следующие действия:

-- 4.1. Извлечь XML/JSON фрагмент из XML/JSON документа
CREATE TABLE IF NOT EXISTS users_part
(
    id INT,
    name VARCHAR,
    terms_count INT
);

SELECT * FROM users_import, json_populate_record(null::users_part, doc);

SELECT id, terms_count
FROM users_import, json_populate_record(null::users_part, doc)
WHERE terms_count > 500;

-- Оператор -> возвращает поле объекта JSON как JSON.
-- -> - выдаёт поле объекта JSON по ключу.
SELECT * FROM users_import;

SELECT doc->'id' AS id, doc->'name' AS name, doc->'terms_count' AS terms_count
FROM users_import;



-- 4.2. Извлечь значения конкретных узлов или атрибутов XML/JSON документа
CREATE TABLE IF NOT EXISTS context_tmp
(
	doc jsonb
);

INSERT INTO context_tmp VALUES ('{"id":0, "content": {"term_1":"absorbtion", "term_2":"Sentence structure"}}');
INSERT INTO context_tmp VALUES ('{"id":1, "content": {"term_1":"none", "term_2":"verb"}}');
INSERT INTO context_tmp VALUES ('{"id":2, "content": {"term_1":"noun", "term_2":"none"}}');

SELECT * FROM context_tmp;

SELECT doc->'id' AS id, doc->'content'->'term_1' AS term_1
FROM context_tmp;

-- 4.3. Выполнить проверку существования узла или атрибута
-- Проверка вхождения — важная особенность типа jsonb, не имеющая аналога для типа json.
-- Эта проверка определяет, входит ли один документ jsonb в другой.

-- В данном примере проверятся существование инвенторя у пользователя с id=u_id.
CREATE OR REPLACE FUNCTION get_context(u_id jsonb)
RETURNS VARCHAR AS '
    SELECT CASE
               WHEN count.cnt > 0
               THEN ''true''
               ELSE ''false''
               END
    FROM (
             SELECT COUNT(doc -> ''id'') cnt
             FROM context_tmp
             WHERE doc -> ''id'' @> u_id
         ) AS count;
' LANGUAGE sql;

INSERT INTO context_tmp VALUES ('{"id":3, "content": {"term_2":"Sentence structure"}}');
INSERT INTO context_tmp VALUES ('{"id":4, "content": {"term_1":"none"}}');
INSERT INTO context_tmp VALUES ('{"id":5, "content": {}}');

SELECT *, get_context(doc->'id') AS id, 
		  get_context(doc->'term_1') AS term_1
FROM context_tmp;

SELECT get_context('1');

-- 4.4. Изменить XML/JSON документ
SELECT * FROM context_tmp;
-- Особенность конкатенации json заключается в перезаписывании.
SELECT doc || '{"id": 19}'::jsonb
FROM context_tmp;

-- Перезаписываем значение json поля.
UPDATE context_tmp
SET doc = doc || '{"id": 33}'::jsonb
WHERE (doc->'id')::INT = 3;

SELECT * FROM context_tmp;

-- 4.5. Разделить XML/JSON документ на несколько строк по узлам
CREATE TABLE IF NOT EXISTS terms_by_users
(
	doc JSON
);

INSERT INTO terms_by_users VALUES ('[{"user_id": 0, "terms_count": 100},
  								  	 {"user_id": 1, "terms_count": 0}, 
								  	 {"user_id": 2, "terms_count": 150}]');

SELECT * FROM terms_by_users;

-- jsonb_array_elements - Разворачивает массив JSON в набор значений JSON.
SELECT jsonb_array_elements(doc::jsonb) AS elems
FROM terms_by_users;

 
 
-- Защита:

--touch /home/json/test.json
--chmod 777 /home/json/test.json

COPY
(
	SELECT row_to_json(res) 
	FROM (
		  SELECT DISTINCT registration_date, text
		  FROM public.terms_en
		  WHERE registration_date BETWEEN '2021-12-01' AND '2021-12-31' 
		 ) AS res
)
TO '/home/json/test.json';
 
 
 
 