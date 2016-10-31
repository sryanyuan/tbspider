package tmodel

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// model
type IModel interface {
	Prepare() error // for creating table
	TableName() string
}

// db connections
var (
	sharedDBConn       *sql.DB
	registeredModelMap map[string]IModel
)

func init() {
	registeredModelMap = make(map[string]IModel)
}

func registerModel(modelName string, model IModel) {
	registeredModelMap[modelName] = model
}

func Init(dataSourceName string) error {
	db, err := sql.Open("sqlite3", dataSourceName)
	if nil != err {
		return err
	}

	if err = db.Ping(); nil != err {
		return err
	}
	sharedDBConn = db

	// prepare all models
	for _, v := range registeredModelMap {
		err = v.Prepare()
		if nil != err {
			return err
		}
	}

	return nil
}

func Shutdown() {
	if nil != sharedDBConn {
		sharedDBConn.Close()
		sharedDBConn = nil
	}
}

func getDBConn() *sql.DB {
	return sharedDBConn
}
