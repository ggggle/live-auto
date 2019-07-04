package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type downloader struct {
	// 直播流URL
	url string
	// 保存位置
	liveFilePath string
	// 是否使用ffmpeg
	needFFmpeg bool
}

func NewDownloader(url, filePath string) *downloader {
	if "" != url && "" != filePath {
		return &downloader{
			url:          url,
			liveFilePath: filePath,
			needFFmpeg:   false,
		}
	}
	return nil
}

func (self *downloader) Start() {
	resp, err := httpDownload(self.url)
	if nil != err {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	f_handle, err := createFile(self.liveFilePath)
	if nil != err {
		fmt.Println(err)
		return
	}
	defer f_handle.Close()
	buffer := make([]byte, bytes.MinRead)
	go func() {
		for true {
			t, err := resp.Body.Read(buffer)
			// fmt.Println(t)
			if t > 0 {
				if _, err := f_handle.Write(buffer[:t]); nil != err {
					fmt.Println(err)
					// 写入文件错误
					break
				}
			}
			if nil != err {
				fmt.Println(err)
				// 直播流到结尾 正常结束
				if io.EOF == err {

				}
				break
			}
		}
	}()
	select {}
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
