copy public.users(name, registration_date, first_addition_date, last_addition_date, terms_count) from '/home/dblabs/generator/data/users.csv' delimiter ',';
-- select * from public.users;

copy public.models(class, words_count, registration_date, text) from '/home/dblabs/generator/data/models.csv' delimiter ',';
-- select * from public.models;

copy public.terms_ru(model_id, class, words_count, registration_date, text) from '/home/dblabs/generator/data/terms_ru.csv' delimiter ',';
-- select * from public.terms_ru;

copy public.terms_en(model_id, class, words_count, registration_date, text) from '/home/dblabs/generator/data/terms_en.csv' delimiter ',';
-- select * from public.terms_en;

copy public.contexts(terms_count, words_count, registration_date, text) from '/home/dblabs/generator/data/contexts.csv' delimiter ',';
-- select * from public.contexts;

copy public.parts_of_speech(name, registration_date) from '/home/dblabs/generator/data/pos.csv' delimiter ',';
-- select * from public.parts_of_speech;

copy public.users_and_terms_ru(user_id, term_id, registration_date) from '/home/dblabs/generator/data/users_terms_ru.csv' delimiter ',';
-- select * from public.users_and_terms_ru;

copy public.users_and_terms_en(user_id, term_id, registration_date) from '/home/dblabs/generator/data/users_terms_en.csv' delimiter ',';
-- select * from public.users_and_terms_en;

copy public.contexts_and_terms_ru(context_id, term_id) from '/home/dblabs/generator/data/contexts_terms_ru.csv' delimiter ',';
-- select * from public.contexts_and_terms_ru;

copy public.contexts_and_terms_en(context_id, term_id) from '/home/dblabs/generator/data/contexts_terms_en.csv' delimiter ',';
-- select * from public.contexts_and_terms_en;

copy public.models_and_pos(model_id, pos_id) from '/home/dblabs/generator/data/models_pos.csv' delimiter ',';
-- select * from public.models_and_pos;
