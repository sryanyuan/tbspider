package tconfig

import (
	"log"
	"os"

	"github.com/cihub/seelog"
)

//	global const variables
var (
	kDefaultLogSetting = `
	<seelog minlevel="debug">
    	<outputs formatid="main">
			<rollingfile namemode="postfix" type="date" filename="log/app.log" datepattern="060102" maxrolls="30"/>
       		<console />
    	</outputs>
    	<formats>
        	<format id="main" format="%Date/%Time [%LEV] %Msg (%File:%Line %FuncShort)%n"/>
    	</formats>
	</seelog>
	`
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Init log module
func InitLogModule() error {
	var err error

	logFilePath := "./conf/log.conf"
	var logger seelog.LoggerInterface
	logFileExist, _ := pathExists(logFilePath)

	if !logFileExist {
		//	using the default setting
		log.Printf("[WRN] Can't open %s, using the default log setting", logFilePath)
		logger, err = seelog.LoggerFromConfigAsString(kDefaultLogSetting)
		if nil != err {
			log.Println("[ERR] Error on initialize log")
			return err
		}
	} else {
		logger, err = seelog.LoggerFromConfigAsFile(logFilePath)
		if nil != err {
			log.Println("[ERR] Error on initialize log")
			return err
		}
	}
	seelog.ReplaceLogger(logger)
	return nil
}
