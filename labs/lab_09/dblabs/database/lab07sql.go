package database

import (
	"dblabs/database/types"
	"encoding/json"
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/buger/jsonparser"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"os"
	"strings"
	"time"
)

func (d *Database) SelectTermsRuByClass(class int) ([]types.Term, error) {
	terms, err := d.SelectTerms("ru")
	if err != nil {
		return terms, err
	}

	terms = lo.Filter[types.Term](terms, func(t types.Term, index int) bool {
		return t.Class == class
	})

	return terms, nil
}

func (d *Database) SelectTermsEnByWordsCount(wordsCount int) ([]types.Term, error) {
	terms, err := d.SelectTerms("en")
	if err != nil {
		return terms, err
	}

	models, err := d.SelectModels()
	if err != nil {
		return terms, err
	}

	models = lo.Filter[types.Model](models, func(t types.Model, index int) bool {
		return t.WordsCount == wordsCount
	})

	modelsID := make([]int, len(models))

	for i := range models {
		modelsID[i] = models[i].ID
	}

	terms = lo.Filter[types.Term](terms, func(t types.Term, index int) bool {
		return lo.Contains[int](modelsID, t.ModelID)
	})

	return terms, nil
}

func (d *Database) SelectVeryLongTerms() ([]types.Term, error) {
	terms, err := d.SelectAllTerms()
	if err != nil {
		return terms, err
	}

	resMap := lop.GroupBy[types.Term, string](terms, func(t types.Term) string {
		if t.WordsCount < 3 {
			return "Short"
		} else if t.WordsCount < 5 {
			return "Normal"
		} else if t.WordsCount < 8 {
			return "Long"
		}
		return "Very long"
	})

	return resMap["Very long"], nil
}

func (d *Database) SelectTermsOfModel(modelID int) ([]string, error) {
	terms, err := d.SelectAllTerms()
	if err != nil {
		return nil, err
	}

	models, err := d.SelectModels()
	if err != nil {
		return nil, err
	}

	res := make(map[int][]string)

	From(models).GroupJoin(From(terms),
		func(m interface{}) interface{} { return m.(types.Model).ID },
		func(t interface{}) interface{} { return t.(types.Term).ModelID },
		func(outer interface{}, inners []interface{}) interface{} {
			t := make([]string, len(inners))

			for i, inner := range inners {
				t[i] = inner.(types.Term).Text
			}

			return KeyValue{
				Key:   outer.(types.Model).ID,
				Value: t,
			}
		},
	).ToMap(&res)

	return res[modelID], nil
}

func (d *Database) SelectTermsWithPOS(partOfSpeech string) ([]string, error) {
	terms, err := d.SelectAllTerms()
	if err != nil {
		return nil, err
	}

	models, err := d.SelectModels()
	if err != nil {
		return nil, err
	}

	models = lo.Filter[types.Model](models, func(m types.Model, index int) bool {
		return strings.Contains(m.Text, partOfSpeech)
	})

	var res []string

	From(models).Join(From(terms),
		func(m interface{}) interface{} { return m.(types.Model).ID },
		func(t interface{}) interface{} { return t.(types.Term).ModelID },
		func(outer interface{}, inner interface{}) interface{} {
			return inner.(types.Term).Text
		},
	).ToSlice(&res)

	return res, nil
}

func (d *Database) ReadPosFromJSON(filename string) ([]types.PartOfSpeech, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var partsOfSpeech []types.PartOfSpeech
	err = json.Unmarshal(data, &partsOfSpeech)
	if err != nil {
		return nil, err
	}

	return partsOfSpeech, nil
}

func (d *Database) WritePosToJSON(partsOfSpeech []types.PartOfSpeech, filename string) error {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(partsOfSpeech)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) WritePosTableToJSON(filename string) error {
	partsOfSpeech, err := d.SelectAllPos()
	if err != nil {
		return err
	}

	return d.WritePosToJSON(partsOfSpeech, filename)
}

func (d *Database) SetPosRegDateInJSON(id int, newDate time.Time, filename string) error {
	partsOfSpeech, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	newDateByte, err := newDate.MarshalJSON()
	if err != nil {
		return err
	}

	partsOfSpeech, err = jsonparser.Set(partsOfSpeech, newDateByte, fmt.Sprintf("[%d]", id-1), "registration_date")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, partsOfSpeech, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertPosToJSON(pos types.PartOfSpeech, filename string) error {
	partsOfSpeech, err := d.ReadPosFromJSON(filename)
	if err != nil {
		return err
	}

	pos.ID = len(partsOfSpeech)
	partsOfSpeech = append(partsOfSpeech, pos)

	return d.WritePosToJSON(partsOfSpeech, filename)
}
