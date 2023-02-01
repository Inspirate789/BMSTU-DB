SELECT e.department, COUNT(distinct e.id)
FROM employees e
         INNER JOIN times t ON (t.employee_id = e.id)
WHERE ((t.time > '09:00:00') AND (t.type = 1))
GROUP by e.department;
