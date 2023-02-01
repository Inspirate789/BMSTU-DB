package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	var dep string

	for rows.Next() {
		err = rows.Scan(&dep)
		if err != nil {
			return err
		}
		fmt.Println(dep)
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

	var id int

	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

func SQLquery3(sqlDB *sql.DB, date string) error {
	script, err := Asset("sql/query3.sql")
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(string(script), date)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dep string

	for rows.Next() {
		err = rows.Scan(&dep)
		if err != nil {
			return err
		}
		fmt.Println(dep)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func task1(sqlDB *sql.DB, date string) error {
	err := recreate(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти все отделы, в которых работает более 10 сотрудников.")
	err = SQLquery1(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти сотрудников, которые не выходят с рабочего места в течение всего рабочего дня.")
	err = SQLquery2(sqlDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти все отделы, в которых есть сотрудники, опоздавшие в определённую дату. ")

	return SQLquery3(sqlDB, date)
}

func GORMquery1(gormDB *gorm.DB) error {
	rows, err := gormDB.Select("department").Table("employees").Group("department").Having("count(id) > 10").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var dep string

	for rows.Next() {
		err = rows.Scan(&dep)
		if err != nil {
			return err
		}
		fmt.Println(dep)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func GORMquery2(gormDB *gorm.DB) error {
	rows, err := gormDB.Select("distinct employee_id").Table("times").Where("type = 2 and date = '2022-12-13'").Group("employee_id, date").Having("count(*) = 1").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int

	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

func GORMquery3(gormDB *gorm.DB, date string) error {
	subQuery := gormDB.Select("employee_id").Table("times").Where("type = 1 and date = ?", date).Group("employee_id").Having("min(time) > '09:00'")
	rows, err := gormDB.Select("distinct department").Table("employees").Where("id in (?)", subQuery).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var dep string

	for rows.Next() {
		err = rows.Scan(&dep)
		if err != nil {
			return err
		}
		fmt.Println(dep)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func task2(gormDB *gorm.DB, date string) error {
	fmt.Println("Найти все отделы, в которых работает более 10 сотрудников.")
	err := GORMquery1(gormDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти сотрудников, которые не выходят с рабочего места в течение всего рабочего дня.")
	err = GORMquery2(gormDB)
	if err != nil {
		return err
	}

	fmt.Println("Найти все отделы, в которых есть сотрудники, опоздавшие в определённую дату. ")

	return GORMquery3(gormDB, date)
}

func main() {
	data, err := os.ReadFile("AuthDB.json")
	if err != nil {
		panic(err)
	}

	var dbRequest AuthDB
	err = json.Unmarshal(data, &dbRequest)
	if err != nil {
		panic(err)
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

	var date string
	fmt.Print("Enter the date: ")
	_, err = fmt.Scanf("%s", &date)
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================================GORM=========================================")
	err = task2(gormDB, date)
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================================SQL==========================================")
	err = task1(sqlDB, date)
	if err != nil {
		panic(err)
	}

	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}
