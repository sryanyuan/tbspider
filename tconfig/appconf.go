package tconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cihub/seelog"
	"github.com/go-validator/validator"
	"github.com/sryanyuan/tbspider/tconstant"
)

type AppConfig struct {
	DBResultAddress string `validate:"nonzero"`
	MaxWorkers      int
	WorkerName      string
	ProxyAddress    string
	SpiderKeyword   string
}

// global share config
var shareConfig *AppConfig

func StoreConfig(conf *AppConfig) *AppConfig {
	if nil != conf {
		shareConfig = conf
	}

	return shareConfig
}

// InitAppConfig
func InitAppConfig(path string) (*AppConfig, error) {
	f, err := os.Open(path)
	if nil != err {
		seelog.Error("Can't open config file:", path)
		return nil, err
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if nil != err {
		seelog.Error("Can't read config file:", path)
		return nil, err
	}

	//	parse the json
	var config AppConfig
	if err = json.Unmarshal(fileBytes, &config); nil != err {
		return nil, err
	}

	// validate
	if err = validator.Validate(&config); nil != err {
		return nil, err
	}

	if 0 == config.MaxWorkers {
		config.MaxWorkers = tconstant.DefaultMaxWorkers
	}

	StoreConfig(&config)
	return &config, err
}
