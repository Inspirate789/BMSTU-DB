select avg(words_count)
from (select * from public.terms_ru union select * from public.terms_en) as terms_all
group by class
order by class asc;
