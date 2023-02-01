package generator

import (
	"dblabs/database/types"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Inspirate789/go-randomdata"
	"github.com/tjarratt/babble"
	"math/rand"
	"os"
	"sync"
)

func GeneratePos(dataPath string, posCount int) ([]types.PartOfSpeech, error) {
	csvFile, err := os.Create(dataPath + "/pos.csv")
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	partsOfSpeech := make([]types.PartOfSpeech, posCount)

	for i := range partsOfSpeech {
		partsOfSpeech[i] = types.RandPartOfSpeech(i)

		if err = csvWriter.Write(types.PosToSlice(partsOfSpeech[i])); err != nil {
			return nil, err
		}
	}

	return partsOfSpeech, nil
}

func GenerateModels(dataPath string, modelsCount int, partsOfSpeech []types.PartOfSpeech) ([]types.Model, error) {
	csvFile, err := os.Create(dataPath + "/models.csv")
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	models := make([]types.Model, modelsCount)

	for i := range models {
		models[i] = types.RandModel(partsOfSpeech)

		if err = csvWriter.Write(types.ModelToSlice(models[i])); err != nil {
			return nil, err
		}
	}

	return models, nil
}

func GenerateUsers(dataPath string, usersCount int) ([]types.User, error) {
	csvFile, err := os.Create(dataPath + "/users.csv")
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	users := make([]types.User, usersCount)

	for i := range users {
		users[i] = types.RandUser()

		if err = csvWriter.Write(types.UserToSlice(users[i])); err != nil {
			return nil, err
		}
	}

	return users, nil
}

func GenerateTerms(dataPath string, TermsCount int, models []types.Model, locale string) ([]types.Term, error) {
	csvFile, err := os.Create(fmt.Sprintf(dataPath+"/terms_%s.csv", locale))
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	terms := make([]types.Term, TermsCount)
	babbler := babble.NewBabbler()
	babbler.Separator = " "

	for i := range terms {
		terms[i] = types.RandTerm(models, babbler)

		if err = csvWriter.Write(types.TermToSlice(terms[i])); err != nil {
			return nil, err
		}
	}

	return terms, nil
}

func GenerateContexts(dataPath string, contextsCount int, terms []types.Term) ([]types.Context, error) {
	csvFile, err := os.Create(dataPath + "/contexts.csv")
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	contexts := make([]types.Context, contextsCount)

	for i := range contexts {
		contexts[i] = types.RandContext(terms)

		row := []string{
			fmt.Sprint(contexts[i].TermsCount),
			fmt.Sprint(contexts[i].WordsCount),
			contexts[i].RegDate,
			contexts[i].Text,
		}

		if err = csvWriter.Write(row); err != nil {
			return nil, err
		}
	}

	return contexts, nil
}

func GenerateUsersAndTerms(dataPath string, linksCount int, users []types.User, terms []types.Term, locale string) error {
	csvFile, err := os.Create(fmt.Sprintf(dataPath+"/users_terms_%s.csv", locale))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for i := 0; i < linksCount; i++ {
		row := []string{
			fmt.Sprint(rand.Intn(len(users)) + 1),
			fmt.Sprint(rand.Intn(len(terms)) + 1),
			randomdata.FullDateInRange("2021-01-01", "2022-08-31"),
		}

		if err = csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func GenerateContextsAndTerms(dataPath string, linksCount int, contexts []types.Context, terms []types.Term, locale string) error {
	csvFile, err := os.Create(fmt.Sprintf(dataPath+"/contexts_terms_%s.csv", locale))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for i := 0; i < linksCount; i++ {
		row := []string{
			fmt.Sprint(rand.Intn(len(contexts)) + 1),
			fmt.Sprint(rand.Intn(len(terms)) + 1),
		}

		if err = csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func GenerateModelsAndPos(dataPath string, linksCount int, models []types.Model, partsOfSpeech []types.PartOfSpeech) error {
	csvFile, err := os.Create(dataPath + "/models_pos.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for i := 0; i < linksCount; i++ {
		row := []string{
			fmt.Sprint(rand.Intn(len(models)) + 1),
			fmt.Sprint(rand.Intn(len(partsOfSpeech)) + 1),
		}

		if err = csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func generationStep1(dataPath string, maxCount int) ([]types.PartOfSpeech, []types.Model, []types.User, error) {
	var resErr, tmpErr1, tmpErr2 error
	var users []types.User
	var partsOfSpeech []types.PartOfSpeech
	var models []types.Model
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		users, tmpErr1 = GenerateUsers(dataPath, maxCount)
		if tmpErr1 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating users: %v", tmpErr1))
			mx.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		partsOfSpeech, tmpErr2 = GeneratePos(dataPath, 25)
		if tmpErr2 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating parts of speech: %v", tmpErr2))
			mx.Unlock()
			return
		}

		models, tmpErr2 = GenerateModels(dataPath, maxCount/5, partsOfSpeech)
		if tmpErr2 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating models: %v", tmpErr2))
			mx.Unlock()
		}
	}()

	wg.Wait()

	return partsOfSpeech, models, users, resErr
}

func generationStep2(dataPath string, maxCount int, partsOfSpeech []types.PartOfSpeech, models []types.Model) ([]types.Term, []types.Term, error) {
	var resErr, tmpErr1, tmpErr2, tmpErr3 error
	var termsRu, termsEn []types.Term
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		termsRu, tmpErr1 = GenerateTerms(dataPath, maxCount, models, "ru")
		if tmpErr1 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating terms (ru): %v", tmpErr1))
			mx.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		termsEn, tmpErr2 = GenerateTerms(dataPath, maxCount, models, "en")
		if tmpErr2 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating terms (en): %v", tmpErr2))
			mx.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		tmpErr3 = GenerateModelsAndPos(dataPath, maxCount, models, partsOfSpeech)
		if tmpErr3 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating models&pos: %v", tmpErr3))
			mx.Unlock()
		}
	}()

	wg.Wait()

	return termsRu, termsEn, resErr
}

func generationStep3(dataPath string, maxCount int, termsRu []types.Term, termsEn []types.Term, users []types.User) error {
	var resErr, tmpErr1, tmpErr2, tmpErr3, tmpErr4 error
	var contexts []types.Context
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		contexts, tmpErr1 = GenerateContexts(dataPath, maxCount/5, append(termsRu, termsEn...))
		if tmpErr1 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating contexts: %v", tmpErr1))
			mx.Unlock()
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			tmpErr1 = GenerateContextsAndTerms(dataPath, maxCount, contexts, termsRu, "ru")
			if tmpErr1 != nil {
				mx.Lock()
				resErr = errors.New(fmt.Sprintf("failed generating contexts&terms (ru): %v", tmpErr1))
				mx.Unlock()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			tmpErr2 = GenerateContextsAndTerms(dataPath, maxCount, contexts, termsEn, "en")
			if tmpErr2 != nil {
				mx.Lock()
				resErr = errors.New(fmt.Sprintf("failed generating contexts&terms (en): %v", tmpErr2))
				mx.Unlock()
			}
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tmpErr3 = GenerateUsersAndTerms(dataPath, maxCount, users, termsRu, "ru")
		if tmpErr3 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating users&terms (ru): %v", tmpErr3))
			mx.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tmpErr4 = GenerateUsersAndTerms(dataPath, maxCount, users, termsEn, "en")
		if tmpErr4 != nil {
			mx.Lock()
			resErr = errors.New(fmt.Sprintf("failed generating users&terms (en): %v", tmpErr4))
			mx.Unlock()
		}
	}()

	wg.Wait()

	return resErr
}

func GenerateDataParallel(dataPath string, maxCount int) error {
	partsOfSpeech, models, users, err := generationStep1(dataPath, maxCount)
	if err != nil {
		return err
	}

	termsRu, termsEn, err := generationStep2(dataPath, maxCount, partsOfSpeech, models)
	if err != nil {
		return err
	}

	err = generationStep3(dataPath, maxCount, termsRu, termsEn, users)
	if err != nil {
		return err
	}

	return nil
}

func GenerateDataSequential(dataPath string, maxCount int) error {
	users, err := GenerateUsers(dataPath, maxCount)
	if err != nil {
		return err
	}

	partsOfSpeech, err := GeneratePos(dataPath, 25)
	if err != nil {
		return err
	}

	models, err := GenerateModels(dataPath, maxCount/5, partsOfSpeech)
	if err != nil {
		return err
	}

	termsRu, err := GenerateTerms(dataPath, maxCount, models, "ru")
	if err != nil {
		return err
	}

	termsEn, err := GenerateTerms(dataPath, maxCount, models, "en")
	if err != nil {
		return err
	}

	contexts, err := GenerateContexts(dataPath, maxCount/5, append(termsRu, termsEn...))
	if err != nil {
		return err
	}

	err = GenerateModelsAndPos(dataPath, maxCount, models, partsOfSpeech)
	if err != nil {
		return err
	}

	err = GenerateContextsAndTerms(dataPath, maxCount, contexts, termsRu, "ru")
	if err != nil {
		return err
	}

	err = GenerateContextsAndTerms(dataPath, maxCount, contexts, termsEn, "en")
	if err != nil {
		return err
	}

	err = GenerateUsersAndTerms(dataPath, maxCount, users, termsRu, "ru")
	if err != nil {
		return err
	}

	err = GenerateUsersAndTerms(dataPath, maxCount, users, termsEn, "en")
	if err != nil {
		return err
	}

	return nil
}
