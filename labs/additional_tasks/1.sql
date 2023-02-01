drop schema if exists lab2 cascade;
create schema lab2;

create table if not exists lab2.table1(
	id int,
	var1 text,
	valid_from_dttm date,
	valid_to_dttm date
);

create table if not exists lab2.table2(
	id int,
	var2 text,
	valid_from_dttm date,
	valid_to_dttm date
);

insert into lab2.table1(id, var1, valid_from_dttm, valid_to_dttm) values(1, 'A', '2018-09-01', '2018-09-15');
insert into lab2.table1(id, var1, valid_from_dttm, valid_to_dttm) values(1, 'B', '2018-09-16', '5999-12-31');
-- insert into lab2.table1(id, var1, valid_from_dttm, valid_to_dttm) values(2, 'A', '2018-09-01', '2018-09-15');
-- insert into lab2.table1(id, var1, valid_from_dttm, valid_to_dttm) values(2, 'B', '2018-09-16', '5999-12-31');

insert into lab2.table2(id, var2, valid_from_dttm, valid_to_dttm) values(1, 'A', '2018-09-01', '2018-09-18');
insert into lab2.table2(id, var2, valid_from_dttm, valid_to_dttm) values(1, 'B', '2018-09-19', '5999-12-31');
-- insert into lab2.table2(id, var2, valid_from_dttm, valid_to_dttm) values(2, 'A', '2018-09-01', '2018-09-18');
-- insert into lab2.table2(id, var2, valid_from_dttm, valid_to_dttm) values(2, 'B', '2018-09-19', '5999-12-31');

select * from lab2.table1;
select * from lab2.table2;



-- explain
select id, var1, var2, valid_from_dttm, valid_to_dttm
from (
	select id, var1, var2, valid_from_dttm, valid_to_dttm, ROW_NUMBER() OVER(PARTITION BY id, valid_from_dttm)
	from (
		select id, var1, var2, valid_from_dttm, valid_to_dttm
		from (
			(select lab2.table1.id, var1, var1 as var2, lab2.table1.valid_from_dttm, lab2.table1.valid_to_dttm
			 from lab2.table1 join lab2.table2 on lab2.table1.id = lab2.table2.id)
			union
			(select lab2.table1.id, var1, var2, lab2.table1.valid_from_dttm, lab2.table2.valid_to_dttm
			 from lab2.table1 join lab2.table2 on lab2.table1.id = lab2.table2.id)
			union
			(select lab2.table1.id, var1, var2, lab2.table2.valid_from_dttm, lab2.table1.valid_to_dttm
			 from lab2.table1 join lab2.table2 on lab2.table1.id = lab2.table2.id)
			union
			(select lab2.table1.id, var2 as var1, var2, lab2.table2.valid_from_dttm, lab2.table2.valid_to_dttm
			 from lab2.table1 join lab2.table2 on lab2.table1.id = lab2.table2.id)
		) as combinations -- В этой таблице нет строк с одинаковыми датами и разными комбинациями var1 и var2
		where valid_from_dttm < valid_to_dttm
		order by id, valid_from_dttm asc, valid_to_dttm asc
	) as numbered_combinations
) as linked_combinations
where row_number = 1;



-- Более оптимальный вариант:
-- explain
select lab2.table1.id, lab2.table1.var1, lab2.table2.var2,
	case when lab2.table1.valid_from_dttm <= lab2.table2.valid_from_dttm 
				then lab2.table2.valid_from_dttm
				else lab2.table1.valid_from_dttm
	end valid_from_dttm,
	case when lab2.table1.valid_to_dttm >= lab2.table2.valid_to_dttm 
				then lab2.table2.valid_to_dttm
				else lab2.table1.valid_to_dttm
	end valid_to_dttm
from lab2.table1 full outer join lab2.table2 on lab2.table1.id = lab2.table2.id
		and lab2.table1.valid_to_dttm >= lab2.table2.valid_from_dttm
		and lab2.table2.valid_to_dttm >= lab2.table1.valid_from_dttm
order by id, valid_from_dttm;




