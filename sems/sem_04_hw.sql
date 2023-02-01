drop table if exists params;
create table if not exists params(
	id int,
	p_name text, 
	p_value text
);

insert into params values(1, 'name', 'Julia');
insert into params values(1, 'gender', 'F');
insert into params values(2, 'name', 'Ivan');

select * from params;

select n.id, n.p_value as name, g.p_value as gender
from
		(select id, p_value
		from params
		where p_name = 'name') as n
	full join
		(select id, p_value
		from params
		where p_name = 'gender') as g
	on n.id = g.id
;

-- Или так:

with names(id, name) as (
	select id, p_value
    from params
    where p_name = 'name'
),
gender(id, gender) as (
 	select id, p_value
    from params
    where p_name = 'gender'
)
select n.id, 
       n.name as name, 
       g.gender as gender
from names as n full join gender as g 
     on n.id = g.id
