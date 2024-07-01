package database

import (
	"fmt"
	"reflect"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) CreateChirp(body string, author int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	var keyMax int
	for _, chirp := range dbStructure.Chirps {
		if chirp.ID > keyMax {
			keyMax = chirp.ID
		}
	}

	id := keyMax + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorId: author,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}
	return chirp, nil
}

// authentication should have already happened before getting here
func (db *DB) DeleteChirp(id, author int) (string, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return "", err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return "", ErrNotExist
	}

	r := reflect.ValueOf(&chirp).Elem()
	rt := r.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		rv := reflect.ValueOf(&chirp)
		value := reflect.Indirect(rv).FieldByName(field.Name)
		fmt.Printf("--> field: %v, value: %v \n", field.Name, value)
	}

	if chirp.AuthorId != author {
		return "wrong user", nil
	}

	delete(dbStructure.Chirps, id)
	db.writeDB(dbStructure)
	return "success", nil
}

func (chrp *Chirp) GetAuthor() int {
	return chrp.AuthorId
}
