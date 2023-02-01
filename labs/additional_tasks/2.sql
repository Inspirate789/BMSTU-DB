DROP TABLE IF EXISTS Employee;
CREATE TABLE IF NOT EXISTS Employee(
    id INT,
    FIO TEXT,
    date_of_status DATE,
    status VARCHAR
);

INSERT INTO Employee(id, FIO, date_of_status, status)
VALUES
(1, 'Иванов Иван Иванович', '2022-12-12', 'Работа offline'),
(1, 'Иванов Иван Иванович', '2022-12-13', 'Работа offline'),
(1, 'Иванов Иван Иванович', '2022-12-14', 'Больничный'),
(1, 'Иванов Иван Иванович', '2022-12-15', 'Больничный'),
(1, 'Иванов Иван Иванович', '2022-12-16', 'Удалённая работа'),
(1, 'Иванов Иван Иванович', '2022-12-19', 'Больничный'),
(1, 'Иванов Иван Иванович', '2022-12-20', 'Больничный'),
(2, 'Петров Пётр Петрович', '2022-12-12', 'Работа offline'),
(2, 'Петров Пётр Петрович', '2022-12-13', 'Работа offline'),
(2, 'Петров Пётр Петрович', '2022-12-14', 'Удалённая работа'),
(2, 'Петров Пётр Петрович', '2022-12-15', 'Удалённая работа'),
(2, 'Петров Пётр Петрович', '2022-12-16', 'Работа offline');

SELECT * FROM Employee;

WITH numbered AS (
    SELECT ROW_NUMBER() OVER(
        PARTITION BY FIO, status
        ORDER BY date_of_status
    ) AS i, FIO, status, date_of_status
    FROM Employee
)
SELECT FIO, status, MIN(date_of_status) as date_from, MAX(date_of_status) as date_to
FROM numbered n
GROUP BY FIO, status, date_of_status - MAKE_INTERVAL(days => CAST(i AS INT)) -- (days => i::int)
ORDER BY FIO, date_from;


