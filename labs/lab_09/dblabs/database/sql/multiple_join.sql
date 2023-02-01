select te.id, te.model_id, te.class, te.words_count, te.registration_date, te.text
from public.terms_en te join ( select id
                               from public.models
                               order by class desc
                               limit 1 ) as M on M.id = te.model_id
union
select tr.id, tr.model_id, tr.class, tr.words_count, tr.registration_date, tr.text
from public.terms_ru tr join ( select id
                               from public.models
                               order by words_count desc
                               limit 1 ) as M on M.id = tr.model_id;
