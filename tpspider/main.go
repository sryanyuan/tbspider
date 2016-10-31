package main

import (
	"flag"
	"log"

	"github.com/cihub/seelog"
	"github.com/sryanyuan/tbspider/tconfig"
	"github.com/sryanyuan/tbspider/tconstant"
	"github.com/sryanyuan/tbspider/tmodel"
	"github.com/sryanyuan/tbspider/tworker"
)

func main() {
	// initialize log
	err := tconfig.InitLogModule()
	if nil != err {
		log.Println("[ERR] ", err.Error())
		return
	}

	defer func() {
		seelog.Flush()
	}()

	// parse flag
	configPath := flag.String("configpath", "", "-configpath <json config file path>")
	flag.Parse()

	if len(*configPath) == 0 {
		// using default path
		*configPath = tconstant.DefaultConfigFile
		seelog.Warn("No config file specified (-configpath <config file>) , using default path : ", *configPath)
	}

	// read the config
	var appConfig *tconfig.AppConfig
	if appConfig, err = tconfig.InitAppConfig(*configPath); nil != err {
		seelog.Error("Cannot read config from config file : ", *configPath, " , error : ", err)
		return
	}

	// init model
	if err = tmodel.Init(appConfig.DBResultAddress); nil != err {
		seelog.Error("Model initialize error : ", err)
		return
	}

	// run all workers
	pool := tworker.NewWorkerPool()
	if err = pool.InitWithWorkerCount(appConfig.MaxWorkers); nil != err {
		seelog.Error("Pool initialize failed : ", err.Error())
		return
	}

	//	listen for SIG
	/*sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh*/

	// do clean up
	pool.WaitWorkersDone()
	tmodel.Shutdown()
}
