package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/Inspirate789/go-randomdata"
)

type PartOfSpeech struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	RegDate string `json:"registration_date"`
}

func RandPartOfSpeech(index int) PartOfSpeech {
	var pos PartOfSpeech
	pos.Name = string('a' + rune(index))
	pos.RegDate = randomdata.FullDateInRange("2019-01-01", "2019-12-31")

	return pos
}

func PosToSlice(pos PartOfSpeech) []string {
	return []string{pos.Name, pos.RegDate}
}

func ScanPos(rows *sql.Rows) (PartOfSpeech, error) {
	var pos PartOfSpeech
	err := rows.Scan(&pos.ID, &pos.Name, &pos.RegDate)

	return pos, err
}

// Value make the PartOfSpeech struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (pos PartOfSpeech) Value() (driver.Value, error) {
	return json.Marshal(pos)
}

// Scan make the PartOfSpeech struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (pos *PartOfSpeech) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pos)
}
