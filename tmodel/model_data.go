package tmodel

import (
	"fmt"
	"os"

	"github.com/cihub/seelog"
)

const (
	spiderRecordModelName = "SpiderRecordModel"
)

type SpiderRecordModel struct {
	ID           int
	ResourceID   string
	ResourceTag  string
	ResourceType string
	ResourceSig  string
	ResourceImg  string
	WorkerName   string
	Title        string
	Source       string
	Size         int
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
		resource_id INTEGER UNIQUE,
		resource_tag VARCHAR(128),
		resource_type VARCHAR(64),
		resource_sig VARCHAR(64),
		resource_img VARCHAR(128),
        worker_name VARCHAR(128),
        title VARCHAR(512),
        source VARCHAR(128),
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
		?,
		?,
		?,
		?
	)`, s.ResourceID, s.ResourceTag, s.ResourceType, s.ResourceSig, s.ResourceImg, s.WorkerName, s.Title, s.Source, s.Size)
	if nil != err {
		return err
	}

	return nil
}

func DumpSpiderRecordToFileByTag(fileName string, tag string) error {
	f, err := os.Create(fileName)
	if nil != err {
		return err
	}
	defer f.Close()

	db := getDBConn()
	if nil == db {
		return fmt.Errorf("Nil db connection")
	}

	rows, err := db.Query("SELECT resource_img, source FROM "+spiderRecordModelName+" WHERE resource_tag = ?", tag)
	if nil != err {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var img, source string
		if err = rows.Scan(&img, &source); nil != err {
			return err
		}

		seelog.Debug("write ", img, source)

		f.WriteString(img)
		f.WriteString("\r\n")
		f.WriteString(source)
		f.WriteString("\r\n")
	}

	return nil
}

func init() {
	var record SpiderRecordModel
	registerModel(spiderRecordModelName, &record)
}
