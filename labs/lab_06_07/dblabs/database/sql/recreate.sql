drop schema public cascade;
create schema public;

-- Пользователи
-- drop table if exists public.users;
create table if not exists public.users(
    id int generated always as identity primary key,
    name text not null,
    registration_date date not null check (first_addition_date >= registration_date),
    first_addition_date date not null check (last_addition_date >= first_addition_date),
    last_addition_date date not null check (terms_count >= 0),
    terms_count int default 0
);



-- Структурные модели терминов
create table if not exists public.models(
    id int generated always as identity primary key,
    class int not null check (class >= 1 and class <= 3),
    words_count int not null check (words_count >= 0),
    registration_date date not null,
    text text not null
);



-- Единицы русского языка
create table if not exists public.terms_ru(
    id int generated always as identity primary key,
    model_id int,
    foreign key (model_id) references public.models(id),
    class int not null check (class >= 1 and class <= 3),
    words_count int not null check (words_count >= 0),
    registration_date date not null,
    text text not null
);



-- Единицы английского языка
create table if not exists public.terms_en(
    id int generated always as identity primary key,
    model_id int,
    foreign key (model_id) references public.models(id),
    class int not null check (class >= 1 and class <= 3),
    words_count int not null check (words_count >= 0),
    registration_date date not null,
    text text not null
);



-- Контексты употребления
-- drop table if exists public.contexts;
create table if not exists public.contexts(
    id int generated always as identity primary key,
    terms_count int not null check (terms_count >= 0),
    words_count int not null check (terms_count >= 0),
    registration_date date not null,
    text text not null
);



-- Части речи (справочная таблица)
-- drop table if exists public.parts_of_speech;
create table if not exists public.parts_of_speech(
    id int generated always as identity primary key,
    name text not null,
    registration_date date not null
);



-- Таблица-связка (пользователи и единицы русского языка)
-- drop table if exists public.users_and_terms_ru;
create table if not exists public.users_and_terms_ru(
    user_id int,
    foreign key (user_id) references public.users(id),
    term_id int,
    foreign key (term_id) references public.terms_ru(id),
    registration_date date
);



-- Таблица-связка (пользователи и единицы английского языка)
-- drop table if exists public.users_and_terms_en;
create table if not exists public.users_and_terms_en(
    user_id int,
    foreign key (user_id) references public.users(id),
    term_id int,
    foreign key (term_id) references public.terms_en(id),
    registration_date date
);



-- Таблица-связка (контексты употребления и единицы русского языка)
-- drop table if exists public.contexts_and_terms_ru;
create table if not exists public.contexts_and_terms_ru(
    context_id int,
    foreign key (context_id) references public.contexts(id),
    term_id int,
    foreign key (term_id) references public.terms_ru(id)
);



-- Таблица-связка (контексты употребления и единицы английского языка)
-- drop table if exists public.contexts_and_terms_en;
create table if not exists public.contexts_and_terms_en(
    context_id int,
    foreign key (context_id) references public.contexts(id),
    term_id int,
    foreign key (term_id) references public.terms_en(id)
);



-- Таблица-связка (структурные модели терминов и части речи)
-- drop table if exists public.models_and_pos;
create table if not exists public.models_and_pos(
    model_id int,
    foreign key (model_id) references public.models(id),
    pos_id int,
    foreign key (pos_id) references public.parts_of_speech(id)
);
