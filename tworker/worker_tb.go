package tworker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/cihub/seelog"
	"github.com/sryanyuan/tbspider/tconfig"
)

// tumblr xml struct
type tbXmlVideoInfo struct {
	XmlName xml.Name `xml:"video"`
	Options string   `xml:"data-crt-options,attr"`
}

type tbXmlVideoPlayer struct {
	XmlName   xml.Name       `xml:"video-player"`
	VideoInfo tbXmlVideoInfo `xml:"video"`
	MaxWidth  string         `xml:"max-width,attr"`
}

type tbXmlVideoSource struct {
	ContentType string `xml:"content-type"`
	Extension   string `xml:"extension"`
	Width       int    `xml:"width"`
	Height      int    `xml:"height"`
}

type tbXmlPostItem struct {
	XmlName     xml.Name           `xml:"post"`
	ID          string             `xml:"id,attr"`
	Type        string             `xml:"type,attr"`
	Video       []tbXmlVideoPlayer `xml:"video-player"`
	VideoSource tbXmlVideoSource   `xml:"video-source"`
}

type tbXmlPosts struct {
	XmlName xml.Name        `xml:"posts"`
	Start   int             `xml:"start,attr"`
	Total   int             `xml:"total,attr"`
	Posts   []tbXmlPostItem `xml:"post"`
}

type tbXmlRoot struct {
	XmlName xml.Name   `xml:"tumblr"`
	Version string     `xml:"version,attr"`
	Posts   tbXmlPosts `xml:"posts"`
}

// shared work task
type tbWorkTask struct {
	url string
}

var (
	sharedTbTaskQueue []*tbWorkTask
)

type WorkerTb struct {
	pool     *WorkerPool
	workerID int
}

func init() {
	registerWorker("tumblr", &WorkerTb{})
}

func getTbDataFromFile(fp string) ([]byte, error) {
	file, err := os.Open(fp)
	if nil != err {
		return nil, err
	}
	defer file.Close()
	rspData, err := ioutil.ReadAll(file)
	if nil != err {
		return nil, err
	}
	return rspData, nil
}

func getTbDataFromHTTP(name string, proxy string) ([]byte, error) {
	rsp, err := GetByProxy(fmt.Sprintf("http://%s.tumblr.com/api/read?num=2", name), proxy)
	if nil != err {
		return nil, err
	}

	rspData, err := ioutil.ReadAll(rsp.Body)
	if nil != err {
		return nil, err
	}
	defer rsp.Body.Close()

	return rspData, nil
}

func (w *WorkerTb) linfo(args ...interface{}) {
	l := fmt.Sprintln(args...)
	seelog.Info("WorkerTb[", w.workerID, "] : ", l)
}

func (w *WorkerTb) Init(workerID int, pool *WorkerPool) error {
	//config := tconfig.StoreConfig(nil)
	w.pool = pool
	w.workerID = workerID

	// here we initialize work task once
	if nil == sharedTbTaskQueue {
		seelog.Info("Initialize task queue, it may takes some time, please wait ...")
		rspData, err := getTbDataFromFile("get.log")
		if nil != err {
			return err
		}
		// get the root element
		var root tbXmlRoot
		if err = xml.Unmarshal(rspData, &root); nil != err {
			return err
		}
		// get the video source by regexp
		regV := regexp.MustCompile(`.*?file_(.*?)" type="video/mp4"`)
		regV.FindAllStringSubmatch(string(rspData), -1)

		// write result
		sharedTbTaskQueue = make([]*tbWorkTask, 0, 32)
	}

	return nil
}

func (w *WorkerTb) New() IWorker {
	n := &WorkerTb{}
	return n
}

func (w *WorkerTb) Run() {
	w.linfo("Running ...")
	defer func() {
		w.pool.wg.Done()
	}()

	// we get the worker count
	workerCount := tconfig.StoreConfig(nil).MaxWorkers
	totalTaskCount := len(sharedTbTaskQueue)
	currentTaskCount := totalTaskCount / workerCount
	currentWorkingStartIndex := w.workerID * currentTaskCount
	reminderTaskCount := totalTaskCount % workerCount
	if w.workerID == workerCount-1 {
		// last one, should do the left work
		currentTaskCount += reminderTaskCount
	}

	for taskIndex := currentWorkingStartIndex; taskIndex < currentWorkingStartIndex+currentTaskCount; taskIndex++ {
		//task := sharedTbTaskQueue[taskIndex]

		// do get
	}

	w.linfo("Done ...")
}
