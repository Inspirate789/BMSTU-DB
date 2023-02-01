package types

import (
	"database/sql"
	"fmt"
	"github.com/Inspirate789/go-randomdata"
	"math/rand"
)

type User struct {
	ID, TermsCount                             int
	Name, RegDate, FirstTermDate, LastTermDate string
}

func RandUser() User {
	var user User
	user.Name = randomdata.FullName(randomdata.Male)
	user.RegDate = randomdata.FullDateInRange("2021-01-01", "2021-12-31")
	user.FirstTermDate = user.RegDate
	user.LastTermDate = randomdata.FullDateInRange("2022-01-01", "2022-08-31")
	user.TermsCount = rand.Intn(500) + 2

	return user
}

func UserToSlice(user User) []string {
	return []string{
		user.Name,
		user.RegDate,
		user.FirstTermDate,
		user.LastTermDate,
		fmt.Sprint(user.TermsCount),
	}
}

func ScanUser(rows *sql.Rows) (User, error) {
	var user User
	err := rows.Scan(&user.ID, &user.Name, &user.RegDate, &user.FirstTermDate, &user.LastTermDate, &user.TermsCount)

	return user, err
}
