-- Пользователи
drop table if exists public.users;
create table if not exists public.users(
	id int primary key,
	name text,
	registration_date date,
	first_addition_date date,
	last_addition_date date,
	terms_count int default 0
);

alter table public.users add constraint reg_date_constr check (first_addition_date >= registration_date);
alter table public.users add constraint add_date_constr check (last_addition_date >= first_addition_date);
alter table public.users add constraint terms_count_constr check (terms_count >= 0);

copy public.users(id, name, registration_date, first_addition_date, last_addition_date, terms_count) from '/home/data/users.csv' delimiter ',';
select * from public.users;

-- Единицы русского языка
drop table if exists public.terms_ru;
create table if not exists public.terms_ru(
	id int primary key,
	model_id int,
	class int,
	words_count int,
	registration_date date,
	text text
);

alter table public.terms_ru add constraint term_class_constr check (class >= 1 and class <= 3);
alter table public.terms_ru add constraint words_count_constr check (words_count >= 0);

copy public.terms_ru(id, model_id, class, words_count, registration_date, text) from '/home/data/terms_ru.csv' delimiter ',';
select * from public.terms_ru;

-- Единицы английского языка
drop table if exists public.terms_en;
create table if not exists public.terms_en(
	id int primary key,
	model_id int,
	class int,
	words_count int,
	registration_date date,
	text text
);

alter table public.terms_en add constraint term_class_constr check (class >= 1 and class <= 3);
alter table public.terms_en add constraint words_count_constr check (words_count >= 0);

copy public.terms_en(id, model_id, class, words_count, registration_date, text) from '/home/data/terms_en.csv' delimiter ',';
select * from public.terms_en;

-- Структурные модели терминов
drop table if exists public.models;
create table if not exists public.models(
	id int primary key,
	class int,
	words_count int,
	registration_date date,
	text text
);

alter table public.models add constraint model_class_constr check (class >= 1 and class <= 3);
alter table public.models add constraint words_count_constr check (words_count >= 0);

copy public.models(id, class, words_count, registration_date, text) from '/home/data/models.csv' delimiter ',';
select * from public.models;


-- Контексты употребления
drop table if exists public.contexts;
create table if not exists public.contexts(
	id int primary key,
	terms_count int,
	words_count int,
	registration_date date,
	text text
);

alter table public.contexts add constraint terms_count_constr check (terms_count >= 0);
alter table public.contexts add constraint words_count_constr check (words_count >= 0);

copy public.contexts(id, terms_count, words_count, registration_date, text) from '/home/data/contexts.csv' delimiter ',';
select * from public.contexts;


-- Части речи (справочная таблица)
drop table if exists public.parts_of_speech;
create table if not exists public.parts_of_speech(
	id int primary key,
	name text,
	registration_date date
);

copy public.parts_of_speech(id, name, registration_date) from '/home/data/pos.csv' delimiter ',';
select * from public.parts_of_speech;


-- Таблица-связка (пользователи и единицы русского языка)
drop table if exists public.users_and_terms_ru;
create table if not exists public.users_and_terms_ru(
	user_id int,
	term_id int,
	registration_date date
);

copy public.users_and_terms_ru(user_id, term_id, registration_date) from '/home/data/users_and_terms_ru.csv' delimiter ',';
select * from public.users_and_terms_ru;

-- Таблица-связка (пользователи и единицы английского языка)
drop table if exists public.users_and_terms_en;
create table if not exists public.users_and_terms_en(
	user_id int,
	term_id int,
	registration_date date
);

copy public.users_and_terms_en(user_id, term_id, registration_date) from '/home/data/users_and_terms_en.csv' delimiter ',';
select * from public.users_and_terms_en;


-- Таблица-связка (контексты употребления и единицы русского языка)
drop table if exists public.contexts_and_terms_ru;
create table if not exists public.contexts_and_terms_ru(
	context_id int,
	term_id int
);

copy public.contexts_and_terms_ru(context_id, term_id) from '/home/data/contexts_and_terms_ru.csv' delimiter ',';
select * from public.contexts_and_terms_ru;

-- Таблица-связка (контексты употребления и единицы английского языка)
drop table if exists public.contexts_and_terms_en;
create table if not exists public.contexts_and_terms_en(
	context_id int,
	term_id int
);

copy public.contexts_and_terms_en(context_id, term_id) from '/home/data/contexts_and_terms_en.csv' delimiter ',';
select * from public.contexts_and_terms_en;

-- Таблица-связка (структурные модели терминов и части речи)
drop table if exists public.models_and_pos;
create table if not exists public.models_and_pos(
	model_id int,
	pos_id int
);

copy public.models_and_pos(model_id, pos_id) from '/home/data/models_and_pos.csv' delimiter ',';
select * from public.models_and_pos;
