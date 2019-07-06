package src

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"github.com/hr3lxphr6j/bililive-go/src/lib/utils"
	"github.com/sirupsen/logrus"
	"net/url"
	"path/filepath"
	"time"
)

type Recorder struct {
	RecordConfig
	Live api.Live
	Stop chan struct{}
}

type RecordConfig struct {
	// 循环
	Loop bool
	// TODO
	BeginTime int
	EndTime   int
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
	return &Recorder{
		Live:         live,
		RecordConfig: config,
		Stop:         make(chan struct{}, 1),
	}, nil
}

func (self *Recorder) Start() {
	go self.run()
	Logger.WithFields(self.Live.GetInfoMap()).Info("Record Start")
}

func (self *Recorder) run() {
	for {
		info, err := self.Live.GetInfo()
		if nil != err || !info.Status {
			time.Sleep(time.Duration(G_Config.CheckInterval) * time.Second)
			continue
		}
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
		downloader := NewDownloader(live_url.String(), live_file_path)
		go downloader.Start()
		select {
		case cb_t := <-downloader.CBChannel:
			switch cb_t.Code {
			case STOP_SELF, WRITE_ERROR, UNSTART_ERROR:
				return
			default:
				continue
			}
		case <-self.Stop:
			downloader.Stop = true
			return
		}
	}
}
