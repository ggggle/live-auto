package main

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"live-auto/src"
	"net/url"
	"time"
)

func main() {
	fmt.Println("main")

	src.Logger.Error("123")
	time.Sleep(time.Second)
	src.Logger.Error("222")
	return
	u, _ := url.Parse("https://www.douyu.com/233233")
	live, err := api.NewLive(u)
	if nil != err {
		src.Logger.Error(err.Error())
		return
	}
	l, err := live.GetStreamUrls()
	download := src.NewDownloader(l[0].String(), "./test/123.flv")
	download.Start()
}
