package database

import (
	"dblabs/database/types"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"
)

var initUpdatePostgres = false
var partsOfSpeechPostgres5s []types.PartOfSpeech
var postgresMx = sync.Mutex{}

var initUpdateRedis = false
var partsOfSpeechRedis5s []types.PartOfSpeech
var redisMx = sync.Mutex{}

func (d *Database) StorePosToRedis(partsOfSpeech []types.PartOfSpeech) error {
	data, err := json.Marshal(partsOfSpeech)
	if err != nil {
		return err
	}

	return d.redisClient.Set("pos_table", data, 0).Err()
}

func (d *Database) LoadPosFromRedis() ([]types.PartOfSpeech, error) {
	data, err := d.redisClient.Get("pos_table").Result()
	if err != nil {
		return nil, err
	}

	var partsOfSpeech []types.PartOfSpeech
	err = json.Unmarshal([]byte(data), &partsOfSpeech)
	if err != nil {
		return nil, err
	}

	return partsOfSpeech, nil
}

func (d *Database) SelectPosFromRedis(key string) (types.PartOfSpeech, error) {
	var partOfSpeech types.PartOfSpeech

	data, err := d.redisClient.Get(key).Result()
	if err != nil {
		return partOfSpeech, err
	}

	err = json.Unmarshal([]byte(data), &partOfSpeech)
	if err != nil {
		return partOfSpeech, err
	}

	return partOfSpeech, nil
}

func (d *Database) SelectPosRedisCachedOnce() ([]types.PartOfSpeech, error) {
	partsOfSpeech, err := d.LoadPosFromRedis()
	if err != nil {
		partsOfSpeech, err = d.SelectAllPos()
		if err != nil {
			return nil, err
		}

		err = d.StorePosToRedis(partsOfSpeech)
		if err != nil {
			return nil, err
		}

		log.Println("Value loaded from database.")

		return partsOfSpeech, nil
	}

	log.Println("Value loaded from cache.")

	return partsOfSpeech, nil
}

func (d *Database) SelectPosPostgresCached5s() ([]types.PartOfSpeech, error) {
	ch := make(chan bool)
	if !initUpdatePostgres {
		go func() {
			for {
				res, err := d.SelectAllPos()
				if err == nil {
					postgresMx.Lock()
					partsOfSpeechPostgres5s = res
					postgresMx.Unlock()
					ch <- true
				} else {
					ch <- false
				}
				time.Sleep(5 * time.Second)
			}
		}()

		initUpdatePostgres = true

		letRead := <-ch
		if !letRead {
			return nil, errors.New("cannot read query result")
		}
	}

	postgresMx.Lock()
	partsOfSpeech := partsOfSpeechPostgres5s
	postgresMx.Unlock()

	return partsOfSpeech, nil
}

func (d *Database) SelectPosRedisCached5s() ([]types.PartOfSpeech, error) {
	ch := make(chan bool)
	if !initUpdateRedis {
		res, err := d.SelectAllPos()
		err = d.StorePosToRedis(res)
		if err != nil {
			return nil, err
		}

		go func() {
			for {
				res, err = d.LoadPosFromRedis()
				if err == nil {
					postgresMx.Lock()
					partsOfSpeechRedis5s = res
					postgresMx.Unlock()
					ch <- true
				} else {
					ch <- false
				}
				time.Sleep(5 * time.Second)
			}
		}()

		initUpdateRedis = true

		letRead := <-ch
		if !letRead {
			return nil, errors.New("cannot read query result")
		}
	}

	redisMx.Lock()
	partsOfSpeech := partsOfSpeechRedis5s
	redisMx.Unlock()

	return partsOfSpeech, nil
}

func (d *Database) InsertPosPostgres(partOfSpeech types.PartOfSpeech) error {
	script, err := Asset("database/sql/insert_pos.sql")
	if err != nil {
		return err
	}

	_, err = d.sqlDB.Query(string(script), partOfSpeech.Name, partOfSpeech.RegDate)

	return err
}

func (d *Database) UpdatePosPostgres(id int) error {
	script, err := Asset("database/sql/update_pos.sql")
	if err != nil {
		return err
	}

	_, err = d.sqlDB.Query(string(script), id)

	return err
}

func (d *Database) DeletePosPostgres(id int) error {
	script, err := Asset("database/sql/delete_pos.sql")
	if err != nil {
		return err
	}

	_, err = d.sqlDB.Query(string(script), id)

	return err
}

func (d *Database) InsertPosRedis(key string, partOfSpeech types.PartOfSpeech) error {
	data, err := json.Marshal(partOfSpeech)
	if err != nil {
		return err
	}

	return d.redisClient.Set(key, data, 0).Err()
}

func (d *Database) UpdatePosRedis(key string) error {
	data, err := d.redisClient.Get(key).Result()
	if err != nil {
		return err
	}

	var partOfSpeech types.PartOfSpeech
	err = json.Unmarshal([]byte(data), &partOfSpeech)
	if err != nil {
		return err
	}

	partOfSpeech.Name = "a"

	newData, err := json.Marshal(partOfSpeech)
	if err != nil {
		return err
	}

	return d.redisClient.Set(key, newData, 0).Err()
}

func (d *Database) DeletePosRedis(key string) error {
	return d.redisClient.Del(key).Err()
}
