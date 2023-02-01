insert into public.terms_ru overriding user value -- or overriding system value
values (null, $1, $2, $3, $4, $5)
returning id;