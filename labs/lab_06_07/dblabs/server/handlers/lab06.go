package handlers

import (
	"dblabs/database"
	"dblabs/database/types"
	"errors"
	"fmt"
	"github.com/tjarratt/babble"
	"strconv"
	"strings"
)

type handler60 struct{}

func (h handler60) Title() string {
	return "Recreate all tables"
}

func (h handler60) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.RecreateTables()
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	lid, err := res.LastInsertId()
	if err != nil {
		builder.WriteString(fmt.Sprintf("last insert id: %d; ", lid))
	}
	ra, err := res.RowsAffected()
	if err != nil {
		builder.WriteString(fmt.Sprintf("rows affected: %d; ", ra))
	}

	return builder.String(), nil
}

type handler61 struct{}

func (h handler61) Title() string {
	return "Fill all tables"
}

func (h handler61) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.FillTables()
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	lid, err := res.LastInsertId()
	if err != nil {
		builder.WriteString(fmt.Sprintf("last insert id: %d; ", lid))
	}
	ra, err := res.RowsAffected()
	if err != nil {
		builder.WriteString(fmt.Sprintf("rows affected: %d; ", ra))
	}

	return builder.String(), nil
}

type handler62 struct{}

func (h handler62) Title() string {
	return "Select average terms words count for every class"
}

func (h handler62) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectAvgWordsCount()
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for i, num := range res {
		builder.WriteString(fmt.Sprintf("average words count for class %d: %f\n", i+1, num))
	}

	return builder.String(), nil
}

type handler63 struct{}

func (h handler63) Title() string {
	return "Select terms with longest model"
}

func (h handler63) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectTermsWithLongestModel()
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for _, str := range res {
		builder.WriteString(str.Text + "\n")
	}

	return builder.String(), nil
}

type handler64 struct{}

func (h handler64) Title() string {
	return "Select average words count in models a given class"
}

func (h handler64) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	class, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	res, err := db.SelectModelsAvgWordsCount(class)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(float64(res), 'f', -1, 32), nil
}

type handler65 struct{}

func (h handler65) Title() string {
	return "Select metadata from table with a given name"
}

func (h handler65) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	columns, dataTypes, err := db.SelectTableMetadataByQuery(args[0])
	if err != nil {
		return "", err
	}

	if len(columns) != len(dataTypes) {
		return "", errors.New("database returned an incorrect response")
	}

	builder := strings.Builder{}
	for i := range columns {
		builder.WriteString(fmt.Sprintf("%s: type %s\n", columns[i], dataTypes[i]))
	}

	return builder.String(), nil
}

type handler66 struct{}

func (h handler66) Title() string {
	return "Select date of last user registration"
}

func (h handler66) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectLatestUserRegDate()
	if err != nil {
		return "", err
	}

	return res, nil
}

type handler67 struct{}

func (h handler67) Title() string {
	return "Get username by user ID"
}

func (h handler67) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	res, err := db.GetUser(id)
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for _, str := range res {
		builder.WriteString(str.Name + "\n")
	}

	return builder.String(), nil
}

type handler68 struct{}

func (h handler68) Title() string {
	return "Insert random term (en)"
}

func (h handler68) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	var partsOfSpeech = []types.PartOfSpeech{types.RandPartOfSpeech(0)}
	var models = []types.Model{types.RandModel(partsOfSpeech)}
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	term := types.RandTerm(models, babbler)

	err := db.InsertTermEn(term)
	if err != nil {
		return term.Text, err
	}

	return term.Text, nil
}

type handler69 struct{}

func (h handler69) Title() string {
	return "Select username of current DB user"
}

func (h handler69) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectCurrentDBUser()
	if err != nil {
		return "", err
	}

	return res, nil
}
