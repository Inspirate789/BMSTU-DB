package database

import (
	"dblabs/database/types"
)

func (d *Database) GormSelectTermsRuByClass(class int) ([]types.Term, error) {
	var terms []types.Term

	result := d.gormDB.Where("class = ?", class).Find(&terms)
	if result.Error != nil {
		return nil, result.Error
	}

	return terms, nil
}

func (d *Database) GormSelectTermsRuOfModel(modelID int) ([]string, error) {
	rows, err := d.gormDB.Table("terms_ru").Select("terms_ru.text").Joins("join models on terms_ru.model_id = models.id and models.id = ?", modelID).Rows()
	if err != nil {
		return nil, err
	}

	var terms []string

	for rows.Next() {
		var term string
		err = rows.Scan(&term)
		if err != nil {
			return terms, err
		}
		terms = append(terms, term)
	}
	if err = rows.Err(); err != nil {
		return terms, err
	}

	return terms, nil
}

func (d *Database) GormInsertTermEn(term types.Term) error {
	result := d.gormDB.Create(&term)

	return result.Error
}

func (d *Database) GormUpdateUsername(id int, newName string) error {
	result := d.gormDB.Model(&types.User{}).Where("id = ?", id).Update("name", newName)

	return result.Error
}

func (d *Database) GormDeleteUser(id int) error {
	result := d.gormDB.Table("users_and_terms_ru").Where("user_id = ?", id).Delete(&struct{}{})
	if result.Error != nil {
		return result.Error
	}

	result = d.gormDB.Table("users_and_terms_en").Where("user_id = ?", id).Delete(&struct{}{})
	if result.Error != nil {
		return result.Error
	}

	result = d.gormDB.Where("id = ?", id).Delete(&types.User{})

	return result.Error
}

func (d *Database) GormInsertTermEnByProc(term types.Term) error {
	script, err := Asset("database/sql/stored_procedure_init.sql")
	if err != nil {
		return err
	}

	result := d.gormDB.Exec(string(script))
	if result.Error != nil {
		return result.Error
	}

	script, err = Asset("database/sql/stored_procedure_call.sql")
	if err != nil {
		return err
	}

	result = d.gormDB.Exec(string(script), term.ModelID, term.Class, term.WordsCount, term.RegDate, term.Text)

	return result.Error
}
