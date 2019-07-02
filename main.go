package main

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"live-auto/src"
	"net/url"
)

func main() {
	err := src.ReadConfig("./config.xml")
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(src.G_Config)
	u, _ := url.Parse("https://www.douyu.com/265438")
	live, err := api.NewLive(u)
	if nil != err {
		fmt.Println(err)
		return
	}
	l, err := live.GetStreamUrls()
	download := src.NewDownloader(l[0].String(), "./小缘/123.flv")
	download.Start()
}
