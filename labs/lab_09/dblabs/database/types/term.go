package types

import (
	"database/sql"
	"fmt"
	"github.com/Inspirate789/go-randomdata"
	"github.com/tjarratt/babble"
	"math/rand"
)

type Term struct {
	ID, ModelID, Class, WordsCount int
	RegDate, Text                  string
}

func (Term) TableName() string {
	return "terms_ru"
}

func RandTerm(models []Model, babbler babble.Babbler) Term {
	var term Term

	term.ModelID = rand.Intn(len(models)) + 1
	term.Class = models[term.ModelID-1].Class
	term.RegDate = randomdata.FullDateInRange("2021-01-01", "2022-08-31")
	term.WordsCount = models[term.ModelID-1].WordsCount
	babbler.Count = term.WordsCount
	term.Text = babbler.Babble()

	return term
}

func TermToSlice(term Term) []string {
	return []string{
		fmt.Sprint(term.ModelID),
		fmt.Sprint(term.Class),
		fmt.Sprint(term.WordsCount),
		term.RegDate,
		term.Text,
	}
}

func ScanTerm(rows *sql.Rows) (Term, error) {
	var term Term
	err := rows.Scan(&term.ID, &term.ModelID, &term.Class, &term.WordsCount, &term.RegDate, &term.Text)

	return term, err
}
