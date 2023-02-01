package handlers

import (
	"dblabs/database"
	"dblabs/database/types"
	"errors"
	"fmt"
	"github.com/tjarratt/babble"
	"strconv"
)

type handler00 struct{}

func (h handler00) Title() string {
	return "Insert random term (ru)"
}

func (h handler00) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	var partsOfSpeech = []types.PartOfSpeech{types.RandPartOfSpeech(0)}
	var models = []types.Model{types.RandModel(partsOfSpeech)}
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	term := types.RandTerm(models, babbler)

	id, err := db.InsertTermRu(term)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

type handler01 struct{}

func (h handler01) Title() string {
	return "Select term (ru) by ID"
}

func (h handler01) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	text, err := db.SelectTermRu(id)
	if err != nil {
		return "", err
	}

	return text, nil
}

type handler02 struct{}

func (h handler02) Title() string {
	return "Delete term (ru) by ID"
}

func (h handler02) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("incorrect args")
	}

	return db.DeleteTermRu(id)
}
