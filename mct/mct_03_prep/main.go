package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

//go:generate go-bindata -pkg main -o sqlscripts.go ./sql

type AuthDB struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func recreate(sqlDB *sql.DB) error {
	script, err := Asset("sql/task1.sql")
	if err != nil {
		return err
	}

	_, err = sqlDB.Exec(string(script))

	return nil
}

func SQLquery1(sqlDB *sql.DB) error {
	script, err := Asset("sql/query1.sql")
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(string(script))
	if err != nil {
		return err
	}
	defer rows.Close()

	var date, dep string
	var count int

	for rows.Next() {
		err = rows.Scan(&date, &count, &dep)
		if err != nil {
			return err
		}
		fmt.Println(date, count, dep)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func SQLquery2(sqlDB *sql.DB) error {
	script, err := Asset("sql/query2.sql")
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(string(script))
	if err != nil {
		return err
	}
	defer rows.Close()

	var avg float32

	rows.Next()
	err = rows.Scan(&avg)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	fmt.Println(avg)

	return nil
}

func SQLquery3(sqlDB *sql.DB) error {
	script, err := Asset("sql/query3.sql")
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(string(script))
	if err != nil {
		return err
	}
	defer rows.Close()

	var dep string
	var count int

	for rows.Next() {
		err = rows.Scan(&dep, &count)
		if err != nil {
			return err
		}
		fmt.Println(dep, count)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func task1(sqlDB *sql.DB) error {
	err := recreate(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти отделы, в которых хоть один сотрудник опаздывает больше 3-х раз в неделю.")
	err = SQLquery1(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти средний возраст сотрудников, не находящихся на рабочем месте 8 часов в день.")
	err = SQLquery2(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Вывести все отделы и количество сотрудников хоть раз опоздавших за всю историю учета.")

	return SQLquery3(sqlDB)
}

func GORMquery1(gormDB *gorm.DB) error {
	subQuery := gormDB.Select("employee_id, EXTRACT(WEEK FROM date) AS date_part, COUNT(*) as late_cnt").Table("times t").Where("type = 1").Group("employee_id, date_part").Having("min(time) > '9:00'")
	rows, err := gormDB.Select("date_part, COUNT(employee_id) AS cnt, department").Table("(?) AS tmp", subQuery).Joins("join employees e on tmp.employee_id = e.id").Where("late_cnt > 3").Group("date_part, department").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var date, dep string
	var count int

	for rows.Next() {
		err = rows.Scan(&date, &count, &dep)
		if err != nil {
			return err
		}
		fmt.Println(date, count, dep)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func GORMquery2(gormDB *gorm.DB) error {
	subQuery1 := gormDB.Select("employee_id, date, time," +
		"type," +
		"lag(time) over (partition by employee_id, date order by time) as prev_time," +
		"time-lag(time) over (partition by employee_id, date order by time) as tmp_dur").Table("times t").Order("employee_id, date, time")
	subQuery2 := gormDB.Select("distinct on (employee_id, date) employee_id, date, sum(tmp_dur) over (partition by employee_id, date) as day_dur").Table("(?) as small_durations", subQuery1)
	rows, err := gormDB.Select("avg(EXTRACT(YEAR FROM CURRENT_DATE) - EXTRACT(YEAR FROM date_of_birth))").Table("employees").Joins("join (?) as day_durations on employees.id = day_durations.employee_id", subQuery2).Where("day_durations.day_dur < '08:00:00'").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var avg float32

	rows.Next()
	err = rows.Scan(&avg)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	fmt.Println(avg)

	return nil
}

func GORMquery3(gormDB *gorm.DB) error {
	rows, err := gormDB.Select("e.department, COUNT(distinct e.id)").Table("employees AS e").Joins("inner join times AS t on t.employee_id = e.id").Where("t.time > '09:00:00' AND t.type = 1").Group("e.department").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var dep string
	var count int

	for rows.Next() {
		err = rows.Scan(&dep, &count)
		if err != nil {
			return err
		}
		fmt.Println(dep, count)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func task2(gormDB *gorm.DB) error {
	fmt.Println("Найти отделы, в которых хоть один сотрудник опаздывает больше 3-х раз в неделю.")
	err := GORMquery1(gormDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти средний возраст сотрудников, не находящихся на рабочем месте 8 часов в день.")
	err = GORMquery2(gormDB)
	if err != nil {
		return err
	}

	fmt.Println("Вывести все отделы и количество сотрудников хоть раз опоздавших за всю историю учета.")

	return GORMquery3(gormDB)
}

func main() {
	data, err := os.ReadFile("AuthDB.json")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var dbRequest AuthDB
	err = json.Unmarshal(data, &dbRequest)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbRequest.Host, dbRequest.Port, dbRequest.Username, dbRequest.Password, dbRequest.DBName, dbRequest.SSLMode)
	sqlDB, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================================SQL==========================================")
	err = task1(sqlDB)
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================================GORM=========================================")
	err = task2(gormDB)
	if err != nil {
		panic(err)
	}

	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}
