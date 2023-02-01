// ВНИМАНИЕ: в ходе написания этого скрипта была локально отредактирована
// библиотека "github.com/Pallinder/go-randomdata"!!!!!
// Выходной формат даты был заменён на входной: (34-37 строки файла random_data.go)
// const (
// 	DateInputLayout  = "2006-01-02"
// 	DateOutputLayout = "2006-01-02" // "Monday 2 Jan 2006" - ORIGINAL!!!
// )

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/Pallinder/go-randomdata"
	"github.com/antelman107/cyrillic_translit"
	"github.com/tjarratt/babble"
)

var dict = map[string]string{
	"а": "a",
	"б": "b",
	"в": "v",
	"г": "g",
	"д": "d",
	"е": "e",
	"ё": "yo",
	"ж": "j",
	"з": "z",
	"и": "i",
	"й": "y",
	"к": "k",
	"л": "l",
	"м": "m",
	"н": "n",
	"о": "o",
	"п": "p",
	"р": "r",
	"с": "s",
	"т": "t",
	"у": "u",
	"ф": "f",
	"х": "kh",
	"ц": "ts",
	"ч": "ch",
	"ш": "sh",
	"щ": "sch",
	"ъ": "",
	"ы": "i",
	"ь": "",
	"э": "e",
	"ю": "yu",
	"я": "ya",
	"А": "A",
	"Б": "B",
	"В": "V",
	"Г": "G",
	"Д": "D",
	"Е": "E",
	"Ё": "YO",
	"Ж": "ZH",
	"З": "Z",
	"И": "I",
	"Й": "Y",
	"К": "K",
	"Л": "L",
	"М": "M",
	"Н": "N",
	"О": "O",
	"П": "P",
	"Р": "R",
	"С": "S",
	"Т": "T",
	"У": "U",
	"Ф": "F",
	"Х": "KH",
	"Ц": "TS",
	"Ч": "CH",
	"Ш": "SH",
	"Щ": "SCH",
	"Ъ": "",
	"Ы": "I",
	"Ь": "",
	"Э": "E",
	"Ю": "YU",
	"Я": "YA",
} // NEED TO REVERSE!!!!!

type user struct {
	id, terms_count                                 int
	name, reg_date, first_term_date, last_term_date string
}

type context struct {
	id, terms_count, words_count int
	reg_date, text               string
}

type term struct {
	id, model_id, class, freq, words_count int
	reg_date, text                         string
}

type model struct {
	id, class, words_count int
	reg_date, text         string
}

type part_of_speech struct {
	id       int
	reg_date string
	name     rune
}

func main() {
	babbler := babble.NewBabbler()
	babbler.Separator = " "

	csv_file, err := os.Create("pos.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer := csv.NewWriter(csv_file)

	const pos_count = 25
	var parts_of_speech [pos_count]part_of_speech

	// const words_count = 1000
	// const pos_words_count = words_count / pos_count
	// words := make(map[part_of_speech][pos_words_count]string)

	for i := 0; i < pos_count; i++ {
		parts_of_speech[i].id = i
		parts_of_speech[i].name = 'a' + rune(i)
		parts_of_speech[i].reg_date = randomdata.FullDateInRange("2019-01-01", "2019-12-31")

		row := []string{
			fmt.Sprint(parts_of_speech[i].id),
			string(parts_of_speech[i].name),
			parts_of_speech[i].reg_date,
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("models.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	const models_count = 1000
	var models [models_count]model

	for i := 0; i < models_count; i++ {
		models[i].id = i
		models[i].words_count = rand.Intn(10) + 1

		switch {
		case models[i].words_count == 1:
			models[i].class = 1
		case models[i].words_count <= 7:
			models[i].class = 2
		default:
			models[i].class = 3
		}

		models[i].reg_date = randomdata.FullDateInRange("2020-01-01", "2020-12-31")
		models[i].text = string(parts_of_speech[rand.Intn(pos_count)].name)

		for j := 1; j < models[i].words_count; j++ {
			models[i].text += "+" + string(parts_of_speech[rand.Intn(pos_count)].name)
		}

		row := []string{
			fmt.Sprint(models[i].id),
			fmt.Sprint(models[i].class),
			fmt.Sprint(models[i].words_count),
			models[i].reg_date,
			models[i].text,
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("users.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	const users_count = 5000
	var users [users_count]user

	for i := 0; i < users_count; i++ {
		users[i].id = i
		users[i].name = randomdata.FullName(randomdata.Male)
		users[i].reg_date = randomdata.FullDateInRange("2021-01-01", "2021-12-31")
		users[i].first_term_date = users[i].reg_date
		users[i].last_term_date = randomdata.FullDateInRange("2022-01-01", "2022-08-31")
		users[i].terms_count = rand.Intn(500) + 2

		row := []string{
			fmt.Sprint(users[i].id),
			users[i].name,
			users[i].reg_date,
			users[i].first_term_date,
			users[i].last_term_date,
			fmt.Sprint(users[i].terms_count),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("terms_en.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	const terms_count = 5000
	var terms_ru, terms_en [terms_count]term

	for i := 0; i < terms_count; i++ {
		terms_en[i].id = i
		terms_en[i].model_id = rand.Intn(models_count)
		terms_en[i].class = models[terms_en[i].model_id].class
		terms_en[i].reg_date = randomdata.FullDateInRange("2021-01-01", "2022-08-31")
		terms_en[i].words_count = rand.Intn(10) + 1
		babbler.Count = terms_en[i].words_count
		terms_en[i].text = babbler.Babble()

		row := []string{
			fmt.Sprint(terms_en[i].id),
			fmt.Sprint(terms_en[i].model_id),
			fmt.Sprint(terms_en[i].class),
			fmt.Sprint(terms_en[i].words_count),
			terms_en[i].reg_date,
			terms_en[i].text,
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("terms_ru.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	for i := 0; i < terms_count; i++ {
		terms_ru[i].id = i
		terms_ru[i].model_id = rand.Intn(models_count)
		terms_ru[i].class = models[terms_ru[i].model_id].class
		terms_ru[i].reg_date = randomdata.FullDateInRange("2021-01-01", "2022-08-31")
		terms_ru[i].words_count = rand.Intn(10) + 1
		babbler.Count = terms_ru[i].words_count
		terms_ru[i].text = cyrillic_translit.DoDict(babbler.Babble(), dict)

		row := []string{
			fmt.Sprint(terms_ru[i].id),
			fmt.Sprint(terms_ru[i].model_id),
			fmt.Sprint(terms_ru[i].class),
			fmt.Sprint(terms_ru[i].words_count),
			terms_ru[i].reg_date,
			terms_ru[i].text,
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("contexts.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	const contexts_count = 1000
	var contexts [contexts_count]context

	for i := 0; i < contexts_count; i++ {
		contexts[i].id = i
		contexts[i].reg_date = randomdata.FullDateInRange("2022-09-01", "2022-09-17")
		contexts[i].terms_count = rand.Intn(10) + 1

		for j := 0; j < contexts[i].terms_count; j++ {
			var cur_term term

			if j%2 == 0 {
				cur_term = terms_ru[rand.Intn(terms_count)]
			} else {
				cur_term = terms_en[rand.Intn(terms_count)]
			}

			contexts[i].text += " " + cur_term.text
			contexts[i].words_count += cur_term.words_count
		}

		row := []string{
			fmt.Sprint(contexts[i].id),
			fmt.Sprint(contexts[i].terms_count),
			fmt.Sprint(contexts[i].words_count),
			contexts[i].reg_date,
			contexts[i].text,
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	// Linking tables:

	csv_file, err = os.Create("users_and_terms_ru.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	const links_count = 5000

	for i := 0; i < links_count; i++ {
		row := []string{
			fmt.Sprint(users[rand.Intn(users_count)].id),
			fmt.Sprint(terms_ru[rand.Intn(terms_count)].id),
			randomdata.FullDateInRange("2021-01-01", "2022-08-31"),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("users_and_terms_en.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	for i := 0; i < links_count; i++ {
		row := []string{
			fmt.Sprint(users[rand.Intn(users_count)].id),
			fmt.Sprint(terms_en[rand.Intn(terms_count)].id),
			randomdata.FullDateInRange("2021-01-01", "2022-08-31"),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("contexts_and_terms_ru.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	for i := 0; i < links_count; i++ {
		row := []string{
			fmt.Sprint(contexts[rand.Intn(contexts_count)].id),
			fmt.Sprint(terms_ru[rand.Intn(terms_count)].id),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("contexts_and_terms_en.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	for i := 0; i < links_count; i++ {
		row := []string{
			fmt.Sprint(contexts[rand.Intn(contexts_count)].id),
			fmt.Sprint(terms_en[rand.Intn(terms_count)].id),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()

	csv_file, err = os.Create("models_and_pos.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csv_writer = csv.NewWriter(csv_file)

	for i := 0; i < links_count; i++ {
		row := []string{
			fmt.Sprint(models[rand.Intn(models_count)].id),
			fmt.Sprint(parts_of_speech[rand.Intn(pos_count)].id),
		}

		csv_writer.Write(row)
	}

	csv_writer.Flush()
	csv_file.Close()
}
