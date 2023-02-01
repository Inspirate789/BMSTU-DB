delete from public.terms_ru
where public.terms_ru.id = $1
returning public.terms_ru.text;
