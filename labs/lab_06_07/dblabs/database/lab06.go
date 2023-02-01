package database

import (
	"database/sql"
	"dblabs/database/types"
	"dblabs/generator"
	"errors"
)

func (d *Database) RecreateTables() (sql.Result, error) {
	script, err := Asset("database/sql/recreate.sql")
	if err != nil {
		return nil, err
	}

	return d.sqlDB.Exec(string(script))
}

func (d *Database) FillTables() (sql.Result, error) {
	err := generator.GenerateData("generator/data")
	if err != nil {
		return nil, err
	}

	script, err := Asset("database/sql/fill.sql")
	if err != nil {
		return nil, err
	}

	return d.sqlDB.Exec(string(script))
}

func (d *Database) SelectAvgWordsCount() ([]float32, error) {
	script, err := Asset("database/sql/scalar_query.sql")
	if err != nil {
		return nil, err
	}

	rows, err := d.sqlDB.Query(string(script))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	avgs := make([]float32, 3)

	for rows.Next() {
		var avg float32
		err = rows.Scan(&avg)
		if err != nil {
			return avgs, err
		}
		avgs = append(avgs, avg)
	}
	if err = rows.Err(); err != nil {
		return avgs, err
	}

	if len(avgs) < 4 {
		return avgs, errors.New("empty result")
	}

	return avgs[3:], nil
}

func (d *Database) SelectTermsWithLongestModel() ([]types.Term, error) {
	script, err := Asset("database/sql/multiple_join.sql")
	if err != nil {
		return nil, err
	}

	rows, err := d.sqlDB.Query(string(script))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []types.Term

	for rows.Next() {
		term, err := types.ScanTerm(rows)
		if err != nil {
			return res, err
		}

		res = append(res, term)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (d *Database) SelectModelsAvgWordsCount(class int) (float32, error) {
	script, err := Asset("database/sql/cte.sql")
	if err != nil {
		return 0, err
	}

	rows, err := d.sqlDB.Query(string(script), class)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var avg float32

	rows.Next()
	err = rows.Scan(&avg)
	if err != nil {
		return avg, err
	}

	if err = rows.Err(); err != nil {
		return avg, err
	}

	return avg, nil
}

func (d *Database) SelectTableMetadataByQuery(tableName string) ([]string, []string, error) {
	script, err := Asset("database/sql/metadata_query.sql")
	if err != nil {
		return nil, nil, err
	}

	rows, err := d.sqlDB.Query(string(script), tableName)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var columns, dataTypes []string
	var tmpColumn, tmpType string

	for rows.Next() {
		err = rows.Scan(&tmpColumn, &tmpType)
		if err != nil {
			return columns, dataTypes, err
		}

		columns = append(columns, tmpColumn)
		dataTypes = append(dataTypes, tmpType)
	}
	if err = rows.Err(); err != nil {
		return columns, dataTypes, err
	}

	return columns, dataTypes, nil
}

func (d *Database) SelectLatestUserRegDate() (string, error) {
	script, err := Asset("database/sql/scalar_function_init.sql")
	if err != nil {
		return "", err
	}

	_, err = d.sqlDB.Query(string(script))
	if err != nil {
		return "", err
	}

	script, err = Asset("database/sql/scalar_function_call.sql")
	if err != nil {
		return "", err
	}

	rows, err := d.sqlDB.Query(string(script))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var lastRegDate string

	rows.Next()
	err = rows.Scan(&lastRegDate)
	if err != nil {
		return lastRegDate, err
	}

	if err = rows.Err(); err != nil {
		return lastRegDate, err
	}

	return lastRegDate, nil
}

func (d *Database) GetUser(id int) ([]types.User, error) {
	script, err := Asset("database/sql/table_function_init.sql")
	if err != nil {
		return nil, err
	}

	_, err = d.sqlDB.Query(string(script))
	if err != nil {
		return nil, err
	}

	script, err = Asset("database/sql/table_function_call.sql")
	if err != nil {
		return nil, err
	}

	rows, err := d.sqlDB.Query(string(script), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []types.User

	for rows.Next() {
		user, err := types.ScanUser(rows)
		if err != nil {
			return res, err
		}
		res = append(res, user)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (d *Database) InsertTermEn(term types.Term) error {
	script, err := Asset("database/sql/stored_procedure_init.sql")
	if err != nil {
		return err
	}

	_, err = d.sqlDB.Query(string(script))
	if err != nil {
		return err
	}

	script, err = Asset("database/sql/stored_procedure_call.sql")
	if err != nil {
		return err
	}

	_, err = d.sqlDB.Query(string(script), term.ModelID, term.Class, term.WordsCount, term.RegDate, term.Text)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectCurrentDBUser() (string, error) {
	script, err := Asset("database/sql/system_function.sql")
	if err != nil {
		return "", err
	}

	rows, err := d.sqlDB.Query(string(script))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var username string

	rows.Next()
	err = rows.Scan(&username)
	if err != nil {
		return username, err
	}

	if err = rows.Err(); err != nil {
		return username, err
	}

	return username, nil
}
