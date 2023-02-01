package database

import (
	"dblabs/database/types"
	"errors"
)

func (d *Database) SelectTerms(lang string) ([]types.Term, error) {
	if lang != "ru" && lang != "en" {
		return nil, errors.New("incorrect language")
	}

	rows, err := d.sqlDB.Query("SELECT * FROM public.terms_" + lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []types.Term

	for rows.Next() {
		term, err := types.ScanTerm(rows)
		if err != nil {
			return terms, err
		}
		terms = append(terms, term)
	}
	if err = rows.Err(); err != nil {
		return terms, err
	}

	return terms, err
}

func (d *Database) SelectAllTerms() ([]types.Term, error) {
	rows, err := d.sqlDB.Query("SELECT * FROM public.terms_ru UNION SELECT * FROM public.terms_en")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []types.Term

	for rows.Next() {
		term, err := types.ScanTerm(rows)
		if err != nil {
			return terms, err
		}
		terms = append(terms, term)
	}
	if err = rows.Err(); err != nil {
		return terms, err
	}

	return terms, err
}

func (d *Database) SelectModels() ([]types.Model, error) {
	rows, err := d.sqlDB.Query("SELECT * FROM public.models")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []types.Model

	for rows.Next() {
		model, err := types.ScanModel(rows)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, err
}

func (d *Database) SelectAllPos() ([]types.PartOfSpeech, error) {
	rows, err := d.sqlDB.Query("SELECT * FROM public.parts_of_speech")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partsOfSpeech []types.PartOfSpeech

	for rows.Next() {
		pos, err := types.ScanPos(rows)
		if err != nil {
			return nil, err
		}
		partsOfSpeech = append(partsOfSpeech, pos)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return partsOfSpeech, err
}

func (d *Database) SelectPos(id int) ([]types.PartOfSpeech, error) {
	rows, err := d.sqlDB.Query("SELECT * FROM public.parts_of_speech WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partsOfSpeech []types.PartOfSpeech

	for rows.Next() {
		pos, err := types.ScanPos(rows)
		if err != nil {
			return nil, err
		}
		partsOfSpeech = append(partsOfSpeech, pos)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return partsOfSpeech, err
}
