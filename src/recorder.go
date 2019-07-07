package src

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"github.com/hr3lxphr6j/bililive-go/src/lib/utils"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Recorder struct {
	RecordConfig
	Live      api.Live
	stop      chan struct{}
	Uploaders []Uploader
	IndexID   int // 递增int标识
}

type RecordConfig struct {
	// 循环
	Loop bool
	// TODO
	BeginTime       int
	EndTime         int
	EnableUploaders []UploaderType
	AutoRemove      bool // 上传成功后删除本地文件
}

func NewRecorder(live_url string, config RecordConfig) (*Recorder, error) {
	u, err := url.Parse(live_url)
	if nil != err {
		return nil, err
	}
	live, err := api.NewLive(u)
	if nil != err {
		return nil, err
	}
	recorder := Recorder{
		Live:         live,
		RecordConfig: config,
		stop:         make(chan struct{}, 1),
		Uploaders:    make([]Uploader, 0),
	}
	for _, _type := range recorder.EnableUploaders {
		recorder.Uploaders = append(recorder.Uploaders, NewUploader(_type, &recorder))
	}
	Logger.WithFields(logrus.Fields{
		"live_url": live_url,
		"config":   config,
	}).Info("Add Recorder")
	return &recorder, nil
}

func (self *Recorder) Start() {
	go self.run()
	Logger.WithFields(self.Live.GetInfoMap()).Info("Record Monitor Start")
}

func (self *Recorder) Stop() {
	Logger.WithFields(self.Live.GetInfoMap()).Info("Record Monitor Stop")
	self.Loop = false
	self.stop <- struct{}{}
}

func (self *Recorder) run() {
	for {
		info, err := self.Live.GetInfo()
		if nil != err || !info.Status {
			// 未开播
			time.Sleep(time.Duration(G_Config.CheckInterval) * time.Second)
			continue
		}
		// 开播后获取直播流
		urls, err := self.Live.GetStreamUrls()
		if nil != err || 0 == len(urls) {
			Logger.WithFields(logrus.Fields{
				ERROR_CONTENT_DEF: err.Error(),
				"url":             self.Live.GetRawUrl(),
			}).Error("直播流获取失败")
			time.Sleep(time.Second * 5)
			continue
		}
		var (
			platformName   = utils.ReplaceIllegalChar(self.Live.GetPlatformCNName())
			hostName       = utils.ReplaceIllegalChar(self.Live.GetCachedInfo().HostName)
			roomName       = utils.ReplaceIllegalChar(self.Live.GetCachedInfo().RoomName)
			fileName       = fmt.Sprintf("[%s][%s][%s].flv", time.Now().Format("2006-01-02 15-04-05"), hostName, roomName)
			outputPath     = filepath.Join(G_Config.OutPutDirPath, platformName, hostName)
			live_file_path = filepath.Join(outputPath, fileName)
			live_url       = urls[0]
		)
		Logger.WithFields(self.Live.GetInfoMap()).Info("直播开始录制")
		downloader := NewDownloader(live_url.String(), live_file_path)
		go downloader.Start()
	SELECT:
		select {
		case cb_t := <-downloader.CBChannel:
			Logger.WithFields(logrus.Fields{
				"cb_code": cb_t.Code,
			}).Debug("download return")
			self.doUpload(live_file_path)
			switch cb_t.Code {
			case LIVE_STREAM_NIL:
				os.Remove(live_file_path)
				continue
			case WRITE_ERROR:
				self.Loop = false
			}
			if self.Loop {
				continue
			} else {
				return
			}
		case <-self.stop:
			downloader.Stop = true
			// 主动停止 等待下载协程返回后处理上传
			goto SELECT
		}
	}

}

// 调用上传
func (self *Recorder) doUpload(file_path string) {
	for i := range self.Uploaders {
		if nil != self.Uploaders[i] {
			go self.Uploaders[i].DoUpload(file_path)
		}
	}
}
