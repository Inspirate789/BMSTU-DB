select avg(EXTRACT(YEAR FROM CURRENT_DATE) - EXTRACT(YEAR FROM date_of_birth))
from employees join
     (
         select distinct on (employee_id, date) employee_id, date, sum(tmp_dur) over (partition by employee_id, date) as day_dur
         from
             (
                 select employee_id, date, time,
                        type,
                        lag(time) over (partition by employee_id, date order by time) as prev_time,
                        time-lag(time) over (partition by employee_id, date order by time) as tmp_dur
                 from times t
                 order by employee_id, date, time
             ) as small_durations
     ) as day_durations
     on employees.id = day_durations.employee_id
where day_durations.day_dur < '08:00:00';
