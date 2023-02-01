drop table if exists public.dinos;
create table public.dinos(
	name text primary key,
	type text,
	period text,
	threat_level int default 0,
	speed  int default 0,
	protection int default 0,
	weight numeric(15,2),
	length numeric(15,2)
);

select * from public.dinos;

-- Добавление атрибута
ALTER TABLE public.dinos ADD COLUMN info text;

UPDATE public.dinos
SET info = 'Test test test'
where length > 10;

ALTER TABLE public.dinos DROP COLUMN info;

-- Ограничения
ALTER TABLE public.dinos ADD CONSTRAINT w_constr CHECK (weight > 0);

insert into public.dinos(name, type, weight, length) values('Диплодок', 'травоядный', -1, 35);
insert into public.dinos(name, type, weight, length) values('Диплодок', 'травоядный', 80000, 35);

-- Ключи
create table public.type(
	id int,
	name text
);
insert into public.type values(1, 'плотоядный');
insert into public.type values(2, 'травоядный');

create table public.period(
	id int,
	name text
);
insert into public.period values(1, 'меловой');
insert into public.period values(2, 'юрский');

ALTER TABLE public.dinos DROP COLUMN type;
ALTER TABLE public.dinos DROP COLUMN period;

ALTER TABLE public.dinos ADD COLUMN "type" int;
ALTER TABLE public.dinos ADD COLUMN "period" int;

UPDATE public.dinos
SET type = 1
where name in ('Тиранозавр', 'Велоцираптор', 'Аллозавр', 'Спинозавр', 'Птеранодон');

UPDATE public.dinos
SET type = 2
where name in ('Трицератопс','Брахиозавр','Анкилозавр','Стегозавр','Игуанодон','Пахицефалозавр','Паразаурофол','Диплодок');

alter table public.type add constraint pk_type primary key (id)
alter table public.dinos add constraint fk_type foreign key (type) references public.type(id)

insert into public.dinos(name, type, weight, length) values('Диплодок_2', 2, 1, 35);
