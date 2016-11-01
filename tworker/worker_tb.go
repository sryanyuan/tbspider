package tworker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/cihub/seelog"
	"github.com/sryanyuan/tbspider/tconfig"
	"github.com/sryanyuan/tbspider/tmodel"
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
	XmlName   xml.Name `xml:"post"`
	ID        string   `xml:"id,attr"`
	Type      string   `xml:"type,attr"`
	Slug      string   `xml:"slug,attr"`
	Signature string
	//Video       []tbXmlVideoPlayer `xml:"video-player"`
	//VideoSource tbXmlVideoSource   `xml:"video-source"`
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
	tbPostItem *tbXmlPostItem
}

var (
	sharedTbTaskQueue     []*tbWorkTask
	sharedTbTaskTotalSize int
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

func getTbDataFromHTTP(name string, proxy string, start int, num int) ([]byte, error) {
	rsp, err := GetByProxy(fmt.Sprintf("http://%s.tumblr.com/api/read?start=%d&num=%d", name, start, num), proxy)
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
	config := tconfig.StoreConfig(nil)
	config.DBResultAddress = config.DBResultAddress
	w.pool = pool
	w.workerID = workerID

	// here we initialize work task once
	if 0 == sharedTbTaskTotalSize {
		seelog.Info("Initialize task queue, it may takes some time, please wait ...")
		rspData, err := getTbDataFromFile("./get.log")
		//rspData, err := getTbDataFromHTTP(config.SpiderKeyword, config.ProxyAddress, 0, 1)
		if nil != err {
			return err
		}
		// get the root element
		var root tbXmlRoot
		if err = xml.Unmarshal(rspData, &root); nil != err {
			return err
		}
		// get the video total number
		sharedTbTaskTotalSize = root.Posts.Total
		if 0 == sharedTbTaskTotalSize {
			return fmt.Errorf("Get empty post number")
		}
		seelog.Info("Get post number done, size : ", sharedTbTaskTotalSize)
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
	config := tconfig.StoreConfig(nil)
	workerCount := config.MaxWorkers
	totalTaskCount := sharedTbTaskTotalSize
	currentTaskCount := totalTaskCount / workerCount
	currentWorkingStartIndex := w.workerID * currentTaskCount
	currentWorkingStartIndex = currentWorkingStartIndex
	reminderTaskCount := totalTaskCount % workerCount
	if w.workerID == workerCount-1 {
		// last one, should do the left work
		currentTaskCount += reminderTaskCount
	}

	// do get
	//rspData, err := getTbDataFromFile("./get.log")
	rspData, err := getTbDataFromHTTP(config.SpiderKeyword, config.ProxyAddress, currentWorkingStartIndex, currentTaskCount)
	if nil != err {
		seelog.Error(err)
		return
	}
	// get the root element
	var root tbXmlRoot
	if err = xml.Unmarshal(rspData, &root); nil != err {
		seelog.Error(err)
		return
	}
	// get the video source by regexp
	//regV := regexp.MustCompile(`.*?/(\d{10,20})/tumblr_(.*?)" type="video/mp4"`)
	regV := regexp.MustCompile(`.*?(\d{2}).media.tumblr.com.*?&gt;\s*&lt;source.*?/(\d{10,20})/tumblr_(.*?)" type="video/mp4"`)
	vResults := regV.FindAllStringSubmatch(string(rspData), -1)
	videoSignatures := make(map[string]string)
	for _, v := range vResults {
		if len(v) == 4 {
			videoSignatures[v[2]] = v[3] + "&" + v[1]
		}
	}
	// write result
	sharedTbTaskQueue = make([]*tbWorkTask, 0, len(videoSignatures))
	// fill the root post item from video signatures
	for i, post := range root.Posts.Posts {
		if post.Type == "video" {
			if vinfo, ok := videoSignatures[post.ID]; ok {
				slist := strings.Split(vinfo, "&")
				if len(slist) != 2 {
					w.linfo("invalid vinfo")
					continue
				}
				vsource := slist[0]
				root.Posts.Posts[i].Signature = vsource
				var record tmodel.SpiderRecordModel
				record.ResourceID = post.ID
				record.ResourceTag = config.SpiderKeyword
				record.ResourceType = post.Type
				record.ResourceSig = vsource
				record.ResourceImg = fmt.Sprintf("http://%s.media.tumblr.com/previews/tumblr_%s_filmstrip.jpg", slist[1], vsource)
				record.Source = fmt.Sprintf("https://vt.tumblr.com/tumblr_%s.mp4", vsource)
				record.WorkerName = "tumblr"
				record.Title = post.Slug

				err = tmodel.InsertSpiderRecord(&record)
				if nil != err {
					w.linfo(err.Error())
				}
			}
		}
	}

	w.linfo("Done ...")
}
