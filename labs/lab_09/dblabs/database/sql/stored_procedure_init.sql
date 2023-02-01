create or replace procedure insert_term_en
(
    id int,
    model_id int,
    class int,
    words_count int,
    registration_date date,
    text text
)
as $$
begin
    insert into public.terms_en overriding user value -- or overriding system value
    values (id, model_id, class, words_count, registration_date, text);
end;
$$ language plpgsql;
