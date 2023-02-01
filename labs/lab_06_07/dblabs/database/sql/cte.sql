with CTE (id, words_count) AS (
    select id, words_count
    from public.models
    where class = $1
)
select AVG(words_count)
from CTE;
