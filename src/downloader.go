package src

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type DownloadCBCodeType int

// downloader返回值定义
const (
	NORMAL_END  DownloadCBCodeType = iota // 通常的情况  live流中断、结束
	WRITE_ERROR                           // 写入硬盘错误 空间满了
	NEXT_PIECE                            // 单片段达到大小上限
	UNSTART_ERROR
	STOP_SELF // 主动停止
	LIVE_STREAM_NIL
)

type downloader struct {
	// 直播流URL
	url string
	// 保存位置
	liveFilePath string
	// 是否使用ffmpeg
	needFFmpeg bool
	CBChannel  chan downloadCB_t
	// 停止
	Stop bool
}

type downloadCB_t struct {
	Code DownloadCBCodeType
}

func NewDownloader(url, filePath string) *downloader {
	if "" != url && "" != filePath {
		return &downloader{
			url:          url,
			liveFilePath: filePath,
			needFFmpeg:   false,
			CBChannel:    make(chan downloadCB_t, 1),
			Stop:         false,
		}
	}
	return nil
}

func (self *downloader) Start() {
	download_cb := downloadCB_t{
		Code: NORMAL_END,
	}
	// 通知结果
	defer func() {
		self.CBChannel <- download_cb
	}()
	resp, err := httpDownload(self.url)
	if nil != err {
		Logger.WithFields(logrus.Fields{
			ERROR_CONTENT_DEF: err.Error(),
			"url":             self.url,
		}).Error("下载错误")
		download_cb.Code = UNSTART_ERROR
		return
	}
	defer resp.Body.Close()
	f_handle, err := createFile(self.liveFilePath)
	if nil != err {
		Logger.WithFields(logrus.Fields{
			ERROR_CONTENT_DEF: err.Error(),
			"file_path":       self.liveFilePath,
		}).Error("创建本地文件错误")
		download_cb.Code = UNSTART_ERROR
		return
	}
	defer f_handle.Close()
	buffer := make([]byte, bytes.MinRead)
	once_end := make(chan struct{})
	go func() {
		for !self.Stop {
			t, err := resp.Body.Read(buffer)
			if t > 0 {
				if _, err := f_handle.Write(buffer[:t]); nil != err {
					// 写入文件错误
					Logger.WithFields(logrus.Fields{
						ERROR_CONTENT_DEF: err.Error(),
					}).Error("直播文件写入错误")
					download_cb.Code = WRITE_ERROR
					break
				}
				if G_Config.PieceSize > 0 {
					if file_info, _ := f_handle.Stat(); file_info.Size() >= G_Config.PieceSize {
						Logger.WithFields(logrus.Fields{
							"file_size": file_info.Size(),
						}).Info("切片存储")
						download_cb.Code = NEXT_PIECE
						break
					}
				}
			}
			// io.EOF直播流到结尾 可能真的结束了或者网络波动导致流断开
			if nil != err {
				if io.EOF != err {
					Logger.WithFields(logrus.Fields{
						ERROR_CONTENT_DEF: err.Error(),
					}).Info("直播流中断")
				}
				break
			}
		}
		if self.Stop {
			download_cb.Code = STOP_SELF
		}
		once_end <- struct{}{}
	}()
	<-once_end
	if f_info, _ := f_handle.Stat(); 0 == f_info.Size() {
		// 开直播了但是还没有流
		time.Sleep(time.Second * 5)
		download_cb.Code = LIVE_STREAM_NIL
		return
	}
	Logger.WithFields(logrus.Fields{
		"file_path": self.liveFilePath,
	}).Debug("once_end")
	return
}

func httpDownload(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if nil != err {
		return nil, err
	}
	req.Header.Add("User-Agent", "Chrome/59.0.3071.115")
	client := http.Client{}
	resp, err = client.Do(req)
	return
}
