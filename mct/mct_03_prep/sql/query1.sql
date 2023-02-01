SELECT date_part, COUNT(employee_id) AS cnt, department
FROM
    (
        select employee_id, EXTRACT(WEEK FROM date) AS date_part, COUNT(*) as late_cnt
        from times t
        where type = 1
        group by employee_id, date_part
        having min(time) > '9:00'
    ) AS tmp join employees e on tmp.employee_id = e.id
where late_cnt > 3
GROUP BY date_part, department;
