-- 1. Инструкция SELECT, использующая предикат сравнения.
SELECT DISTINCT ru.text, ru.words_count, en.text, en.words_count
FROM public.terms_ru ru JOIN public.terms_en AS en ON ru.model_id = en.model_id
WHERE ru.class = en.class
	  AND ru.class = 2
ORDER BY ru.words_count; 

-- 2. Инструкция SELECT, использующая предикат BETWEEN.
SELECT DISTINCT registration_date, text
FROM public.terms_en
WHERE registration_date BETWEEN '2021-12-01' AND '2021-12-31'; 

-- 3. Инструкция SELECT, использующая предикат LIKE.
SELECT DISTINCT public.terms_ru.text, public.models.text
FROM public.terms_ru JOIN public.models ON public.models.id = public.terms_ru.model_id
WHERE public.models.text LIKE '%e%';

-- 4. Инструкция SELECT, использующая предикат IN с вложенным подзапросом.
SELECT id, model_id, text
FROM public.terms_en
WHERE model_id IN (SELECT id
					 FROM public.models 
					 WHERE words_count = 5)
	  AND class = 2;

-- 5. Инструкция SELECT, использующая предикат EXISTS с вложенным подзапросом.
SELECT id, registration_date, text
FROM public.terms_ru
WHERE EXISTS (SELECT public.terms_ru.id 
			  FROM public.terms_ru LEFT OUTER JOIN public.models
			  	   ON public.terms_ru.model_id = public.models.id
			  WHERE public.models.words_count = 10
); 

-- 6. Инструкция SELECT, использующая предикат сравнения с квантором.
SELECT id, words_count, registration_date, text
FROM public.terms_en
WHERE registration_date  > ALL ( SELECT registration_date
						         FROM public.terms_en
						         WHERE words_count = 7 );

-- 7. Инструкция SELECT, использующая агрегатные функции в выражениях столбцов.
SELECT AVG(words_count) AS Actual_AVG,
 	   SUM(words_count) / COUNT(id) AS Calc_AVG
FROM ( SELECT id, words_count
 	   FROM public.terms_ru
 	   GROUP BY id
) AS total_terms_ru;

-- 8. Инструкция SELECT, использующая скалярные подзапросы в выражениях столбцов.
SELECT id, words_count,
	   ( SELECT AVG(words_count)
	 	 FROM public.terms_en
	 	 WHERE public.terms_en.id = public.models.id ) AS avg_words_count, -- WHERE 0 < public.models.id ) AS avg_words_count,
	   ( SELECT MAX(words_count)
	 	 FROM public.terms_en
	 	 WHERE public.terms_en.id = public.models.id ) AS max_words_count, -- WHERE 0 < public.models.id ) AS max_words_count,
	   text
FROM public.models
WHERE class = 3;

-- 9. Инструкция SELECT, использующая простое выражение CASE.
SELECT public.terms_ru.id, public.terms_ru.text,
 	   CASE DATE_PART('year', public.terms_ru.registration_date)
		    WHEN DATE_PART('year', NOW()) THEN 'This Year'
		    WHEN DATE_PART('year', NOW()) - 1 THEN 'Last year'
		    ELSE CONCAT(CAST(DATE_PART('year', NOW()) - DATE_PART('year', public.terms_ru.registration_date) AS varchar(5)), ' years ago')
 	   END AS Registration
FROM public.terms_ru JOIN public.models ON public.terms_ru.model_id = public.models.id;

-- 10. Инструкция SELECT, использующая поисковое выражение CASE.
SELECT id, text,
 	   CASE
		   WHEN words_count < 3 THEN 'Short'
		   WHEN words_count < 5 THEN 'Normal'
		   WHEN words_count < 8 THEN 'Long'
		   ELSE 'Very Long'
 	   END AS length 
FROM public.terms_en;

-- 11. Создание новой временной локальной таблицы из результирующего набора данных инструкции SELECT.
DROP TABLE IF EXISTS public.longest_models;

SELECT id, words_count, text
INTO public.longest_models
FROM public.models
WHERE class > 1
GROUP BY id, words_count;

SELECT * from public.longest_models;

-- 12. Инструкция SELECT, использующая вложенные коррелированные подзапросы в качестве производных таблиц в предложении FROM.
SELECT public.terms_en.id, public.terms_en.class, public.terms_en.words_count, text as longest
FROM public.terms_en JOIN ( SELECT id, class, words_count
					   		FROM public.models m 
					   		ORDER BY class DESC
					   		LIMIT 1 ) AS M ON M.id = public.terms_en.model_id
UNION
SELECT public.terms_en.id, public.terms_en.class, public.terms_en.words_count, text as longest
FROM public.terms_en JOIN ( SELECT id, class, words_count
					   		FROM public.models m 
					   		ORDER BY words_count DESC
					   		LIMIT 1 ) AS M ON M.id = public.terms_en.model_id;

-- 13. Инструкция SELECT, использующая вложенные подзапросы с уровнем вложенности 3.
SELECT id, words_count, text
FROM public.terms_ru
WHERE id = ( SELECT id
			 FROM public.models
			 GROUP BY id
			 HAVING SUM(words_count) = ( SELECT MAX(SW)
										 FROM ( SELECT SUM(words_count) as SW
											    FROM public.models
											    GROUP BY id
										 	  ) AS models_wc_sum
			 						   )
			 LIMIT 1
		   );

-- 14. Инструкция SELECT, консолидирующая данные с помощью предложения GROUP BY, но без предложения HAVING.
SELECT te.id, te.words_count, te.text, m.words_count AS model_words_count
FROM public.terms_en te LEFT OUTER JOIN public.models m ON te.model_id = m.id
WHERE te.class = 2
GROUP BY te.id, te.words_count, m.words_count;

-- 15. Инструкция SELECT, консолидирующая данные с помощью предложения GROUP BY и предложения HAVING.
SELECT id, words_count
FROM public.terms_ru
GROUP BY id
HAVING words_count > ( SELECT AVG(words_count)
 					   FROM public.terms_ru);

-- 16. Однострочная инструкция INSERT, выполняющая вставку в таблицу одной строки значений.
INSERT INTO public.models(id, class, words_count, registration_date, text)
VALUES (1001, 3, 10, '2022-03-05', 'alter table public models add constraint')
RETURNING *;

-- 17. Многострочная инструкция INSERT, выполняющая вставку в таблицу результирующего набора данных вложенного подзапроса.
INSERT INTO public.contexts(id, terms_count, words_count, registration_date, text)
SELECT id, 1 AS terms_count, words_count, NOW() AS registration_date, text
FROM public.terms_ru
WHERE id > ( SELECT MAX(id)
 		 	 FROM public.contexts )
RETURNING *;

-- 18. Простая инструкция UPDATE.
UPDATE public.terms_en
SET words_count = words_count * 100
WHERE id = 5
RETURNING *;

-- 19. Инструкция UPDATE со скалярным подзапросом в предложении SET.
UPDATE public.models
SET words_count = ( SELECT AVG(words_count)
 				  	FROM public.terms_ru
 				  	WHERE model_id = 10 )
WHERE id = 10
RETURNING *;

-- 20. Простая инструкция DELETE.
DELETE FROM public.terms_en
WHERE words_count > 3
RETURNING *;

-- 21. Инструкция DELETE с вложенным коррелированным подзапросом в предложении WHERE.
DELETE FROM public.terms_ru
WHERE id IN ( SELECT public.terms_ru.id
			  FROM public.terms_ru LEFT OUTER JOIN public.models
				   ON public.terms_ru.model_id = public.models.id
			  WHERE public.models.registration_date > '2020-09-01')
RETURNING *;

-- 22. Инструкция SELECT, использующая простое обобщенное табличное выражение
WITH CTE (id, words_count) AS (
	 SELECT id, words_count
	 FROM public.models 
	 WHERE class = 2
	 GROUP BY id
)
SELECT AVG(words_count) AS "Среднее количество слов в моделях терминов"
FROM CTE;

-- 23. Инструкция SELECT, использующая рекурсивное обобщенное табличное выражение.
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
ORDER BY level DESC;



-- 24. Оконные функции. Использование конструкций MIN/MAX/AVG OVER()
SELECT te.id, te.class, te.text, m.words_count,
	   AVG(m.words_count) OVER(PARTITION BY te.class) AS avg_words_count,
	   MIN(m.words_count) OVER(PARTITION BY te.class) AS min_words_count,
	   MAX(m.words_count) OVER(PARTITION BY te.class) AS max_words_count
FROM public.terms_en te LEFT OUTER JOIN public.models m 
	 ON te.model_id = m.id
ORDER BY te.class DESC;

-- 25. Оконные фнкции для устранения дублей
-- Придумать запрос, в результате которого в данных появляются полные дубли.
-- Устранить дублирующиеся строки с использованием функции ROW_NUMBER().
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
WHERE row_number = 1;

-- Защита
SELECT COUNT(id) AS count_ru, 0 AS count_en
	FROM ( SELECT public.terms_ru.id, public.terms_ru.text
		   FROM public.terms_ru JOIN public.models ON public.models.id = public.terms_ru.model_id
		   WHERE public.models.text LIKE '%a%'
	) AS terms_ru_with_nouns 
UNION 
SELECT 0 as count_ru, COUNT(id) AS count_en
	FROM ( SELECT public.terms_en.id, public.terms_en.text
		   FROM public.terms_en JOIN public.models ON public.models.id = public.terms_en.model_id
		   WHERE public.models.text LIKE '%a%'
) AS terms_en_with_nouns; 
				   		
				   		
					   		
					   		
					   		
