drop schema public cascade;
create schema public;

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

-- Структурные модели терминов
drop table if exists public.terms_ru;
drop table if exists public.terms_en;
drop table if exists public.models;
create table if not exists public.models(
	id int primary key,
	class int,
	words_count int,
	registration_date date,
	text text
);

-- Единицы русского языка
create table if not exists public.terms_ru(
	id int primary key,
	model_id int,
	foreign key (model_id) references public.models(id),
	class int,
	words_count int,
	registration_date date,
	text text
);

-- Единицы английского языка
create table if not exists public.terms_en(
	id int primary key,
	model_id int,
	foreign key (model_id) references public.models(id),
	class int,
	words_count int,
	registration_date date,
	text text
);


-- Контексты употребления
drop table if exists public.contexts;
create table if not exists public.contexts(
	id int primary key,
	terms_count int,
	words_count int,
	registration_date date,
	text text
);


-- Части речи (справочная таблица)
drop table if exists public.parts_of_speech;
create table if not exists public.parts_of_speech(
	id int primary key,
	name text,
	registration_date date
);


-- Таблица-связка (пользователи и единицы русского языка)
drop table if exists public.users_and_terms_ru;
create table if not exists public.users_and_terms_ru(
	user_id int,
	term_id int,
	registration_date date
);

-- Таблица-связка (пользователи и единицы английского языка)
drop table if exists public.users_and_terms_en;
create table if not exists public.users_and_terms_en(
	user_id int,
	term_id int,
	registration_date date
);


-- Таблица-связка (контексты употребления и единицы русского языка)
drop table if exists public.contexts_and_terms_ru;
create table if not exists public.contexts_and_terms_ru(
	context_id int,
	term_id int
);

-- Таблица-связка (контексты употребления и единицы английского языка)
drop table if exists public.contexts_and_terms_en;
create table if not exists public.contexts_and_terms_en(
	context_id int,
	term_id int
);

-- Таблица-связка (структурные модели терминов и части речи)
drop table if exists public.models_and_pos;
create table if not exists public.models_and_pos(
	model_id int,
	pos_id int
);
