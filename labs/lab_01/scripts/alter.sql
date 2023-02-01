-- Пользователи
alter table public.users add constraint reg_date_constr check (first_addition_date >= registration_date);
alter table public.users add constraint add_date_constr check (last_addition_date >= first_addition_date);
alter table public.users add constraint terms_count_constr check (terms_count >= 0);

-- Структурные модели терминов
alter table public.models add constraint model_class_constr check (class >= 1 and class <= 3);
alter table public.models add constraint words_count_constr check (words_count >= 0);

-- insert into public.models(id, class, words_count, registration_date, text) values(1001, 3, 10, '2022-03-05', 'alter table public models add constraint'); -- Correct
-- insert into public.models(id, class, words_count, registration_date, text) values(1002, 3, -3, '2022-03-05', 'alter table public models add constraint'); -- Incorrect

-- Единицы русского языка
alter table public.terms_ru add constraint term_class_constr check (class >= 1 and class <= 3);
alter table public.terms_ru add constraint words_count_constr check (words_count >= 0);

-- Единицы английского языка
alter table public.terms_en add constraint term_class_constr check (class >= 1 and class <= 3);
alter table public.terms_en add constraint words_count_constr check (words_count >= 0);

-- Контексты употребления
alter table public.contexts add constraint terms_count_constr check (terms_count >= 0);
alter table public.contexts add constraint words_count_constr check (words_count >= 0);
