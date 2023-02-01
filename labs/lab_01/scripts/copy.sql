-- Пользователи
copy public.users(id, name, registration_date, first_addition_date, last_addition_date, terms_count) from '/home/data/users.csv' delimiter ',';
select * from public.users;

-- Структурные модели терминов
copy public.models(id, class, words_count, registration_date, text) from '/home/data/models.csv' delimiter ',';
select * from public.models;

alter table public.models add column description text;

update public.models set description = 'Test test test' where words_count >= 6;

alter table public.models drop column description;

-- Единицы русского языка
copy public.terms_ru(id, model_id, class, words_count, registration_date, text) from '/home/data/terms_ru.csv' delimiter ',';
select * from public.terms_ru;

-- Единицы английского языка
copy public.terms_en(id, model_id, class, words_count, registration_date, text) from '/home/data/terms_en.csv' delimiter ',';
select * from public.terms_en;

-- Контексты употребления
copy public.contexts(id, terms_count, words_count, registration_date, text) from '/home/data/contexts.csv' delimiter ',';
select * from public.contexts;

-- Части речи (справочная таблица)
copy public.parts_of_speech(id, name, registration_date) from '/home/data/pos.csv' delimiter ',';
select * from public.parts_of_speech;

-- Таблица-связка (пользователи и единицы русского языка)
copy public.users_and_terms_ru(user_id, term_id, registration_date) from '/home/data/users_and_terms_ru.csv' delimiter ',';
select * from public.users_and_terms_ru;

-- Таблица-связка (пользователи и единицы английского языка)
copy public.users_and_terms_en(user_id, term_id, registration_date) from '/home/data/users_and_terms_en.csv' delimiter ',';
select * from public.users_and_terms_en;

-- Таблица-связка (контексты употребления и единицы русского языка)
copy public.contexts_and_terms_ru(context_id, term_id) from '/home/data/contexts_and_terms_ru.csv' delimiter ',';
select * from public.contexts_and_terms_ru;

-- Таблица-связка (контексты употребления и единицы английского языка)
copy public.contexts_and_terms_en(context_id, term_id) from '/home/data/contexts_and_terms_en.csv' delimiter ',';
select * from public.contexts_and_terms_en;

-- Таблица-связка (структурные модели терминов и части речи)
copy public.models_and_pos(model_id, pos_id) from '/home/data/models_and_pos.csv' delimiter ',';
select * from public.models_and_pos;
