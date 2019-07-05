package main

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"live-auto/src"
	"net/url"
)

func main() {
	// defer src.Logger.Exit(0)
	u, _ := url.Parse("https://www.douyu.com/2550505")
	live, err := api.NewLive(u)
	if nil != err {
		src.Logger.Error(err.Error())
		return
	}
	l, err := live.GetStreamUrls()
	download := src.NewDownloader(l[0].String(), "./test/123.flv")
	go download.Start()
	select {
	case cb_t := <-download.CBChannel:
		fmt.Println(cb_t)
	}
	fmt.Println("main")
}
