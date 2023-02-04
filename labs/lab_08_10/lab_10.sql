SELECT
CURRENT_DATE as time,
count(*) as class_1
FROM terms_ru
WHERE class = 1;

SELECT
CURRENT_DATE as time,
count(*) as class_2
FROM terms_ru
WHERE class = 2;

SELECT
CURRENT_DATE as time,
count(*) as class_3
FROM terms_ru
WHERE class = 3;

SELECT registration_date, COUNT(*)
FROM users
GROUP BY registration_date;

SELECT count(*)
FROM terms_ru;

SELECT count(*)
FROM terms_en;

SELECT *
FROM parts_of_speech;
