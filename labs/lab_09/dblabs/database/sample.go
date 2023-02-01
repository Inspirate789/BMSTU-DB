package database

import "dblabs/database/types"

func (d *Database) InsertTermRu(term types.Term) (int, error) {
	script, err := Asset("database/sql/sample_insert.sql")
	if err != nil {
		return 0, err
	}

	rows, err := d.sqlDB.Query(string(script), term.ModelID, term.Class, term.WordsCount, term.RegDate, term.Text)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (d *Database) SelectTermRu(id int) (string, error) {
	script, err := Asset("database/sql/sample_select.sql")
	if err != nil {
		return "", err
	}

	rows, err := d.sqlDB.Query(string(script), id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var text string
	rows.Next()
	err = rows.Scan(&text)
	if err != nil {
		return text, err
	}

	return text, nil
}

func (d *Database) DeleteTermRu(id int) (string, error) {
	script, err := Asset("database/sql/sample_delete.sql")
	if err != nil {
		return "", err
	}

	rows, err := d.sqlDB.Query(string(script), id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var text string
	rows.Next()
	err = rows.Scan(&text)
	if err != nil {
		return text, err
	}

	return text, nil
}
