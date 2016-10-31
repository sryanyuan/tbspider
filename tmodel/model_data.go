package tmodel

import (
	"fmt"
)

const (
	spiderRecordModelName = "SpiderRecordModel"
)

type SpiderRecordModel struct {
	ID          int
	ResourceID  int
	ResourceTag string
	WorkerName  string
	Title       string
	Source      string
	Size        int
}

func (s *SpiderRecordModel) TableName() string {
	return spiderRecordModelName
}

func (s *SpiderRecordModel) Prepare() error {
	db := getDBConn()
	if nil == db {
		return fmt.Errorf("Nil db connection")
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + spiderRecordModelName + `(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		resource_id INTEGER UNIQUE KEY,
		resource_tag VARCHAR(128),
        worker_name VARCHAR(128),
        title VARCHAR(256),
        source VARCHAR(512),
        size INTEGER
    )`)
	if nil != err {
		return err
	}

	return nil
}

func InsertSpiderRecord(s *SpiderRecordModel) error {
	db := getDBConn()
	if nil == db {
		return fmt.Errorf("Nil db connection")
	}

	_, err := db.Exec("INSERT INTO "+spiderRecordModelName+` VALUES (
		null,
		?,
		?,
		?,
		?,
		?,
		?
	)`, s.ResourceID, s.ResourceTag, s.WorkerName, s.Title, s.Source, s.Size)
	if nil != err {
		return err
	}

	return nil
}

func init() {
	var record SpiderRecordModel
	registerModel(spiderRecordModelName, &record)
}
