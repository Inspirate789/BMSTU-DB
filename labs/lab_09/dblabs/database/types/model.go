package types

import (
	"database/sql"
	"fmt"
	"github.com/Inspirate789/go-randomdata"
	"math/rand"
	"strings"
)

type Model struct {
	ID, Class, WordsCount int
	RegDate, Text         string
}

func (Model) TableName() string {
	return "models"
}

func RandModel(partsOfSpeech []PartOfSpeech) Model {
	var model Model

	model.WordsCount = rand.Intn(10) + 1

	switch {
	case model.WordsCount == 1:
		model.Class = 1
	case model.WordsCount <= 7:
		model.Class = 2
	default:
		model.Class = 3
	}

	model.RegDate = randomdata.FullDateInRange("2020-01-01", "2020-12-31")

	builder := strings.Builder{}
	builder.WriteString(partsOfSpeech[rand.Intn(len(partsOfSpeech))].Name)
	for j := 1; j < model.WordsCount; j++ {
		builder.WriteRune('+')
		builder.WriteString(partsOfSpeech[rand.Intn(len(partsOfSpeech))].Name)
	}
	model.Text = builder.String()

	return model
}

func ModelToSlice(model Model) []string {
	return []string{
		fmt.Sprint(model.Class),
		fmt.Sprint(model.WordsCount),
		model.RegDate,
		model.Text,
	}
}

func ScanModel(rows *sql.Rows) (Model, error) {
	var model Model
	err := rows.Scan(&model.ID, &model.Class, &model.WordsCount, &model.RegDate, &model.Text)

	return model, err
}
