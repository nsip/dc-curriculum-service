// db.go

//
// captures all db interactions
//

package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
)

var buntDB *buntdb.DB

//
// opens the database on the file system
// and loads all curriculm json resources
//
func initialiseDB() {

	// Open the data.db file. It will be created if it doesn't exist.
	var dbErr error
	buntDB, dbErr = buntdb.Open("./db/data.db")
	if dbErr != nil {
		log.Fatal("error opening database: ", dbErr)
	}
	log.Println("database opened.")

	// walk the directory of json files & commit them to db
	root := "./nsw"
	err := filepath.Walk(root, visitAndCommit)
	log.Printf("filepath.Walk() for %s complete returned error %v\n", root, err)

}

//
// compacts and closes the db
//
func closeDB() {
	err := buntDB.Shrink()
	if err != nil {
		log.Println("error compacting database: ", err)
	}
	err = buntDB.Close()
	if err != nil {
		log.Println("error closing database: ", err)
	}
	log.Println("...database closed")
}

//
// called for each file the curriculum directory walker encounters
// if a .json file is found, content committed to db.
//
func visitAndCommit(path string, fi os.FileInfo, err error) error {

	//
	// make sure we only list .json files
	//
	if fi.Mode().IsRegular() && strings.HasSuffix(path, ".json") {
		commitErr := commitJSON(path)
		if commitErr != nil {
			return commitErr
		}
		// log.Println("successfully committed: ", path)
	}

	return nil
}

//
// reads the content of a json file and commits it to the
// local db with a key derived from the file path.
//
func commitJSON(path string) error {

	// read the file
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// retrieve bytes into json string
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	jsonString := gjson.ParseBytes(jsonBytes).Raw
	if !gjson.Valid(jsonString) {
		log.Println("could not read valid json from file: ", path)
	}

	// create the key for the json data string
	keyName := strings.Replace(strings.TrimSuffix(path, ".json"), string(os.PathSeparator), "-", -1)
	// log.Println("derived key:\t", keyName)

	// commit to the db
	txErr := buntDB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(keyName, jsonString, nil)
		return err
	})

	return txErr
}

//
// for the given key returns a map[string]interface{} rendering
// of the relevant json in the db
//
func getJSONMap(key string) (map[string]interface{}, error) {

	var result gjson.Result

	err := buntDB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		result = gjson.Parse(val)
		return nil
	})
	if err != nil {
		return nil, err
	}

	jsonMap, ok := result.Value().(map[string]interface{})
	if !ok {
		return nil, errors.New("value from db is not a json object")
	}

	return jsonMap, nil
}
