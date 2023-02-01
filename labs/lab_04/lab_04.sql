select * from pg_language;

create extension plpython3u;

-- Создание строк в plpython3u:
-- \'{str}\' => 'str'

-- 1) Определяемая пользователем скалярная функцию CLR.
-- Получить имя пользователя по id.
CREATE OR REPLACE FUNCTION get_username(user_id INT)
RETURNS VARCHAR
AS $$
	res = plpy.execute(f" \
	    SELECT name \
	    FROM public.users  \
	    WHERE id = {user_id};")
	
	if res:
	    return res[0]['name']
$$ LANGUAGE plpython3u;

SELECT * FROM get_username(5);



-- 2) Пользовательская агрегатная функцию CLR.
-- Получить среднее количество слов в русскоязычных терминах заданного класса.
CREATE OR REPLACE FUNCTION get_avg_words_count_ru(cls INT)
RETURNS FLOAT
AS $$
	# 1-й способ:
	# res = plpy.execute(f" \
	#     SELECT AVG(words_count) \
	#     FROM public.terms_ru      \
	#     WHERE class = {cls};")
	# 
	# return res
	
	# 2-й способ:
	words_count = 0
	count = 0
	res = plpy.execute("SELECT * FROM public.terms_ru")
	
	if not res:
		return -1
	
	for elem in res:
		if elem["class"] == cls:
			words_count += elem["words_count"]
			count += 1
	
	if not count:
		return 0
	
	return words_count / count
$$ LANGUAGE plpython3u;

select get_avg_words_count_ru(3);



-- 3) Определяемая пользователем табличная функция CLR.
-- Получить все русскоязычные термины, соответствующие заданной модели.
CREATE OR REPLACE FUNCTION get_terms_ru(model_id INT)
RETURNS TABLE
(
	id INT,
	model_id INT,
	class INT,
	words_count INT,
	registration_date DATE,
	text TEXT
)
AS $$
	table = plpy.execute(f" \
		SELECT public.terms_ru.id as id, \
			model_id, \
			public.terms_ru.class as class, \
			public.terms_ru.words_count as words_count, \
			public.terms_ru.registration_date as registration_date, \
			public.terms_ru.text as text \
		FROM public.terms_ru JOIN public.models ON public.terms_ru.model_id = public.models.id;")
	# 	WHERE public.models.id = {model_id}	
	
	res = []
	
	for elem in table:
		if elem["model_id"] == model_id:
			res.append(elem)
	
	return res
$$ LANGUAGE plpython3u;

select * from get_terms_ru(19);



-- 4) Хранимая процедура CLR.
-- Добавление англоязычного термина.
CREATE OR REPLACE PROCEDURE add_term_en
(
	id INT,
	model_id INT,
	cls INT,
	words_count INT,
	registration_date DATE,
	t TEXT
)
AS $$
	# 1 способ:
	# plpy.execute(f"INSERT INTO public.terms_en VALUES({id}, {model_id}, {cls}, {words_count}, \'{registration_date}\', \'{t}\');")
	
	# 2 способ:
	# Функция plpy.prepare подготавливает шаблон запроса. Передаётся строка запроса и список типов параметров.
	plan = plpy.prepare("INSERT INTO public.terms_en VALUES($1, $2, $3, $4, $5, $6)", 
						["INT", "INT", "INT","INT", "DATE", "TEXT"])
	res = plpy.execute(plan, [id, model_id, cls, words_count, registration_date, t])
$$ LANGUAGE plpython3u;

call add_term_en(5019, 19, 1, 1, '2022-10-16', 'Andrey');



-- 5) Триггер CLR.
-- "Мягкое" удаление пользователя.
drop VIEW if exists tmp_users;
CREATE VIEW tmp_users AS
SELECT *
FROM users
WHERE id < 19;

drop FUNCTION if exists soft_user_delete();
CREATE OR REPLACE FUNCTION soft_user_delete()
RETURNS TRIGGER
AS $$
	old_id = TD["old"]["id"]
	rv = plpy.execute(f" \
	UPDATE tmp_users \
	SET name = \'deleted_user\'  \
	WHERE tmp_users.id = {old_id}")
	
	return TD["new"]
$$ LANGUAGE plpython3u;

drop trigger if exists user_delete_trigger on tmp_users;
CREATE TRIGGER user_delete_trigger
INSTEAD OF DELETE ON tmp_users
FOR EACH ROW
EXECUTE PROCEDURE soft_user_delete();



delete from tmp_users
where id = 7;

SELECT * FROM tmp_users;



-- 6) Определяемый пользователем тип данных CLR.
-- Тип содержит модель и количество соответствующих ей терминов.
CREATE TYPE model AS
(
	id INT,
	count INT
);

CREATE OR REPLACE FUNCTION get_model(m_id INT)
RETURNS model
AS
$$
	res = plpy.execute(f"      \
	SELECT public.models.id as id, COUNT(model_id) as count \
	FROM (SELECT * FROM public.terms_ru\
		  UNION \
		  SELECT * FROM public.terms_en) as terms \
		 RIGHT JOIN public.models ON terms.model_id = public.models.id              \
	WHERE public.models.id = {m_id}\
	GROUP BY public.models.id;")
	
	# nrows() возвращает количествово строк в результате запроса.
	if (res.nrows()):
		return (res[0]["id"], res[0]["count"])
$$ LANGUAGE plpython3u;

SELECT * FROM get_model('19');



-- Защита:

CREATE OR REPLACE FUNCTION level(model_id INT)
RETURNS INT
AS $$
	table = plpy.execute(f" \
		SELECT id, text \
		FROM public.models")
		
	cur_str = plpy.execute(f" \
		SELECT text \
		FROM public.models \
		WHERE id = {model_id}")[0]["text"]
			
	def get_level(str):
	    level = 0
	    
	    for elem in table:
	        if str.find(elem["text"]) != -1 and str != elem["text"]:
	            level = max(level, get_level(elem["text"]) + 1)
	    
	    return level
	
	return get_level(cur_str)
$$ LANGUAGE plpython3u;

SELECT id, text
FROM public.models;

SELECT level(122);

SELECT id, text, level(id) as level
FROM public.models;

SELECT *
FROM (
	  SELECT id, text, level(id) as level
	  FROM public.models
	 ) as leveled_models
WHERE level = 3;


