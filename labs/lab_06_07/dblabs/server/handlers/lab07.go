package handlers

import (
	"dblabs/database"
	"dblabs/database/types"
	"errors"
	"fmt"
	"github.com/tjarratt/babble"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type handler70 struct{}

func (h handler70) Title() string {
	return "Select terms_ru with a given class (LINQ)"
}

func (h handler70) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	class, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	res, err := db.SelectTermsRuByClass(class)
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

type handler71 struct{}

func (h handler71) Title() string {
	return "Select terms (en) with model of a given words count (LINQ)"
}

func (h handler71) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	wordsCount, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	res, err := db.SelectTermsEnByWordsCount(wordsCount)
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

type handler72 struct{}

func (h handler72) Title() string {
	return "Select very long terms (LINQ)"
}

func (h handler72) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectVeryLongTerms()
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

type handler73 struct{}

func (h handler73) Title() string {
	return "Select terms of model with a given ID (LINQ)"
}

func (h handler73) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	modelID, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	terms, err := db.SelectTermsOfModel(modelID)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for _, term := range terms {
		builder.WriteString(fmt.Sprintf("%s\n", term))
	}

	return builder.String(), nil
}

type handler74 struct{}

func (h handler74) Title() string {
	return "Select terms with a given pos (LINQ)"
}

func (h handler74) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	terms, err := db.SelectTermsWithPOS(args[0])
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for _, term := range terms {
		builder.WriteString(fmt.Sprintf("%s\n", term))
	}

	return builder.String(), nil
}

type handler75 struct{}

func (h handler75) Title() string {
	return "Write parts of speech table to JSON"
}

func (h handler75) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	err := db.WritePosTableToJSON("database/json/pos.json")
	if err != nil {
		return "", err
	}

	return "", nil
}

type handler76 struct{}

func (h handler76) Title() string {
	return "Read parts of speech from JSON"
}

func (h handler76) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	partsOfSpeech, err := db.ReadPosFromJSON("database/json/pos.json") // TODO: add absolute path everywhere (by pwd)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for _, pos := range partsOfSpeech {
		builder.WriteString(fmt.Sprintf("%v\n", pos))
	}

	return builder.String(), nil
}

type handler77 struct{}

func (h handler77) Title() string {
	return "Refresh the registration date of parts of speech in JSON by id"
}

func (h handler77) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	posID, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	curDate := time.Now()
	err = db.SetPosRegDateInJSON(posID, curDate, "database/json/pos.json")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Registration date (%s) is set successfully", curDate), nil
}

type handler78 struct{}

func (h handler78) Title() string {
	return "Insert random part of speech to JSON"
}

func (h handler78) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	pos := types.RandPartOfSpeech(rand.Intn(25))
	err := db.InsertPosToJSON(pos, "database/json/pos.json")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Part of speech (%s) inserted successfully", pos.Name), nil
}

type handler79 struct{}

func (h handler79) Title() string {
	return "Select terms_ru with a given class (GORM)"
}

func (h handler79) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	class, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	res, err := db.GormSelectTermsRuByClass(class)
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

type handler710 struct{}

func (h handler710) Title() string {
	return "Select terms (ru) of model with a given ID (GORM)"
}

func (h handler710) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	modelID, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	terms, err := db.GormSelectTermsRuOfModel(modelID)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for _, term := range terms {
		builder.WriteString(fmt.Sprintf("%s\n", term))
	}

	return builder.String(), nil
}

type handler711 struct{}

func (h handler711) Title() string {
	return "Insert random term (en) (GORM)"
}

func (h handler711) Execute(db *database.Database, args []string) (string, error) {
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

type handler712 struct{}

func (h handler712) Title() string {
	return "Set username by ID (GORM)"
}

func (h handler712) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 2, len(args)))
	}

	userID, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	return "", db.GormUpdateUsername(userID, args[0])
}

type handler713 struct{}

func (h handler713) Title() string {
	return "Delete username by ID (GORM)"
}

func (h handler713) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	return "", db.GormDeleteUser(userID)
}

type handler714 struct{}

func (h handler714) Title() string {
	return "Insert random term (en) by stored procedure (GORM)"
}

func (h handler714) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	var partsOfSpeech = []types.PartOfSpeech{types.RandPartOfSpeech(0)}
	var models = []types.Model{types.RandModel(partsOfSpeech)}
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	term := types.RandTerm(models, babbler)

	err := db.GormInsertTermEnByProc(term)
	if err != nil {
		return term.Text, err
	}

	return term.Text, nil
}
